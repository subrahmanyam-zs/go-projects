package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
)

type pusher struct {
	http.ResponseWriter
	err    error // err to return from Push()
	target string
	opts   *http.PushOptions
}

func (p pusher) Push(target string, opts *http.PushOptions) error {
	// record passed arguments for later inspection
	p.target = target
	p.opts = opts

	return p.err
}

func TestHomeHandler(t *testing.T) {
	app := gofr.New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pusher, ok := w.(http.Pusher)
		if !ok {
			t.Fatal(ok)
		}
		err := pusher.Push("/", nil)
		if err != nil {
			t.Error(err)
		}
	})

	server := httptest.NewTLSServer(handler)
	defer server.Close()

	req, _ := request.NewMock(http.MethodGet, server.URL+"/", nil)
	w := pusher{}
	ctx := gofr.NewContext(responder.NewContextualResponder(w, req), request.NewHTTPRequest(req), app)

	tests := []struct {
		desc string
		push pusher
		err  error
	}{
		{"push without error", pusher{err: nil}, nil},
		{"push with error", pusher{err: &errors.Response{Reason: "test error"}}, &errors.Response{Reason: "test error"}},
	}

	for i, tc := range tests {
		ctx.ServerPush = tc.push

		_, err := HomeHandler(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestServeStatic(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://dummy", nil)

	ctx := gofr.NewContext(responder.NewContextualResponder(w, req), request.NewHTTPRequest(req), nil)

	ctx.SetPathParams(map[string]string{
		"name": "app.js",
	})

	resp, err := ServeStatic(ctx)
	assert.Nil(t, err)

	assert.Equal(t, resp, template.Template{File: "app.js"})
}
