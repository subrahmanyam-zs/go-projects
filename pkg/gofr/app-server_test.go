package gofr

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestServer_Done(t *testing.T) {
	// start a server using Gofr
	k := New()
	k.Server.HTTP.Port = 8080

	go k.Start()
	time.Sleep(time.Second * 3)

	serverUP := false

	// check if server is up
	for i := 0; i < 2; i++ {
		resp, _ := http.Get("http://localhost:8080/.well-known/heartbeat")
		if resp.StatusCode == http.StatusOK {
			serverUP = true
			_ = resp.Body.Close()

			break
		}

		time.Sleep(time.Second)
	}

	if !serverUP {
		t.Errorf("server not up")
	}

	// stop the server
	k.Server.Done()

	serverUP = true

	// check if the server is down
	for i := 0; i < 3; i++ {
		//nolint:bodyclose // there is no response here hence body cannot be closed.
		_, err := http.Get("http://localhost:8080/.well-known/heartbeat")
		// expecting an error since server is down
		if err != nil {
			serverUP = false

			break
		}

		time.Sleep(time.Second)
	}

	if serverUP {
		t.Errorf("server down failed")
	}
}

// This tests if a server can be started again after being stopped.
func TestServer_Done2(t *testing.T) {
	TestServer_Done(t)
	TestServer_Done(t)
}

// Test_AllRouteLog will test logging of all routes of the server along with methods
func Test_AllRouteLog(t *testing.T) {
	k := New()
	k.Server.HTTP.Port = 8080

	b := new(bytes.Buffer)
	k.Logger = log.NewMockLogger(b)

	go k.Start()
	time.Sleep(time.Second * 2)
	assert.Contains(t, b.String(), "GET /.well-known/health-check HEAD /.well-known/health-check ")
	assert.Contains(t, b.String(), "GET /.well-known/heartbeat HEAD /.well-known/heartbeat ")
	assert.Contains(t, b.String(), "GET /.well-known/openapi.json HEAD /.well-known/openapi.json ")
	assert.Contains(t, b.String(), "GET /.well-known/swagger HEAD /swagger ")
	assert.Contains(t, b.String(), "GET /.well-known/swagger/{name} HEAD /swagger/{name}")
	assert.NotContains(t, b.String(), "\"NotFoundHandler\":null,\"MethodNotAllowedHandler\":null,\"KeepContext\":false")
}

// TestRouter_CatchAllRoute tests the CatchAllRoute for the requests to not registered endpoints or invalid routes.
func TestRouter_CatchAllRoute(t *testing.T) {
	app := New()

	app.GET("/dummy", func(ctx *Context) (interface{}, error) {
		return nil, nil
	})

	go app.Start()
	time.Sleep(time.Second * 2)

	tests := []struct {
		desc       string
		endpoint   string
		method     string
		statusCode int
	}{
		{"invalid route", "/dummy1", http.MethodGet, http.StatusNotFound},
		{"substring route", "/dumm", http.MethodPost, http.StatusNotFound},
		{"valid route", "/dummy", http.MethodGet, http.StatusOK},
		{"invalid method", "/dummy", http.MethodDelete, http.StatusMethodNotAllowed},
	}

	for i, test := range tests {
		req, _ := http.NewRequest(test.method, "http://localhost:8000"+test.endpoint, http.NoBody)
		client := http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("TEST[%v] %v\nerror while making request, %v", i, test.desc, err)
			continue
		}

		if resp.StatusCode != test.statusCode {
			t.Errorf("TEST[%v] %v\n expected %v, got %v", i, test.desc, test.statusCode, resp.StatusCode)
		}

		_ = resp.Body.Close()
	}
}
