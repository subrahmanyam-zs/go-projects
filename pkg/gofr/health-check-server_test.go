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

	const port, route = 8086, "/.well-known/health-check"

	srv := healthCheckHandlerServer(ctx, port, route)
	serverURL := "http://localhost:" + strconv.Itoa(port)
	r := httptest.NewRequest(http.MethodGet, serverURL+route, nil)
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Code)
}
