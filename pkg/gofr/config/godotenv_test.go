package config

import (
	"bytes"
	"os"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

// This Test is to check if environment variables are loaded from configs/.env on
// initialization of new gofr object
func TestReadConfig(t *testing.T) {
	os.Setenv("GOFR_ENV", "test")

	testCases := []struct {
		envKey   string
		envValue string
	}{
		{"TEST", "test"},         // Load from original .env file
		{"EXAMPLE", "example"},   // Load from .test.env file
		{"OVERWRITE", "success"}, // The key is present in both but value should be from .test.env
	}

	conf := NewGoDotEnvProvider(t, "../../../configs")

	for _, tc := range testCases {
		val := conf.Get(tc.envKey)
		if val != tc.envValue {
			t.Errorf("Test Failed.\t Expected env value: %v\tGot env value: %v", tc.envValue, val)
		}
	}
}

func TestNewGoDotEnvProvider(t *testing.T) {
	logger := log.NewMockLogger(new(bytes.Buffer))

	os.Unsetenv("APP_NAME")
	os.Unsetenv("GOFR_ENV")

	// testing case where folder doesn't exist
	c := NewGoDotEnvProvider(logger, "./configs")

	if app := c.Get("APP_NAME"); app != "" {
		t.Errorf("FAILED, Expected: %s, Got: %s", "", app)
	}

	// testing case where folder does exist
	NewGoDotEnvProvider(logger, "../../../configs")

	expected := "gofr"

	if got := os.Getenv("APP_NAME"); got != expected {
		t.Errorf("FAILED, Expected: %s, Got: %s", expected, got)
	}
}

func TestGoDotEnvProvider_GetOrDefault(t *testing.T) {
	var (
		key   = "random123"
		value = "value123"
		g     = new(GoDotEnvProvider)
	)

	os.Setenv(key, value)

	if got := g.GetOrDefault(key, "default"); got != value {
		t.Errorf("FAILED, Expected: %v, Got: %v", value, got)
	}

	got := g.GetOrDefault("someKeyThatDoesntExist", "default")
	if got != "default" {
		t.Errorf("FAILED, Expected: default, Got: %v", got)
	}
}
