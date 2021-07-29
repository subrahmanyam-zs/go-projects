package cspauth

import "developer.zopsmart.com/go/gofr/pkg/errors"

var (
	// ErrInvalidBlockSize if aes cipher of not block size 16,24,32,
	ErrInvalidBlockSize    = errors.Error("invalid block size")
	ErrInvalidPKCS7Data    = errors.Error("invalid PKCS7 data (empty or not padded)")
	ErrInvalidPKCS7Padding = errors.Error("invalid padding on input")
)
