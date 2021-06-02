package log

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRemoteLevelLogger(t *testing.T) {
	// test server that returns log level for the app
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		level := struct {
			Data string `json:"data"`
		}{Data: "debug"}
		reBytes, _ := json.Marshal(level)
		_, _ = w.Write(reBytes)
	}))

	defer ts.Close()

	rls.logger = NewMockLogger(io.Discard)
	s := &levelService{url: ts.URL, logger: rls.logger}

	s.updateRemoteLevel()

	if !reflect.DeepEqual(s.level, Debug) {
		t.Errorf("update remote logging failed")
	}
}

func TestRemoteLevelLoggerRequestError(t *testing.T) {
	// test server that returns log level for the app
	b := new(bytes.Buffer)
	l := NewMockLogger(b)
	rls.logger = l
	s := &levelService{url: "", logger: l}

	s.updateRemoteLevel()

	if !strings.Contains(b.String(), "Could not create log service client") {
		t.Errorf("expected error")
	}
}

func TestRemoteLevelLoggerNoResponse(t *testing.T) {
	// test server that returns log level for the app
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))

	defer ts.Close()

	b := new(bytes.Buffer)
	l := NewMockLogger(b)
	rls.logger = l
	s := &levelService{url: ts.URL, logger: l}

	s.updateRemoteLevel()

	expectedLog := "Logging Service returned 404 status. Req: " + ts.URL

	if !strings.Contains(b.String(), expectedLog) {
		t.Errorf("expected error")
	}
}

func TestRemoteLevelLogging(t *testing.T) {
	// test server that returns log level for the app
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		level := struct {
			Data string `json:"data"`
		}{Data: "warn"}
		reBytes, _ := json.Marshal(level)
		_, _ = w.Write(reBytes)
	}))

	defer ts.Close()

	os.Setenv("LOG_SERVICE_URL", ts.URL)

	b := new(bytes.Buffer)
	l := NewMockLogger(b)
	rls.logger = l

	newLevelService(l, "gofr-app")

	time.Sleep(15 * time.Second)

	if rls.level != Warn {
		t.Errorf("expected WARN\tGot %v", rls.level)
	}

	if rls.app != "gofr-app" {
		t.Errorf("expected APP_NAME : test, Got : %v", rls.app)
	}
}
