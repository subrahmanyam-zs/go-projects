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
		{"ak", authCtx, "CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"ak11127983471298348912734", "YzI5dFpTQjUwMzdkOQ==", "CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"ak11127983471298348912734",
			"d1dPQTFUbGZJVzhtcXlRbzNZQ1lSUW5aWGVJK1g3QnF5SEpiUWxBOUo1TkdpaGFCa1hHQy9SVFQ0Y1krNjlSMExPbmpLZHhKaXB4TFFLbXE3dVBFT" +
				"GcyOHprRGp4ckxqRHRHdmExbUxJUTJBc29CN0NOSm9BWDJHaE12TFpBdDRNOWcwZlpuR0RVUFlhNGMrZlR0eDV5QU9FQWhvNzllbHZudUU1Q0p4WHNLN2g5OFF" +
				"hNkIzN2o3cWI3Q0dBRlNYeVNpME95elowU3V5MFpnc1d6UjNjMGtWQVYyNU9hc3orVzdxOHhIWkR2dz1hNTQwZTY=",
			"CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"ak11127983471298348912734", "", "CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"ak11127983471298348912734", "aA==", "CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"ak11127983471298348912734", "c29tZSB", "CSP_SHARED_KEY", http.StatusUnauthorized, "Dummy body"},
		{"mock-app-key", "Z3lGMUFrMWcyaVE3U0crOExjME5LUlJxS2pCci9JYzdtZUhJZlN0WXJxVkMwUWduSEo4UTI0Ykkva" +
			"VRCVEY4SVdiaDB6RFJqeFlrNUlpNmlQR25NaExtZTRYdkk5cXFBVlNxUDByVHRhK2szd3cxcnpxY1liNURvQzJ6YUF0S1dvcHpWUjRlTExyVnhxTnhJYllhSzd" +
			"0U0hwMUU0NkIxQkk2QzltUzNKbXBHS2NuaDFqSU44L2VUd20zNmp0NDl1cS81anNuMGh0bUFwK2luN1F0RWZ5Yzl5bloyMTE1Njk5ZEdpWEpaNmFadC9GRzFHQjZaZ" +
			"ldicUpxUXdoakF0S0FIUXhINmNGRWVDc0RFRlo5NnVPaVF1N1p3eXUvT3VvT2lyc1pMY3B3cFQ1L1Z6U1JiTjhvR0tVOEV5UTF6cFRiSzF3Qit4WUF4L09JK3FhUlVk" +
			"Vzcrb2lQUVBuNEV5UkN6eGhDZEw1SWlJPTkyYjJhYg==",
			"mock-shared-key", http.StatusUnauthorized, "Dummy-body1"},
		{"mock-app-key", "cjV4aGdZdy9nMlVnbDVpQ1VDcWZzaE8ySk5YZ2MrR010anpwRnArRStNVWtQUHA1ZVVkNm1leE5iL2trMnlhNHFYL2wvQ1Mxdk" +
			"EzNG9zRTZLWlhiZ25Fd3p4R3JRUWs3azYvQ2h1aEN3YjVpY1RCVVd3Ty9xT2VEMWlrSkMyUVVYZDhsdTNtbTNGK0dIZlk2ZlVSRE40TExCSVAzYm0zOU9RbW4vWWx4Z1ovQ" +
			"TZidSs2N0RFelhaVVV1SEpXenppVUI3cU5wUWg4U0lqVVlVRXJjcnZUTHdNdVJ4ZEV3aXViblJEYS9SRG5Wc01kTjIzY3ZEbHJXQU1OR2wzUHFaa2hlaWxIcVlraFZ1YWZZbk" +
			"U2Q09EUTVVK1c5aU5QRUk2RCtIN1FsMGtGSXl1OFBjaUhzZnlqWjdDUU04SDRhUmE0U245V24yN1RTQzJEU0tOc1BCc3hKZVZob1dDTGk0T2pIYktTT3dpc0Y0PWJiMGY0Ng==",
			"mock-shared-key", http.StatusUnauthorized, "dummy-body"},
		{"mock-app-key", "Z3lGMUFrMWcyaVE3U0crOExjME5LUlJxS2pCci9JYzdtZUhJZlN0WXJxVkMwUWduSEo4UTI0Ykkva" +
			"VRCVEY4SVdiaDB6RFJqeFlrNUlpNmlQR25NaExtZTRYdkk5cXFBVlNxUDByVHRhK2szd3cxcnpxY1liNURvQzJ6YUF0S1dvcHpWUjRlTExyVnhxTnhJYllhSzd" +
			"0U0hwMUU0NkIxQkk2QzltUzNKbXBHS2NuaDFqSU44L2VUd20zNmp0NDl1cS81anNuMGh0bUFwK2luN1F0RWZ5Yzl5bloyMTE1Njk5ZEdpWEpaNmFadC9GRzFHQjZaZ" +
			"ldicUpxUXdoakF0S0FIUXhINmNGRWVDc0RFRlo5NnVPaVF1N1p3eXUvT3VvT2lyc1pMY3B3cFQ1L1Z6U1JiTjhvR0tVOEV5UTF6cFRiSzF3Qit4WUF4L09JK3FhUlVk" +
			"Vzcrb2lQUVBuNEV5UkN6eGhDZEw1SWlJPTkyYjJhYg==",
			"mock-shared-key", http.StatusOK, "Dummy-body"},
		{"mock-app-key", "QnlSSm53ZWJta1pITDYrU3JWaXkwSkJPTStmcnFDeDFDbWNYSUVWS01PMDBQSTlmT2kwSWFHTHB3Z3BVOWFlS21OVERLM29MT0" +
			"F5aFVyTUUxcGVic0cxbzgzSkMzUC9ZUlA1MkFTR2ljY1BGa0NuVmNRdkxDRmZ4d2lqTDVDQjM2Rks5VWZqbGFyRnp5ckY0TTBHUkt5cnB0Mm1IaGN0UVZ4WVI3bWdYVThxN" +
			"ndzbjZQdWpEV2FGaWNIYkkxRCtHeEpvdWNjaGFEUnMxS1B5UDRvd09hUmFuT1FGajFmSGFQcE1yVmNwYlQxc0JWeGJEOGVGNWJSRTdrMWhwbGlEWFV2Q25kWXpCblJJeXJnSC" +
			"9lbGg5WktJcE1CcEEyd2xaUjBJZUF1RDB6VlUrK09mdkNvTlpxMy9sWC9yc1pSUEo3Z0JIdTJlMk1XaXhZYXRJR0ZRY01jSVViL2o2YlJLTEZqY3NRaG9rT21mQ3ZsZFMrR2J" +
			"ldlh3M2RRUDBLTU84N1E5SEloMHlLRThDUjJWampiU3BtQT09ZWRlYTkz",
			"mock-shared-key", http.StatusOK, ""},
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
