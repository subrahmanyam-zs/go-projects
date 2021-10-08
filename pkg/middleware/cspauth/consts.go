package cspauth

const (
	IVLength                = 16
	EncryptionKeyLen        = 32
	cspEncryptionIterations = 1000
	lenRandomChars          = 6
	minLenAppKey            = 12
	appKeyHeader            = "ak"
	authContextHeader       = "ac"
	securityVersionHeader   = "sv"
	securityTypeHeader      = "st"
	cspSecurityVersion      = "v1"
	cspSecurityType         = "1"
)
