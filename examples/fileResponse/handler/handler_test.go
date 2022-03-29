package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestTemplateHandler(t *testing.T) {
	app := gofr.New()
	rootPath := t.TempDir()
	app.TemplateDir = rootPath + "/../templates"
	r := httptest.NewRequest(http.MethodGet, "http://dummy/test", nil)
	req := request.NewHTTPRequest(r)

	ctx := gofr.NewContext(nil, req, app)
	if _, err := FileHandler(ctx); err != nil {
		t.Errorf("FAILED, got error: %v", err)
	}
}
