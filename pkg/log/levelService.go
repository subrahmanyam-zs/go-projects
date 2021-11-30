package log

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type levelService struct {
	logger       Logger
	init         bool
	level        level
	url          string
	app          string
	failureCount int
}

// nolint
var (
	rls levelService // TODO - remove this
	mu  sync.RWMutex
)

const LevelFetchInterval = 10 // In seconds

func newLevelService(l Logger, appName string) *levelService {
	if !rls.init {
		lvl := getLevel(os.Getenv("LOG_LEVEL"))

		mu.Lock()

		rls.level = lvl

		mu.Unlock()

		rls.url = os.Getenv("LOG_SERVICE_URL")
		rls.app = appName
		rls.logger = l

		if rls.url != "" {
			rls.init = true

			go func() {
				for {
					rls.updateRemoteLevel()
					time.Sleep(LevelFetchInterval * time.Second)
				}
			}()
		}
	}

	return &rls
}

func (s *levelService) updateRemoteLevel() {
	rls.logger.Debugf("Making request to remote logging service %s", s.url)

	req, _ := http.NewRequest(http.MethodGet, s.url+"/configs?serviceName="+s.app, http.NoBody)

	tr := &http.Transport{
		//nolint:gosec // need this to skip TLS verification
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := (&http.Client{Transport: tr}).Do(req)
	if err != nil {
		s.logger.Warnf("Could not create log service client. err:%v", err)
		s.failureCount++

		return
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Warnf("Logging Service returned %d status. Req: %s", resp.StatusCode, req.URL)

		return
	}

	if resp.Body != nil {
		b, _ := io.ReadAll(resp.Body)

		_ = resp.Body.Close()

		if newLevel := s.getRemoteLevel(b); s.level != newLevel {
			s.logger.Debugf("Changing log level from %s to %s because of remote log service", s.level, newLevel)

			s.level = newLevel
		}
	}
}

func (s *levelService) getRemoteLevel(body []byte) level {
	type data struct {
		ServiceName string            `json:"serviceName"`
		Config      map[string]string `json:"config"`
		UserGroup   string            `json:"userGroup"`
	}

	level := struct {
		Data []data `json:"data"`
	}{}

	err := json.Unmarshal(body, &level)
	if err != nil {
		s.logger.Warnf("Logging Service returned %v", err)
	}

	if level.Data != nil {
		logLevel := level.Data[0].Config["LOG_LEVEL"]
		newLevel := getLevel(logLevel)

		return newLevel
	}

	return s.level
}
