package cspauth

import (
	"net/http"
	"testing"
)

func Test_getBody(t *testing.T) {
	req, _ := http.NewRequest("GET", "/dummy", http.NoBody)
	req.Body = nil

	b := GetBodyHash(req)
	if b != "" {
		t.Errorf("expected empty string, got %v", b)
	}
}

func Test_pkcs7Unpad(t *testing.T) {
	tcs := []struct {
		data      []byte
		blockSize int
		err       error
	}{
		{nil, 0, ErrInvalidBlockSize},
		{nil, 1, ErrInvalidPKCS7Data},
		{[]byte{1, 2, 3}, 2, ErrInvalidPKCS7Padding},
		{[]byte{1, 2, 3, 5}, 2, ErrInvalidPKCS7Padding},
	}

	for i, tc := range tcs {
		_, err := pkcs7Unpad(tc.data, tc.blockSize)
		if err != tc.err {
			t.Errorf("TESTCASE[%v]:\nexpected %v, got %v", i, tc.err, err)
		}
	}
}
