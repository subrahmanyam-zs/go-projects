package gofr

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHealthCheckHandlerServer(t *testing.T) {
	k := New()
	ctx := NewContext(nil, nil, k)

	srv := healthCheckHandlerServer(ctx, defaultMetricsPort)
	serverURL := "http://localhost:" + strconv.Itoa(defaultMetricsPort)
	r := httptest.NewRequest(http.MethodGet, serverURL+defaultHealthCheckRoute, nil)
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Code)
}
