package cspauth

import (
	"sync"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

type Cache struct {
	keys map[string]EncryptionKey
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		keys: make(map[string]EncryptionKey),
		mu:   sync.RWMutex{},
	}
}

type EncryptionKey struct {
	encryptionKey []byte // encryptionKey to be used for aes encryption/decryption
	iv            []byte // initial vector(iv) to be used for aes encryption/decryption
}

type CSP struct {
	options *Options
	EncryptionKey
}

func New(logger log.Logger, opts *Options, cache *Cache) (*CSP, error) {
	if err := opts.validate(); err != nil {
		logger.Warnf("Invalid Options, %v", err)
		return nil, err
	}

	csp := &CSP{
		options: opts,
	}

	if val, ok := cache.keys[opts.AppKey]; ok {
		csp.encryptionKey = val.encryptionKey
		csp.iv = val.iv
	} else {
		csp.encryptionKey = createKey([]byte(opts.AppKey), []byte(opts.AppKey[:12]))
		csp.iv = createKey([]byte(opts.SharedKey), []byte(opts.AppKey[:12]))

		cache.mu.Lock()

		cache.keys[opts.AppKey] = EncryptionKey{
			encryptionKey: csp.encryptionKey,
			iv:            csp.iv,
		}

		cache.mu.Unlock()
	}

	return csp, nil
}

// Options used to initialize CSP
type Options struct {
	MachineName string
	IPAddress   string
	AppKey      string
	SharedKey   string
	AppID       string
}

func (o *Options) validate() error {
	if o.SharedKey == "" {
		return ErrEmptySharedKey
	}

	if len(o.AppKey) < minAppKeyLen {
		return ErrEmptyAppKey
	}

	if o.AppID == "" {
		return ErrEmptyAppID
	}

	return nil
}
