package service

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

func Test_getAuthContext(t *testing.T) {
	tcs := []struct {
		appKey    string
		clientID  string
		sharedKey string
		body      string
	}{
		{"ankling123jerkins4junked", "cd1", "CSP_SHARED_KEY", "Dummy body"},
	}

	for i, tc := range tcs {
		opts := &CSPOption{
			SharedKey: tc.sharedKey,
			AppKey:    tc.appKey,
			AppID:     tc.clientID,
		}

		logger := log.NewMockLogger(io.Discard)
		csp, _ := New(logger, opts)
		body := bytes.NewReader([]byte(tc.body))
		req, _ := http.NewRequest("POST", "/dummy", body)
		ac := csp.getAuthContext(http.MethodPost, req.Body)

		if ac == "" {
			t.Errorf("TESTCASE[%v] Expected to be get verified auth context", i)
		}
	}
}

func Test_bodyHash(t *testing.T) {
	tcs := []struct {
		body    io.Reader
		expHash string
	}{
		{bytes.NewReader([]byte("Hello")), "185F8DB32271FE25F561A6FC938B2E264306EC304EDA518007D1764826381969"},
		{nil, ""},
	}

	for i, tc := range tcs {
		hash := getBodyHash(tc.body)

		if hash != tc.expHash {
			t.Errorf("TESTCASE[%v] Expected hash %v, got %v", i, tc.expHash, hash)
		}
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
		{"", "cd1", "CSP_SHARED_KEY", ErrEmptyAppKey},
		{"ankling123jerkins4junked", "cd1", "", ErrEmptySharedKey},
		{"ankling123jerkins4junked", "", "CSP_SHARED_KEY", ErrEmptyAppID},
	}

	for i, tc := range tcs {
		opts := &CSPOption{
			SharedKey: tc.sharedKey,
			AppKey:    tc.appKey,
			AppID:     tc.clientID,
		}

		logger := log.NewMockLogger(io.Discard)
		_, err := New(logger, opts)

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
