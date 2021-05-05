package handler

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
)

func TestTemplateHandler(t *testing.T) {
	k := gofr.New()
	rootPath, _ := os.Getwd()
	k.TemplateDir = rootPath + "/../templates"
	r := httptest.NewRequest("GET", "http://dummy/test", nil)
	req := request.NewHTTPRequest(r)

	c := gofr.NewContext(nil, req, k)
	if _, err := FileHandler(c); err != nil {
		t.Errorf("FAILED, got error: %v", err)
	}
}
