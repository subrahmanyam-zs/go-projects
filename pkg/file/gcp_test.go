package file

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
)

func TestNewGCP(t *testing.T) {
	c := &config.MockConfig{Data: map[string]string{
		"FILE_STORE":              "GCP",
		"GCP_STORAGE_CREDENTIALS": "gcpKey",
		"GCP_STORAGE_BUCKET_NAME": "gcpBucket",
	}}

	_, err := NewWithConfig(c, "test.txt", READ)
	if err == nil {
		t.Errorf("For wrong config GCP client should not be created")
	}
}
func Test_list_gcp(t *testing.T) {
	s := &gcp{}
	expErr := ErrListingNotSupported
	_, err := s.list("test")
	assert.Equalf(t, expErr, err, "Test case failed.\nExpected: %v, got: %v", expErr, err)
}
