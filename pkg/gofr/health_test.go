package gofr

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestHeartBeatHandler(t *testing.T) {
	c := &Context{}
	tests := []struct {
		name    string
		c       *Context
		want    interface{}
		wantErr bool
	}{
		{"test1", c, map[string]string{"status": "UP"}, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := HeartBeatHandler(tt.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("heartBeatHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("heartBeatHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_HeartBeatIntegration(t *testing.T) {
	k := New()
	k.Server.HTTP.Port = 3339
	http.DefaultServeMux = new(http.ServeMux)

	go k.Start()
	time.Sleep(3 * time.Second)

	client := http.Client{}
	url := "http://localhost:3339/.well-known/heartbeat"
	req, _ := request.NewMock(http.MethodGet, url, nil)

	resp, err := client.Do(req)

	if err != nil {
		t.Errorf("got error %s", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("want status code 200 got= %T", resp.StatusCode)
	}

	resp.Body.Close()

	url = "http://localhost:3339/.well-known/health-check"
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	resp, err = client.Do(req)

	if err != nil {
		t.Errorf("got error %s", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("want status code 200 got= %v", resp.StatusCode)
	}

	resp.Body.Close()
}

func Test_server_HeartCheck(t *testing.T) {
	s := New()
	s.Server.HTTP.Port = 3340
	s.Server.HTTPS.CertificateFile = ""

	tests := []struct {
		name string
		want string
	}{
		{"test1", "GET /.well-known/health-check"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			http.DefaultServeMux = new(http.ServeMux)
			go s.Start()
			time.Sleep(3 * time.Second)
			got := fmt.Sprintf("%s", s.Server.Router)

			if reflect.DeepEqual(got, tt.want) {
				t.Errorf(" got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_finalStatus(t *testing.T) {
	upCount := 5
	downCount := 2

	status := finalStatus(upCount, downCount)
	if status != pkg.StatusDegraded {
		t.Errorf("Incorrect status: Expected: %v  Got: %v", pkg.StatusDegraded, status)
	}
}

func Test_GetAppDetails(t *testing.T) {
	cfg := config.NewGoDotEnvProvider(log.NewLogger(), "../../configs")

	testCases := []struct {
		config   Config
		expected types.AppDetails
	}{
		{
			config:   cfg,
			expected: types.AppDetails{Name: cfg.Get("APP_NAME"), Version: cfg.Get("APP_VERSION"), Framework: pkg.Framework},
		},
		{
			config:   &config.MockConfig{Data: map[string]string{"APP_NAME": "sample-app", "APP_VERSION": "test-version"}},
			expected: types.AppDetails{Name: "sample-app", Version: "test-version", Framework: pkg.Framework},
		},
		{
			config:   &config.MockConfig{},
			expected: types.AppDetails{Name: pkg.DefaultAppName, Version: pkg.DefaultAppVersion, Framework: pkg.Framework},
		},
	}

	for i, testCase := range testCases {
		got := getAppDetails(testCase.config)
		assert.Equal(t, testCase.expected, got, i)
	}
}

func Test_HealthCheckHandler(t *testing.T) {
	k := New()
	ctx := NewContext(nil, nil, k)

	healthCheckResponse, err := HealthHandler(ctx)
	if err != nil {
		t.Error(err)
	}

	m, ok := healthCheckResponse.(map[string]interface{})
	if !ok {
		t.Errorf("expected type map[string]interface{} got %T", m)
	}

	// details should not be nil
	if _, ok = m["details"]; !ok {
		t.Errorf("details should not be nil")
	}

	// status should be UP
	if v, ok := m["status"]; ok {
		if v != pkg.StatusUp {
			t.Errorf("status should be UP")
		}
	} else {
		t.Errorf("status should not be nil")
	}
}
