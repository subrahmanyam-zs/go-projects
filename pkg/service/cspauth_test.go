package service

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/middleware"

	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware/cspauth"
)

func Test_getAuthContext(t *testing.T) {
	opts := &CSPOption{
		SharedKey: "CSP_SHARED_KEY",
		AppKey:    "ankling123jerkins4junked",
	}

	logger := log.NewMockLogger(io.Discard)
	csp, _ := NewCSP(logger, opts)
	body := bytes.NewReader([]byte(`{"foo":"bar"}`))
	req, _ := http.NewRequest(http.MethodPost, "/dummy", body)
	ac := csp.getAuthContext(req)

	if ac == "" {
		t.Errorf("Failed to generate authcontext")
	}
}

func Test_validate(t *testing.T) {
	tcs := []struct {
		appKey    string
		clientID  string
		sharedKey string
		err       error
	}{
		{"ankling123jerkins4junked", "cd1", "CSP_SHARED_KEY", nil},
		{"", "cd1", "CSP_SHARED_KEY", middleware.ErrInvalidAppKey},
		{"ankling123jerkins4junked", "cd1", "", ErrEmptySharedKey},
	}

	for i, tc := range tcs {
		opts := &CSPOption{
			SharedKey: tc.sharedKey,
			AppKey:    tc.appKey,
			ClientID:  tc.clientID,
		}

		logger := log.NewMockLogger(io.Discard)
		_, err := NewCSP(logger, opts)

		if err != tc.err {
			t.Errorf("TESTCASE[%v] Expected error %v, got %v", i, tc.err, err)
		}
	}
}

func Test_pkcs7Pad(t *testing.T) {
	tcs := []struct {
		blockSize int
		err       error
	}{
		{0, cspauth.ErrInvalidBlockSize},
		{1, cspauth.ErrInvalidPKCS7Data},
	}

	for i, tc := range tcs {
		_, err := pkcs7Pad(nil, tc.blockSize)
		if err != tc.err {
			t.Errorf("TESTCASE[%v]:\nexpected %v, got %v", i, tc.err, err)
		}
	}
}
