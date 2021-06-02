package config

type MockConfig struct {
	Data map[string]string
}

func (m *MockConfig) Get(key string) string {
	return m.Data[key]
}

func (m *MockConfig) GetOrDefault(key, defaultValue string) string {
	v, ok := m.Data[key]
	if ok {
		return v
	}

	return defaultValue
}
