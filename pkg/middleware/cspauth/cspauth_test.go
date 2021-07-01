package cspauth

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

type MockHandler struct{}

func (r *MockHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestCSPAuth(t *testing.T) {
	authCtx := "bFhON1NFNVJKeVRadVpDVURrWnU5bks2a1JFdm9JWUFyMWpMckROOFdheHFPbFhIV2xTc3VucW5HS2l6LzNnbXd5bW0zSnJWdXI" +
		"vdEF6ajVQTmlZejRmMmwwUHBvNk41Q2ZIckhxZURwM1huZlZ0bWo1VVRPZE12QUJLajBJdFRJNXpQMlZUYTdmL2hBclZIbDNCMjJ5Ym9OL0p" +
		"DVVE5ZVNkYnlzNDBsZy9MUzBQZnYvUzRHRk5zbGRTaW1tbm9LQ2N4R2p5UEhmdWxqZkhHbTRmVTQ3KzJRNjFXL0d5d2pjVFo3NjNlb1h3SEZSa" +
		"0NzbllzcmJzdExWR2V2NGR3QkJvSUNrMWtvUFpPcjZWcDl4a0VOLytJYXd6TGV4VlExaDcxUXh5Q01qV0hCMXB5dWZnNWNCZWx0ZHBHWmYxQWNFb" +
		"WV0STUrcGFqOXVjYXZQeGRSU2FaZmtxTStIVFVqNlRVcUNLM0p1cHhFMnMzMVozc2lEZ3ptUi9xWGJnemxoNjAxZDZj"

	tcs := []struct {
		appKey      string
		clientID    string
		authContext string
		sharedKey   string
		expCode     int
		body        string
	}{
		{"ak", "cd1", authCtx, "CSP_SHARED_KEY", http.StatusBadRequest, "Dummy body"},
		{"ak11127983471298348912734", "cd1", "YzI5dFpTQjUwMzdkOQ==", "", http.StatusBadRequest, "Dummy body"},
		{"ak11127983471298348912734", "", "YzI5dFpTQjUwMzdkOQ==", "CSP_SHARED_KEY", http.StatusBadRequest, "Dummy body"},
		{"ak11127983471298348912734", "cd1", "YzI5dFpTQjUwMzdkOQ==", "CSP_SHARED_KEY", http.StatusForbidden, "Dummy body"},
		{"ankling123jerkins4junked", "cd1", authCtx, "CSP_SHARED_KEY", http.StatusForbidden, "Dummy body"},
		{"ak11127983471298348912734", "cd1",
			"d1dPQTFUbGZJVzhtcXlRbzNZQ1lSUW5aWGVJK1g3QnF5SEpiUWxBOUo1TkdpaGFCa1hHQy9SVFQ0Y1krNjlSMExPbmpLZHhKaXB4TFFLbXE3dVBFT" +
				"GcyOHprRGp4ckxqRHRHdmExbUxJUTJBc29CN0NOSm9BWDJHaE12TFpBdDRNOWcwZlpuR0RVUFlhNGMrZlR0eDV5QU9FQWhvNzllbHZudUU1Q0p4WHNLN2g5OFF" +
				"hNkIzN2o3cWI3Q0dBRlNYeVNpME95elowU3V5MFpnc1d6UjNjMGtWQVYyNU9hc3orVzdxOHhIWkR2dz1hNTQwZTY=",
			"CSP_SHARED_KEY", http.StatusForbidden, "Dummy body"},
		{"ak11127983471298348912734", "cd1", "", "CSP_SHARED_KEY", http.StatusForbidden, "Dummy body"},
		{"ak11127983471298348912734", "cd1", "c29tZSB", "CSP_SHARED_KEY", http.StatusForbidden, "Dummy body"},
		{"ak11127983471298348912734", "cd1", authCtx, "CSP_SHARED_KEY", http.StatusForbidden, "Dummy body1"},
		{"ankling123jerkins4junked", "cd1",
			"Y1YyM0ptSDJqTmRkNlRCNFArbkp5ck9LaWlXc2NCME9WWmNNUXZ6ZFVRR1VoYnhFRmdvNWpsd3daSjFEWDdrMnJuM1d6Yk9Ic045MTV" +
				"YVWFUdEQ1V1d5RHZyZ2phOWU5aUVPcXZsM1JUT1lQanFqVFFVZ0tKT1ZqK0VZRDhMYnpTenZ5dFNzbmVKS1hZdW5JdVBBYU8ySDNqY2toaHJFcUxG" +
				"MEJhajJZN0Y2b2VLRUc0bUoyMGdrazBybDY5NVE1RlZDTW1QdzNkdnZ1TkRTSjlMZmNmSm5DZzFWNnRybm52dG1MQlloSi9LSEZydW8rRm9SOEVNM0Y3Q" +
				"1pDZUFMQVVoL1FYeWR1c1FoV0wxcm9xMVd0SDdjV0FOZU0xSmtoNnVXM3dYRTI4NjlRb1o3cmFtck5YaW9KcUpSczM5cnFXVXlrRHp2T2pGTWV2NHFiL2U2Vz" +
				"IydHNwV04xa3VkY0t2OXNUcFlsUUJZPWJhMGZmOQ==",
			"CSP_SHARED_KEY", http.StatusOK, "Dummy body"},
		{"ak11127983471298348912734", "cd1", authCtx, "CSP_SHARED_KEY", http.StatusOK, "Dummy body"},
	}

	for i, tc := range tcs {
		body := bytes.NewReader([]byte(tc.body))
		req := httptest.NewRequest(http.MethodPost, "/dummy", body)
		req.Header.Set("ak", tc.appKey)
		req.Header.Set("cd", tc.clientID)
		req.Header.Set("ac", tc.authContext)

		w := httptest.NewRecorder()

		logger := log.NewMockLogger(io.Discard)

		handler := CSPAuth(logger, tc.sharedKey)(&MockHandler{})
		handler.ServeHTTP(w, req)

		if w.Code != tc.expCode {
			t.Errorf("TESTCASE[%v]\nexpected code %v,\ngot %v", i, tc.expCode, w.Code)
		}
	}
}

func Test_ExemptPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/.well-known/health-check", nil)

	w := httptest.NewRecorder()

	logger := log.NewMockLogger(io.Discard)

	handler := CSPAuth(logger, "")(&MockHandler{})
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected code %v,\nGot %v", http.StatusOK, w.Code)
	}
}