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
		authContext string
		sharedKey   string
		expCode     int
		body        string
	}{
		{"ak", authCtx, "CSP_SHARED_KEY", http.StatusBadRequest, "Dummy body"},
		{"ak11127983471298348912734", "YzI5dFpTQjUwMzdkOQ==", "CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"ak11127983471298348912734",
			"d1dPQTFUbGZJVzhtcXlRbzNZQ1lSUW5aWGVJK1g3QnF5SEpiUWxBOUo1TkdpaGFCa1hHQy9SVFQ0Y1krNjlSMExPbmpLZHhKaXB4TFFLbXE3dVBFT" +
				"GcyOHprRGp4ckxqRHRHdmExbUxJUTJBc29CN0NOSm9BWDJHaE12TFpBdDRNOWcwZlpuR0RVUFlhNGMrZlR0eDV5QU9FQWhvNzllbHZudUU1Q0p4WHNLN2g5OFF" +
				"hNkIzN2o3cWI3Q0dBRlNYeVNpME95elowU3V5MFpnc1d6UjNjMGtWQVYyNU9hc3orVzdxOHhIWkR2dz1hNTQwZTY=",
			"CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"ak11127983471298348912734", "", "CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"ak11127983471298348912734", "c29tZSB", "CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"mock-app-key", "Z3lGMUFrMWcyaVE3U0crOExjME5LUlJxS2pCci9JYzdtZUhJZlN0WXJxVkMwUWduSEo4UTI0Ykkva" +
			"VRCVEY4SVdiaDB6RFJqeFlrNUlpNmlQR25NaExtZTRYdkk5cXFBVlNxUDByVHRhK2szd3cxcnpxY1liNURvQzJ6YUF0S1dvcHpWUjRlTExyVnhxTnhJYllhSzd" +
			"0U0hwMUU0NkIxQkk2QzltUzNKbXBHS2NuaDFqSU44L2VUd20zNmp0NDl1cS81anNuMGh0bUFwK2luN1F0RWZ5Yzl5bloyMTE1Njk5ZEdpWEpaNmFadC9GRzFHQjZaZ" +
			"ldicUpxUXdoakF0S0FIUXhINmNGRWVDc0RFRlo5NnVPaVF1N1p3eXUvT3VvT2lyc1pMY3B3cFQ1L1Z6U1JiTjhvR0tVOEV5UTF6cFRiSzF3Qit4WUF4L09JK3FhUlVk" +
			"Vzcrb2lQUVBuNEV5UkN6eGhDZEw1SWlJPTkyYjJhYg==",
			"mock-shared-key", http.StatusUnauthorized, "Dummy-body1"},
		{"mock-app-key", "cjV4aGdZdy9nMlVnbDVpQ1VDcWZzaE8ySk5YZ2MrR010anpwRnArRStNWGhvWEppVmRzVGpMTmliVlJPSFBnSDcyU1BRT3NyWlR" +
			"FdnNNdHlRTDE4NGQyZ3dONUZncE5GcXBBTGg5OVFaYjl1QllXYmxRR1NjMXFobk1VM1JUREZ3TVNSV3UwWm00bHJVVUorQW51OC9vbkNvbm9JTWROSnJmaktoWFRrSDQ5NE9" +
			"EQ3VHdGFyK0xlRWNVN1NISUZ3N2tnVzczSU5uS2tBN3E2VTY5bkVyMDFzT21RNVVmdHlwMDRaNjhVcDE0UkJBNkpVd3A3SnRXbmNpcTRhazBPQzBpdURxMjZFS" +
			"G5rYW5yRCs1VViaU8yelVqSmo2eGtYMUhwOGFoM3JCV284V1UwejhZbHJtMnhWTWVBR1pEbVN0ZmlSY2NzZTF0RElwWlFxemh2bmlkRnQ0em00elJi" +
			"Z1BKV0RwbTU1L3ZsYzdheVJPekNUS1ptSC9jMjdSMnFXRFh6ZFFpenhPd0JSeGJBaHVCVEUvdz09NGUwYTI2",
			"mock-shared-key", http.StatusUnauthorized, "dummy-body"},
		{"mock-app-key", "Z3lGMUFrMWcyaVE3U0crOExjME5LUlJxS2pCci9JYzdtZUhJZlN0WXJxVkMwUWduSEo4UTI0Ykkva" +
			"VRCVEY4SVdiaDB6RFJqeFlrNUlpNmlQR25NaExtZTRYdkk5cXFBVlNxUDByVHRhK2szd3cxcnpxY1liNURvQzJ6YUF0S1dvcHpWUjRlTExyVnhxTnhJYllhSzd" +
			"0U0hwMUU0NkIxQkk2QzltUzNKbXBHS2NuaDFqSU44L2VUd20zNmp0NDl1cS81anNuMGh0bUFwK2luN1F0RWZ5Yzl5bloyMTE1Njk5ZEdpWEpaNmFadC9GRzFHQjZaZ" +
			"ldicUpxUXdoakF0S0FIUXhINmNGRWVDc0RFRlo5NnVPaVF1N1p3eXUvT3VvT2lyc1pMY3B3cFQ1L1Z6U1JiTjhvR0tVOEV5UTF6cFRiSzF3Qit4WUF4L09JK3FhUlVk" +
			"Vzcrb2lQUVBuNEV5UkN6eGhDZEw1SWlJPTkyYjJhYg==",
			"mock-shared-key", http.StatusOK, "Dummy-body"},
	}

	for i, tc := range tcs {
		body := bytes.NewReader([]byte(tc.body))
		req := httptest.NewRequest(http.MethodPost, "/dummy", body)
		req.Header.Set("ak", tc.appKey)
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
