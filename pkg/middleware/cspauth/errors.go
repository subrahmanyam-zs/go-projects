package cspauth

import "developer.zopsmart.com/go/gofr/pkg/errors"

var (
	// ErrEmptySharedKey is raised when shared key is empty
	ErrEmptySharedKey = errors.Error("shared key cannot be empty")
	// ErrEmptyAppKey is raised when app key is is not more than 12 bytes
	ErrEmptyAppKey = errors.Error("app key should be more than 12 bytes for successful key generation")
	// ErrEmptyAppID is raised when app id is empty
	ErrEmptyAppID = errors.Error("app id cannot be empty")

	errInvalidBlockSize    = errors.Error("invalid blocksize")
	errInvalidPKCS7Data    = errors.Error("invalid PKCS7 data (empty or not padded)")
	errInvalidPKCS7Padding = errors.Error("invalid padding on input")
)
