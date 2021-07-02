package cspauth

import "developer.zopsmart.com/go/gofr/pkg/errors"

var (
	// ErrEmptyAppKey is raised when app key is is not more than 12 bytes
	ErrEmptyAppKey = errors.Error("app key should be more than 12 bytes")

	ErrInvalidBlockSize    = errors.Error("invalid block size")
	ErrInvalidPKCS7Data    = errors.Error("invalid PKCS7 data (empty or not padded)")
	ErrInvalidPKCS7Padding = errors.Error("invalid padding on input")
)
