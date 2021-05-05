package handler

import (
	"context"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
	"github.com/zopsmart/gofr/pkg/gofr/responder"
)

func TestCustomer_GetByID(t *testing.T) {
	testcases := []struct {
		id       string
		response interface{}
		err      error
	}{
		{"1", "1", nil},
		{"", nil, errors.MissingParam{Param: []string{"id"}}},
	}

	for i := range testcases {
		req := httptest.NewRequest("GET", "http://dummy", nil)
		c := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		c.SetPathParams(map[string]string{"id": testcases[i].id})

		h := New(mockServicer{})

		resp, err := h.Get(c)
		if !reflect.DeepEqual(err, testcases[i].err) {
			t.Errorf("[TEST%d]Failed. Got%v\tExpected %v\n", i+1, err, testcases[i].err)
		}

		if !reflect.DeepEqual(resp, testcases[i].response) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, testcases[i].response)
		}

	}
}

type mockServicer struct {
}

func (m mockServicer) GetBrandByID(ctx context.Context, id string) interface{} {
	if id == "1" {
		return "1"
	}

	return nil
}

func (m mockServicer) PropagateHeaders(headers ...string) {}
