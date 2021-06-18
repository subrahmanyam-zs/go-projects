package cspauth

import (
	"bytes"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockHandler struct{}

func (r *MockHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
			"bFhON1NFNVJKeVRadVpDVURrWnU5bks2a1JFdm9JWUFyMWpMckROOFdheHFPbFhIV2xTc3VucW5HS2l6LzNnbXd5bW0zSnJWdXIvdEF6ajVQTmlZejRmMmwwUHBvNk41Q2ZIckhxZURwM1huZlZ0bWo1VVRPZE12QUJLajBJdFRJNXpQMlZUYTdmL2hBclZIbDNCMjJ5Ym9OL0pDVVE5ZVNkYnlzNDBsZy9MUzBQZnYvUzRHRk5zbGRTaW1tbm9LQ2N4R2p5UEhmdWxqZkhHbTRmVTQ3KzJRNjFXL0d5d2pjVFo3NjNlb1h3SEZSa0NzbllzcmJzdExWR2V2NGR3QkJvSUNrMWtvUFpPcjZWcDl4a0VOLytJYXd6TGV4VlExaDcxUXh5Q01qV0hCMXB5dWZnNWNCZWx0ZHBHWmYxQWNFbWV0STUrcGFqOXVjYXZQeGRSU2FaZmtxTStIVFVqNlRVcUNLM0p1cHhFMnMzMVozc2lEZ3ptUi9xWGJnemxoNjAxZDZj",
			"CSP_SHARED_KEY",
			200,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"d1dPQTFUbGZJVzhtcXlRbzNZQ1lSUW5aWGVJK1g3QnF5SEpiUWxBOUo1TkdpaGFCa1hHQy9SVFQ0Y1krNjlSMExPbmpLZHhKaXB4TFFLbXE3dVBFTGcyOHprRGp4ckxqRHRHdmExbUxJUTJBc29CN0NOSm9BWDJHaE12TFpBdDRNOWcwZlpuR0RVUFlhNGMrZlR0eDV5QU9FQWhvNzllbHZudUU1Q0p4WHNLN2g5OFFhNkIzN2o3cWI3Q0dBRlNYeVNpME95elowU3V5MFpnc1d6UjNjMGtWQVYyNU9hc3orVzdxOHhIWkR2dz1hNTQwZTY=",
			"CSP_SHARED_KEY",
			200,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"",
			"CSP_SHARED_KEY",
			200,
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
			"",
			"bFhON1NFNVJKeVRadVpDVURrWnU5bks2a1JFdm9JWUFyMWpMckROOFdheHFPbFhIV2xTc3VucW5HS2l6LzNnbXd5bW0zSnJWdXIvdEF6ajVQTmlZejRmMmwwUHBvNk41Q2ZIckhxZURwM1huZlZ0bWo1VVRPZE12QUJLajBJdFRJNXpQMlZUYTdmL2hBclZIbDNCMjJ5Ym9OL0pDVVE5ZVNkYnlzNDBsZy9MUzBQZnYvUzRHRk5zbGRTaW1tbm9LQ2N4R2p5UEhmdWxqZkhHbTRmVTQ3KzJRNjFXL0d5d2pjVFo3NjNlb1h3SEZSa0NzbllzcmJzdExWR2V2NGR3QkJvSUNrMWtvUFpPcjZWcDl4a0VOLytJYXd6TGV4VlExaDcxUXh5Q01qV0hCMXB5dWZnNWNCZWx0ZHBHWmYxQWNFbWV0STUrcGFqOXVjYXZQeGRSU2FaZmtxTStIVFVqNlRVcUNLM0p1cHhFMnMzMVozc2lEZ3ptUi9xWGJnemxoNjAxZDZj",
			"CSP_SHARED_KEY",
			400,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"bFhON1NFNVJKeVRadVpDVURrWnU5bks2a1JFdm9JWUFyMWpMckROOFdheHFPbFhIV2xTc3VucW5HS2l6LzNnbXd5bW0zSnJWdXIvdEF6ajVQTmlZejRmMmwwUHBvNk41Q2ZIckhxZURwM1huZlZ0bWo1VVRPZE12QUJLajBJdFRJNXpQMlZUYTdmL2hBclZIbDNCMjJ5Ym9OL0pDVVE5ZVNkYnlzNDBsZy9MUzBQZnYvUzRHRk5zbGRTaW1tbm9LQ2N4R2p5UEhmdWxqZkhHbTRmVTQ3KzJRNjFXL0d5d2pjVFo3NjNlb1h3SEZSa0NzbllzcmJzdExWR2V2NGR3QkJvSUNrMWtvUFpPcjZWcDl4a0VOLytJYXd6TGV4VlExaDcxUXh5Q01qV0hCMXB5dWZnNWNCZWx0ZHBHWmYxQWNFbWV0STUrcGFqOXVjYXZQeGRSU2FaZmtxTStIVFVqNlRVcUNLM0p1cHhFMnMzMVozc2lEZ3ptUi9xWGJnemxoNjAxZDZj",
			"",
			200,
			"Dummy body",
		},
		{
			"ak11127983471298348912734",
			"cd1",
			"c29tZSB",
			"CSP_SHARED_KEY",
			200,
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
	}

	for i, tc := range tcs {
		opts := Options{SharedKey: tc.sharedKey}
		body := bytes.NewReader([]byte(tc.body))
		if tc.body == "" {
			body = nil
		}
		req := httptest.NewRequest("GET", "/dummy", body)
		req.Header.Set("ak", tc.appKey)
		req.Header.Set("cd", tc.clientID)
		req.Header.Set("ac", tc.authContext)
		w := httptest.NewRecorder()

		b := new(bytes.Buffer)
		logger := log.NewMockLogger(b)

		handler := CSPAuth(logger, opts)(&MockHandler{})
		handler.ServeHTTP(w, req)

		if w.Code != tc.expCode {
			t.Errorf("TESTCASE[%v]\nexpected code %v,\ngot %v", i, tc.expCode, w.Code)
		}
	}
}

func Test_getBody(t *testing.T) {
	req, _ := http.NewRequest("GET", "/dummy", nil)
	req.Body = nil

	b := getBody(req)
	if len(b) != 0 {
		t.Errorf("expected empty slice, got %v", b)
	}
}

func Test_pkcs7Pad(t *testing.T) {
	tcs := []struct {
		blockSize int
		err       error
	}{
		{0, errInvalidBlockSize},
		{1, errInvalidPKCS7Data},
	}

	for i, tc := range tcs {
		_, err := pkcs7Pad(nil, tc.blockSize)
		if err != tc.err {
			t.Errorf("TESTCASE[%v]:\nexpected %v, got %v", i, tc.err, err)
		}
	}
}

func Test_pkcs7Unpad(t *testing.T) {
	tcs := []struct {
		data      []byte
		blockSize int
		err       error
	}{
		{nil, 0, errInvalidBlockSize},
		{nil, 1, errInvalidPKCS7Data},
		{[]byte{1, 2, 3}, 2, errInvalidPKCS7Padding},
		{[]byte{1, 2, 3, 5}, 2, errInvalidPKCS7Padding},
	}

	for i, tc := range tcs {
		_, err := pkcs7Unpad(tc.data, tc.blockSize)
		if err != tc.err {
			t.Errorf("TESTCASE[%v]:\nexpected %v, got %v", i, tc.err, err)
		}
	}
}
