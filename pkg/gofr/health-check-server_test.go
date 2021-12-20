package gofr

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandlerServer(t *testing.T) {
	app := NewCMD()

	srv := healthCheckHandlerServer(app.cmd)
	serverURL := "http://localhost:" + strconv.Itoa(defaultMetricsPort)
	r := httptest.NewRequest(http.MethodGet, serverURL+"/.well-known/health-check", nil)
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Code)
}
