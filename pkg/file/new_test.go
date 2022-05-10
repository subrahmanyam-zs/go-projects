package file

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
)

func TestNewFile(t *testing.T) {
	testcases := []struct {
		store    string
		fileName string
		fileMode Mode
		expErr   error
	}{
		{Local, "test.txt", READ, nil},
		{Local, "test.txt", WRITE, nil},
		{Local, "test.txt", READWRITE, nil},
		{Local, "test.txt", APPEND, nil},
		{Azure, "test.txt", READ, nil},
		{Azure, "test.txt", WRITE, nil},
		{Azure, "test.txt", READWRITE, nil},
		{Azure, "test.txt", APPEND, nil},
		{AWS, "test.txt", READWRITE, nil},
		{GCP, "test.txt", WRITE, fmt.Errorf("dialing: google: could not find default " +
			"credentials. See https://developers.google.com/accounts/docs/application-default-credentials for more information.")},
		{SFTP, "test.txt", READ, &net.OpError{}},
		{"invalid file stroage", "test.txt", READ, errors.InvalidFileStorage},
	}

	for _, tc := range testcases {
		c := config.MockConfig{Data: map[string]string{
			"FILE_STORE": tc.store,
		}}
		_, err := NewWithConfig(&c, tc.fileName, tc.fileMode)

		assert.IsType(t, tc.expErr, err)
	}
}
