package log

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRemoteLevelLogger(t *testing.T) {
	tests := []struct {
		desc        string
		level       level
		serviceName string
		body        []byte
	}{
		{"success case", Info, "gofr-sample-api", []byte(`{"data": [{"serviceName": "gofr-sample-api","config": {"LOG_LEVEL": "INFO"}}]}`)},
		{"failure case", Debug, "", nil},
	}

	for i, tc := range tests {
		// test server that returns log level for the app
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(tc.body)
		}))

		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/configs?serviceName="+tc.serviceName, http.NoBody)

		tr := &http.Transport{
			//nolint:gosec // need this to skip TLS verification
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		rls.logger = NewMockLogger(io.Discard)

		s := &levelService{url: ts.URL, logger: rls.logger}

		s.level = Debug

		s.updateRemoteLevel(client, req)

		ts.Close()

		assert.Equal(t, tc.level, s.level, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestRemoteLevelLoggerRequestError(t *testing.T) {
	// test server that returns log level for the app
	b := new(bytes.Buffer)
	l := NewMockLogger(b)

	rls.logger = l

	req, _ := http.NewRequest(http.MethodGet, "", http.NoBody)
	client := &http.Client{}

	s := &levelService{url: "", logger: l}

	s.updateRemoteLevel(client, req)

	assert.Contains(t, b.String(), "Could not create log service client")
}

func TestRemoteLevelLoggerNoResponse(t *testing.T) {
	// test server that returns log level for the app
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/configs?serviceName=", http.NoBody)

	tr := &http.Transport{
		//nolint:gosec // need this to skip TLS verification
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	defer ts.Close()

	b := new(bytes.Buffer)
	l := NewMockLogger(b)

	rls.logger = l

	s := &levelService{url: ts.URL, logger: l}

	s.updateRemoteLevel(client, req)

	expectedLog := "Logging Service returned 404 status. Req: " + ts.URL

	if !strings.Contains(b.String(), expectedLog) {
		t.Errorf("expected error")
	}
}

func TestRemoteLevelLogging(t *testing.T) {
	body := []byte(`{"data": [{"serviceName": "gofr-sample-api","config": {"LOG_LEVEL": "WARN"}}]}`)
	// test server that returns log level for the app
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))

	defer ts.Close()

	t.Setenv("LOG_SERVICE_URL", ts.URL)

	b := new(bytes.Buffer)
	l := NewMockLogger(b)

	rls.logger = l

	newLevelService(l, "gofr-app")

	time.Sleep(15 * time.Second)

	mu.Lock()
	lvl := rls.level
	mu.Unlock()

	if lvl != Warn {
		t.Errorf("expected WARN\tGot %v", lvl)
	}

	if rls.app != "gofr-app" {
		t.Errorf("expected APP_NAME : test, Got : %v", rls.app)
	}
}
