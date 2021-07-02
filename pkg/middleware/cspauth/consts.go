package cspauth

const (
	encryptionBlockSizeBytes = 16
	cspEncryptionIterations  = 1000
	lenRandomChars           = 6
	minAppKeyLen             = 12
	appKeyHeader             = "ak"
	authContextHeader        = "ac"
)
