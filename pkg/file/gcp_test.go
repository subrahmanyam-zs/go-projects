package file

import (
	"testing"

	"github.com/zopsmart/gofr/pkg/gofr/config"
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
