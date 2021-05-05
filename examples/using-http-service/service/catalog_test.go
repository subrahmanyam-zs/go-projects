package service

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/service"
)

func TestCustomer_GetByID(t *testing.T) {
	testcases := []struct {
		id       string
		testID   int
		response interface{}
	}{
		{"1", 0, map[string]interface{}{
			"data": map[string]interface{}{"customer": map[string]interface{}{"id": 1.0, "name": "Red"}}}},
		{"1", 1, nil},
		{"0", 1, nil},
	}

	for i := range testcases {
		h := catalog{httpService: mockServicer{testID: testcases[i].testID}}

		resp := h.GetBrandByID(context.TODO(), testcases[i].id)

		if !reflect.DeepEqual(resp, testcases[i].response) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, testcases[i].response)
		}

	}
}

type mockServicer struct {
	testID int
}

func (m mockServicer) Get(ctx context.Context, api string, params map[string]interface{}) (*service.Response, error) {
	if api == "1" {
		resp := service.Response{
			Body:       []byte(`{"data":{"customer":{"id":1,"name":"Red"}}}`),
			StatusCode: 200,
		}
		return &resp, nil
	}

	return nil, &errors.Response{Reason: "core error"}
}
func (m mockServicer) Bind(resp []byte, i interface{}) error {
	if m.testID == 0 {
		return json.Unmarshal(resp, &i)
	}

	return &errors.Response{Reason: "unmarshal error"}
}

func (m mockServicer) BindStrict(resp []byte, i interface{}) error {
	if m.testID == 0 {
		return json.Unmarshal(resp, &i)
	}

	return &errors.Response{Reason: "unmarshal error"}
}

func (m mockServicer) Post(ctx context.Context, api string, params map[string]interface{}, body []byte) (*service.Response, error) {
	return nil, nil
}

func (m mockServicer) Put(ctx context.Context, api string, params map[string]interface{}, body []byte) (*service.Response, error) {
	return nil, nil
}

func (m mockServicer) Patch(ctx context.Context, api string, params map[string]interface{}, body []byte) (*service.Response, error) {
	return nil, nil
}

func (m mockServicer) Delete(ctx context.Context, api string, body []byte) (*service.Response, error) {
	return nil, nil
}
func (m mockServicer) PropagateHeaders(headers ...string) {}
func (m mockServicer) SetSurgeProtectorOptions(isEnabled bool, customHeartbeatURL string, retryFrequencySeconds int) {
}

func (m mockServicer) GetWithHeaders(ctx context.Context, api string, params map[string]interface{}, headers map[string]string) (*service.Response, error) {
	return nil, nil
}

func (m mockServicer) PostWithHeaders(ctx context.Context, api string, params map[string]interface{}, body []byte, headers map[string]string) (*service.Response, error) {
	return nil, nil
}

func (m mockServicer) PutWithHeaders(ctx context.Context, api string, params map[string]interface{}, body []byte, headers map[string]string) (*service.Response, error) {
	return nil, nil
}

func (m mockServicer) PatchWithHeaders(ctx context.Context, api string, params map[string]interface{}, body []byte, headers map[string]string) (*service.Response, error) {
	return nil, nil
}

func (m mockServicer) DeleteWithHeaders(ctx context.Context, api string, body []byte, headers map[string]string) (*service.Response, error) {
	return nil, nil
}
