package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockHandlerForHeaderPropagation struct{}

// ServeHTTP is used for testing if the request context has traceId
func (r *MockHandlerForHeaderPropagation) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	channel, _ := req.Context().Value(ZopsmartChannelKey).(string)
	tenant, _ := req.Context().Value(ZopsmartTenantKey).(string)
	authorization := req.Context().Value(AuthorizationHeader).(string)
	body := strings.Join([]string{channel, tenant, authorization}, ",")
	_, _ = w.Write([]byte(body))
}

func TestPropagateHeaders(t *testing.T) {
	handler := PropagateHeaders(&MockHandlerForHeaderPropagation{})
	req := httptest.NewRequest("GET", "/dummy", nil)
	req.Header = map[string][]string{"X-Zopsmart-Tenant": {"zopsmart"}, "X-Zopsmart-Channel": {"WEB"}, "Authorization": {"zop"}}
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Body.String() != "WEB,zopsmart,zop" {
		t.Errorf("propagation of headers through context failed. Got %v\tExpected %v", recorder.Body.String(), "WEB,zopsmart")
	}
}
