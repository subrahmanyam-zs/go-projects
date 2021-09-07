package handler

import (
	"net/http"
	"net/http/httptest"
	"reflect"
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
	k := gofr.New()
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
	ctx := gofr.NewContext(responder.NewContextualResponder(w, req), request.NewHTTPRequest(req), k)

	testCases := []struct {
		push    pusher
		wantErr error
	}{
		{pusher{err: nil}, nil},
		{pusher{err: &errors.Response{Reason: "test error"}}, &errors.Response{Reason: "test error"}},
	}

	for _, tt := range testCases {
		ctx.ServerPush = tt.push

		_, err := HomeHandler(ctx)

		assert.Equal(t, tt.wantErr, err)
	}
}

func TestServeStatic(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://dummy", nil)

	c := gofr.NewContext(responder.NewContextualResponder(w, req), request.NewHTTPRequest(req), nil)

	c.SetPathParams(map[string]string{
		"name": "app.js",
	})

	got, err := ServeStatic(c)
	if err != nil {
		t.Errorf("ServeStatic() error = %v, wantErr nil", err)
		return
	}

	if !reflect.DeepEqual(got, template.Template{File: "app.js"}) {
		t.Errorf("ServeStatic() got = %v, want app.js", got)
	}
}
