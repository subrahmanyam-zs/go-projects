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

func (c *cache) get(appKey, sharedKey string) EncryptionKey {
	c.mu.Lock()

	keys, ok := c.keys[appKey]

	c.mu.Unlock()

	if !ok {
		return c.set(appKey, sharedKey)
	}

	return keys
}

func (c *cache) set(appKey, sharedKey string) EncryptionKey {
	encryptionKey := CreateKey([]byte(appKey), []byte(appKey[:12]), EncryptionKeyLen)
	iv := CreateKey([]byte(sharedKey), []byte(appKey[:12]), IVLength)

	c.mu.Lock()

	c.keys[appKey] = EncryptionKey{
		encryptionKey: encryptionKey,
		iv:            iv,
	}

	c.mu.Unlock()

	return EncryptionKey{encryptionKey, iv}
}
