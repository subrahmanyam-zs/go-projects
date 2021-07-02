package cspauth

import "sync"

type cache struct {
	keys map[string]EncryptionKey
	mu   sync.RWMutex
}

type EncryptionKey struct {
	encryptionKey []byte // encryptionKey to be used for aes encryption/decryption
	iv            []byte // initial vector(iv) to be used for aes encryption/decryption
}

func (c *cache) Get(appKey string) EncryptionKey {
	c.mu.Lock()

	keys := c.keys[appKey]

	c.mu.Unlock()

	return keys
}

func (c *cache) Set(appKey, sharedKey string) {
	c.mu.Lock()

	_, ok := c.keys[appKey]

	c.mu.Unlock()

	if !ok {
		encryptionKey := createKey([]byte(appKey), []byte(appKey[:12]))
		iv := createKey([]byte(sharedKey), []byte(appKey[:12]))

		c.mu.Lock()

		c.keys[appKey] = EncryptionKey{
			encryptionKey: encryptionKey,
			iv:            iv,
		}

		c.mu.Unlock()
	}
}
