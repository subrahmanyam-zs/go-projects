package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/model"
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func initializeTest(t *testing.T, method, url string, body []byte) (*store.MockPerson, handler, *gofr.Context) {
	mockStore := store.NewMockPerson(gomock.NewController(t))
	handler := New(mockStore)

	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	r := request.NewHTTPRequest(req)

	return mockStore, handler, gofr.NewContext(nil, r, gofr.New())
}

func TestGetByID(t *testing.T) {
	tcs := []struct {
		id   string
		resp interface{}
		out  interface{}
		err  error
	}{
		{"1", model.Person{ID: "1", Name: "gofr", Email: "gofr@gmail.com"}, model.Person{ID: "1", Name: "gofr", Email: "gofr@gmail.com"}, nil},
		{"2", model.Person{}, nil, errors.DB{}},
	}

	for i, tc := range tcs {
		mockStore, handler, ctx := initializeTest(t, http.MethodGet, "/person"+tc.id, nil)
		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		mockStore.EXPECT().Get(gomock.Any(), tc.id).Return(tc.resp, tc.err)

		out, err := handler.GetByID(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d],failed.", i)

		assert.Equal(t, tc.out, out, "TEST[%d],failed.", i)
	}
}

func TestDelete(t *testing.T) {
	mockStore, handler, ctx := initializeTest(t, http.MethodDelete, "/person", nil)
	ctx.SetPathParams(map[string]string{
		"id": "1",
	})

	mockStore.EXPECT().Delete(ctx, gomock.Any()).Return(nil)

	_, err := handler.Delete(ctx)

	assert.Equal(t, nil, err, "TEST,failed.")
}

func TestCreate(t *testing.T) {
	tcs := []struct {
		body []byte
		err  error
	}{
		{[]byte(`{"id":"1", "name":  "gofr", "email": "gofr@zopsmart.com"}`), nil},
		{[]byte(`{"id":"1", "name":  "gofr", "email": "gofr@zopsmart.com"}`), errors.DB{}},
	}

	for i, tc := range tcs {
		mockStore, handler, ctx := initializeTest(t, http.MethodPost, "/person", tc.body)

		mockStore.EXPECT().Create(ctx, gomock.Any()).Return(tc.err)

		_, err := handler.Create(ctx)

		assert.IsType(t, tc.err, err, "TESTCASE[%d],failed.", i)
	}
}

func TestUpdate(t *testing.T) {
	tcs := []struct {
		body []byte
		err  error
	}{
		{[]byte(`{"id":"1", "name":  "gofr", "email": "gofr@zopsmart.com"}`), nil},
		{[]byte(`{"id":"1", "name":  "gofr", "email": "gofr@zopsmart.com"}`), errors.DB{}},
	}

	for i, tc := range tcs {
		mockStore, handler, ctx := initializeTest(t, http.MethodPut, "/person", tc.body)
		ctx.SetPathParams(map[string]string{
			"id": "1",
		})

		mockStore.EXPECT().Update(ctx, gomock.Any()).Return(tc.err)

		_, err := handler.Update(ctx)

		assert.IsType(t, tc.err, err, "TESTCASE[%d],failed.", i)
	}
}

func Test_BindError(t *testing.T) {
	body := []byte(`{"id": 1, "name":  "gofr", "email": "gofr@zopsmart.com"}`)

	_, handler, ctx := initializeTest(t, http.MethodPut, "/person", body)
	ctx.SetPathParams(map[string]string{
		"id": "1",
	})

	var handlers []gofr.Handler

	handlers = append(handlers, handler.Update, handler.Create)

	for i := range handlers {
		_, err := handlers[i](ctx)
		assert.Error(t, err, "TEST,failed.")
	}
}
