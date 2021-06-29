package cspauth

const (
	securityType             = "1"
	securityVersion          = "V1"
	encryptionBlockSizeBytes = 16
	cspEncryptionIterations  = 1000
	lenRandomChars           = 6
	appKeyHeader             = "ak"
	clientIDHeader           = "cd"
	authContextHeader        = "ac"
	securityTypeHeader       = "st"
	securityVersionHeader    = "sv"
	minAppKeyLen             = 12
)
