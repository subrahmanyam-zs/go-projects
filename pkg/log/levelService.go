package log

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
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
var rls levelService          // TODO - remove this
const LevelFetchInterval = 10 // In seconds

func newLevelService(l Logger, appName string) *levelService {
	if !rls.init {
		rls.level = getLevel(os.Getenv("LOG_LEVEL"))
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
	req, _ := http.NewRequest("GET", s.url+"/level?service="+s.app, nil)

	var tr = &http.Transport{
		//nolint:gosec // need this to skip TLS verification
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := (&http.Client{Transport: tr}).Do(req)
	if err != nil {
		s.logger.Warn("Could not create log service client. err: ", err)
		s.failureCount++

		return
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Warnf("Logging Service returned %d status. Req: %s", resp.StatusCode, req.URL)
		return
	}

	level := struct {
		Data string `json:"data"`
	}{}

	if resp.Body != nil {
		body, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &level)
		_ = resp.Body.Close()
	}

	if newLevel := getLevel(level.Data); s.level != newLevel {
		s.logger.Debugf("Changing log level from %s to %s because of remote log service", s.level, newLevel)
		s.level = newLevel
	}
}
