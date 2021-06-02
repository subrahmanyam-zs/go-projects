package http

import (
	"net/http/httptest"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/sample-grpc/handler/grpc"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
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
