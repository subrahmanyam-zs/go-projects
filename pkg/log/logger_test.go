package log

import "testing"

const defaultAppName = "gofr-app"

func TestNewLogger(t *testing.T) {
	l := newLogger()

	if l.app.Name != defaultAppName {
		t.Errorf("Expected APP_NAME: gofr-app     GOT:  %v", l.app.Name)
	}

	if l.app.Version != "dev" {
		t.Errorf("Expected APP_VERSION: dev    GOT:  %v", l.app.Version)
	}
}
