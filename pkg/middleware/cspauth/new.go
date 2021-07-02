package cspauth

import (
	"sync"
)

// CSP validates Auth Context Header
type CSP struct {
	cache
	sharedKey string
}

func New(sharedKey string) *CSP {
	return &CSP{
		cache: cache{
			keys: make(map[string]EncryptionKey),
			mu:   sync.RWMutex{},
		},
		sharedKey: sharedKey,
	}
}
