package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHandlerForNewRelic struct{}

func (m *mockHandlerForNewRelic) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

func TestNewRelic(t *testing.T) {
	handler := NewRelic("gofr", "6378b0a5bf929e7eb36d480d4e3cd914b74eNRAL")(&mockHandlerForNewRelic{})
	req, _ := http.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	newRelicTxn := req.Context().Value(newRelicTxnKey)

	if newRelicTxn == nil {
		t.Error("NewRelicTxn not injected into the request")
	}
}
