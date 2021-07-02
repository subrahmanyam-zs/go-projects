package cspauth

import "developer.zopsmart.com/go/gofr/pkg/errors"

var (
	// ErrEmptyAppKey is raised when app key is is not more than 12 bytes
	ErrEmptyAppKey = errors.Error("app key should be more than 12 bytes")

	errInvalidBlockSize    = errors.Error("invalid block size")
	errInvalidPKCS7Data    = errors.Error("invalid PKCS7 data (empty or not padded)")
	errInvalidPKCS7Padding = errors.Error("invalid padding on input")
)
