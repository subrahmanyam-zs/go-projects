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
	tcs := []struct {
		appKey      string
		clientID    string
		authContext string
		sharedKey   string
		expCode     int
		body        string
	}{
		{
			"ankling123jerkins4junked",
			"cd1",
			"Y1YyM0ptSDJqTmRkNlRCNFArbkp5ck9LaWlXc2NCME9WWmNNUXZ6ZFVRR1VoYnhFRmdvNWpsd3daSjFEWDdrMnJuM1d6Yk9Ic045MTVYVWFUdEQ1V1d5RHZyZ2phOWU5aUVPcXZsM1JUT1lQanFqVFFVZ0tKT1ZqK0VZRDhMYnpTenZ5dFNzbmVKS1hZdW5JdVBBYU8ySDNqY2toaHJFcUxGMEJhajJZN0Y2b2VLRUc0bUoyMGdrazBybDY5NVE1RlZDTW1QdzNkdnZ1TkRTSjlMZmNmSm5DZzFWNnRybm52dG1MQlloSi9LSEZydW8rRm9SOEVNM0Y3Q1pDZUFMQVVoL1FYeWR1c1FoV0wxcm9xMVd0SDdjV0FOZU0xSmtoNnVXM3dYRTI4NjlRb1o3cmFtck5YaW9KcUpSczM5cnFXVXlrRHp2T2pGTWV2NHFiL2U2VzIydHNwV04xa3VkY0t2OXNUcFlsUUJZPWJhMGZmOQ==",
			"CSP_SHARED_KEY",
			200,
			"Dummy body",
		},
		{
			"ankling123jerkins4junked",
			"cd1",
			"bFhON1NFNVJKeVRadVpDVURrWnU5bks2a1JFdm9JWUFyMWpMckROOFdheHFPbFhIV2xTc3VucW5HS2l6LzNnbXd5bW0zSnJWdXIvdEF6ajVQTmlZejRmMmwwUHBvNk41Q2ZIckhxZURwM1huZlZ0bWo1VVRPZE12QUJLajBJdFRJNXpQMlZUYTdmL2hBclZIbDNCMjJ5Ym9OL0pDVVE5ZVNkYnlzNDBsZy9MUzBQZnYvUzRHRk5zbGRTaW1tbm9LQ2N4R2p5UEhmdWxqZkhHbTRmVTQ3KzJRNjFXL0d5d2pjVFo3NjNlb1h3SEZSa0NzbllzcmJzdExWR2V2NGR3QkJvSUNrMWtvUFpPcjZWcDl4a0VOLytJYXd6TGV4VlExaDcxUXh5Q01qV0hCMXB5dWZnNWNCZWx0ZHBHWmYxQWNFbWV0STUrcGFqOXVjYXZQeGRSU2FaZmtxTStIVFVqNlRVcUNLM0p1cHhFMnMzMVozc2lEZ3ptUi9xWGJnemxoNjAxZDZj",
			"CSP_SHARED_KEY",
			403,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"d1dPQTFUbGZJVzhtcXlRbzNZQ1lSUW5aWGVJK1g3QnF5SEpiUWxBOUo1TkdpaGFCa1hHQy9SVFQ0Y1krNjlSMExPbmpLZHhKaXB4TFFLbXE3dVBFTGcyOHprRGp4ckxqRHRHdmExbUxJUTJBc29CN0NOSm9BWDJHaE12TFpBdDRNOWcwZlpuR0RVUFlhNGMrZlR0eDV5QU9FQWhvNzllbHZudUU1Q0p4WHNLN2g5OFFhNkIzN2o3cWI3Q0dBRlNYeVNpME95elowU3V5MFpnc1d6UjNjMGtWQVYyNU9hc3orVzdxOHhIWkR2dz1hNTQwZTY=",
			"CSP_SHARED_KEY",
			403,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"",
			"CSP_SHARED_KEY",
			403,
			"Dummy body",
		},
		{
			"ak",
			"cd1",
			"bFhON1NFNVJKeVRadVpDVURrWnU5bks2a1JFdm9JWUFyMWpMckROOFdheHFPbFhIV2xTc3VucW5HS2l6LzNnbXd5bW0zSnJWdXIvdEF6ajVQTmlZejRmMmwwUHBvNk41Q2ZIckhxZURwM1huZlZ0bWo1VVRPZE12QUJLajBJdFRJNXpQMlZUYTdmL2hBclZIbDNCMjJ5Ym9OL0pDVVE5ZVNkYnlzNDBsZy9MUzBQZnYvUzRHRk5zbGRTaW1tbm9LQ2N4R2p5UEhmdWxqZkhHbTRmVTQ3KzJRNjFXL0d5d2pjVFo3NjNlb1h3SEZSa0NzbllzcmJzdExWR2V2NGR3QkJvSUNrMWtvUFpPcjZWcDl4a0VOLytJYXd6TGV4VlExaDcxUXh5Q01qV0hCMXB5dWZnNWNCZWx0ZHBHWmYxQWNFbWV0STUrcGFqOXVjYXZQeGRSU2FaZmtxTStIVFVqNlRVcUNLM0p1cHhFMnMzMVozc2lEZ3ptUi9xWGJnemxoNjAxZDZj",
			"CSP_SHARED_KEY",
			400,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"c29tZSB",
			"CSP_SHARED_KEY",
			403,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"bFhON1NFNVJKeVRadVpDVURrWnU5bks2a1JFdm9JWUFyMWpMckROOFdheHFPbFhIV2xTc3VucW5HS2l6LzNnbXd5bW0zSnJWdXIvdEF6ajVQTmlZejRmMmwwUHBvNk41Q2ZIckhxZURwM1huZlZ0bWo1VVRPZE12QUJLajBJdFRJNXpQMlZUYTdmL2hBclZIbDNCMjJ5Ym9OL0pDVVE5ZVNkYnlzNDBsZy9MUzBQZnYvUzRHRk5zbGRTaW1tbm9LQ2N4R2p5UEhmdWxqZkhHbTRmVTQ3KzJRNjFXL0d5d2pjVFo3NjNlb1h3SEZSa0NzbllzcmJzdExWR2V2NGR3QkJvSUNrMWtvUFpPcjZWcDl4a0VOLytJYXd6TGV4VlExaDcxUXh5Q01qV0hCMXB5dWZnNWNCZWx0ZHBHWmYxQWNFbWV0STUrcGFqOXVjYXZQeGRSU2FaZmtxTStIVFVqNlRVcUNLM0p1cHhFMnMzMVozc2lEZ3ptUi9xWGJnemxoNjAxZDZj",
			"CSP_SHARED_KEY",
			200,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"YzI5dFpTQjUwMzdkOQ==",
			"CSP_SHARED_KEY",
			403,
			"Dummy body",
		},
	}

	for i, tc := range tcs {
		opts := Options{SharedKey: tc.sharedKey}
		body := bytes.NewReader([]byte(tc.body))
		req := httptest.NewRequest("GET", "/dummy", body)
		req.Header.Set("ak", tc.appKey)
		req.Header.Set("cd", tc.clientID)
		req.Header.Set("ac", tc.authContext)

		w := httptest.NewRecorder()

		logger := log.NewMockLogger(io.Discard)

		handler := CSPAuth(logger, &opts)(&MockHandler{})
		handler.ServeHTTP(w, req)

		if w.Code != tc.expCode {
			t.Errorf("TESTCASE[%v]\nexpected code %v,\ngot %v", i, tc.expCode, w.Code)
		}
	}
}

func Test_Set(t *testing.T) {
	tcs := []struct {
		appKey    string
		clientID  string
		sharedKey string
		body      string
	}{
		{
			"ankling123jerkins4junked",
			"cd1",
			"CSP_SHARED_KEY",
			"Dummy body",
		},
	}

	for i, tc := range tcs {
		opts := &Options{
			SharedKey: tc.sharedKey,
			AppKey:    tc.appKey,
			AppID:     tc.clientID,
		}

		logger := log.NewMockLogger(io.Discard)
		csp, _ := New(logger, opts)
		body := bytes.NewReader([]byte(tc.body))
		req, _ := http.NewRequest("POST", "/dummy", body)
		csp.Set(req)

		if !csp.Verify(logger, req) {
			t.Errorf("TESTCASE[%v] Expected to be get verified auth context", i)
		}
	}
}
