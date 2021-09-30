package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/sample-grpc/handler/grpc"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"

	"github.com/stretchr/testify/assert"
)

func TestExample_Get(t *testing.T) {
	tests := []struct {
		desc string
		id   string
		resp interface{}
		err  error
	}{
		{"get success case", "1", &grpc.Response{FirstName: "Henry", SecondName: "Marc"}, nil},
		{"get non existent entity case", "2", nil, errors.EntityNotFound{Entity: "name", ID: "2"}},
	}

	for i, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, "http://dummy?id="+tc.id, nil)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, nil)

		resp, err := Get(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}
