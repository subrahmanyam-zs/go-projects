package cspauth

const (
	encryptionBlockSizeBytes = 16
	cspEncryptionIterations  = 1000
	lenRandomChars           = 6
	appKeyHeader             = "ak"
	clientIDHeader           = "cd"
	authContextHeader        = "ac"
	minAppKeyLen             = 12
)
