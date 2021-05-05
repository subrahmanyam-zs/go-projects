package http

import (
	"net/http/httptest"
	"testing"

	"github.com/zopsmart/gofr/examples/sample-grpc/handler/grpc"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
)

func TestExample_Get(t *testing.T) {
	tcs := []struct {
		id   string
		resp interface{}
		err  error
	}{
		{"1", &grpc.Response{FirstName: "First", SecondName: "Second"}, nil},
		{"2", nil, errors.EntityNotFound{Entity: "name", ID: "2"}},
	}

	for _, tc := range tcs {
		var (
			req = httptest.NewRequest("GET", "http://dummy?id="+tc.id, nil)
			r   = request.NewHTTPRequest(req)
			c   = gofr.NewContext(nil, r, nil)
		)

		resp, _ := Get(c)

		if resp == nil && tc.resp != nil {
			t.Errorf("FAILED, Expected: %v, Got: %v", tc.resp, resp)
			continue
		}
	}
}
