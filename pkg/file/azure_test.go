package file

import (
	"encoding/base64"
	"errors"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
)

func TestAzureFileOpen(t *testing.T) {
	testcases := []struct {
		azureAcName string
		azureAccKey string
		fileName    string
		blockSize   string
		parallelism string
		mode        Mode
		expErr      error
	}{
		{"alice", "some-random-key", "test.txt", "", "", READ, base64.CorruptInputError(4)},
		{"^@", "c29tZS1yYW5kb20tdGV4dA==", "test.txt", "", "", READWRITE, &url.Error{
			Op:  "parse",
			URL: "https://^@.blob.core.windows.net/container",
			Err: errors.New("net/url: invalid userinfo"),
		}},
		{"bob", "c29tZS1yYW5kb20tdGV4dA==", "test.txt", "abc", "", READWRITE, &strconv.NumError{
			Func: "Atoi",
			Num:  "abc",
			Err:  errors.New("invalid syntax"),
		}},
		{"bob", "c29tZS1yYW5kb20tdGV4dA==", "test.txt", "", "def", READWRITE, &strconv.NumError{
			Func: "Atoi",
			Num:  "def",
			Err:  errors.New("invalid syntax"),
		}},
		{"bob", "c29tZS1yYW5kb20tdGV4dA==", "test.txt", "4194304", "16", READWRITE, nil},
		{"bob", "c29tZS1yYW5kb20tdGV4dA==", "test.txt", "", "", READWRITE, nil},
	}

	for _, v := range testcases {
		c := &config.MockConfig{Data: map[string]string{
			"FILE_STORE":                "AZURE",
			"AZURE_STORAGE_ACCOUNT":     v.azureAcName,
			"AZURE_STORAGE_ACCESS_KEY":  v.azureAccKey,
			"AZURE_STORAGE_CONTAINER":   "container",
			"AZURE_STORAGE_BLOCK_SIZE":  v.blockSize,
			"AZURE_STORAGE_PARALLELISM": v.parallelism,
		}}

		_, err := NewWithConfig(c, "test.txt", READ)
		assert.Equal(t, v.expErr, err)
	}
}
