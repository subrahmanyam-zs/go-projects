package config

import "testing"

func TestMockConfig_Get(t *testing.T) {
	m := MockConfig{
		Data: map[string]string{
			"key": "value",
		},
	}

	if m.Get("key") != "value" || m.Get("random") != "" {
		t.Error("Get not working")
	}

	if m.GetOrDefault("key", "test") != "value" ||
		m.GetOrDefault("random", "test") != "test" {
		t.Error("GetOrDefault not working")
	}
}
