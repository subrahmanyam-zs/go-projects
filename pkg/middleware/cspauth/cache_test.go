package cspauth

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cache_Get(t *testing.T) {
	tests := []struct {
		description string
		// input
		appKey string
		keys   map[string]EncryptionKey
		// output
		output EncryptionKey
	}{
		{
			appKey: "sample-app-key",
			keys:   map[string]EncryptionKey{"sample-app-key": {[]byte("sample-encryption-key"), []byte("sample-iv")}},
			output: EncryptionKey{[]byte("sample-encryption-key"), []byte("sample-iv")},
		},
	}

	for i, tc := range tests {
		c := &cache{tc.keys, sync.RWMutex{}}

		output := c.Get(tc.appKey)

		assert.Equal(t, tc.output, output, "TEST[%d], failed. %s", i+1, tc.description)
	}
}

func Test_cache_set(t *testing.T) {
	tests := []struct {
		description string
		// input
		appKey    string
		sharedKey string
		// output
		keys EncryptionKey
	}{
		{
			appKey:    "sample-app-key",
			sharedKey: "sample-shared-key",
			keys:      EncryptionKey{},
		},
	}

	for i, tc := range tests {
		c := &cache{make(map[string]EncryptionKey), sync.RWMutex{}}

		c.Set(tc.appKey, tc.sharedKey)
		output := c.Get(tc.appKey)

		assert.Equal(t, tc.keys, output, "TEST[%d], failed. %s", i+1, tc.description)
	}
}
