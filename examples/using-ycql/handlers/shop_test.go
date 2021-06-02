package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"developer.zopsmart.com/go/gofr/examples/using-ycql/entity"
	"developer.zopsmart.com/go/gofr/examples/using-ycql/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func initializeHandlerTest(t *testing.T) (*store.MockShop, Shop, *gofr.Gofr) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	shopStore := store.NewMockShop(ctrl)
	methods := New(shopStore)
	k := gofr.New()

	return shopStore, methods, k
}

func TestShop_Get(t *testing.T) {
	tests := []struct {
		queryParams  string
		expectedResp []entity.Shop
		mockErr      error
	}{
		{"id=1", []entity.Shop{{ID: 1, Name: "PhoenixMall", Location: "Gaya", State: "Bihar"}}, nil},
		{"name=PhoenixMall&location=Gaya", []entity.Shop{{ID: 1, Name: "PhoenixMall", Location: "Gaya", State: "Bihar"}}, nil},
		{"", []entity.Shop{
			{ID: 1, Name: "PhoenixMall", Location: "Gaya", State: "Bihar"},
			{ID: 2, Name: "GarudaMall", Location: "Dhanbad", State: "bihar"},
		}, nil},
	}

	shopStore, shop, k := initializeHandlerTest(t)

	for i, tc := range tests {
		req := httptest.NewRequest("GET", "/shop?"+tc.queryParams, nil)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)

		shopStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.expectedResp)

		resp, err := shop.Get(context)

		if err != tc.mockErr {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.mockErr, err)
		}

		if !reflect.DeepEqual(tc.expectedResp, resp) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.expectedResp, resp)
		}
	}
}

func TestShop_Create(t *testing.T) {
	tests := []struct {
		callGet       bool
		input         string
		expectedResp  interface{}
		mockGetOutput []entity.Shop
		mockErr       error
	}{
		{true, `{"id":4, "name": "UBCity", "location":"HSR", "State":"karnataka"}`,
			[]entity.Shop{{ID: 4, Name: "UBCity", Location: "HSR", State: "karnataka"}}, nil, nil},
		{
			true, `{"id": 3, "name": "UBCity", "location":"Bangalore", "state":"karnataka"}`,
			[]entity.Shop{}, nil, errors.EntityAlreadyExists{},
		},
		{false, `{"id":"3", "name":"UBCity", "location":"Bangalore", "state":"karnataka"}`, nil, nil,
			&json.UnmarshalTypeError{Value: "string", Type: reflect.TypeOf(40), Offset: 9, Struct: "Shop", Field: "id"}},
	}

	shopStore, shop, k := initializeHandlerTest(t)

	for i, tc := range tests {
		in := strings.NewReader(tc.input)
		req := httptest.NewRequest("POST", "/dummy", in)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)

		if tc.callGet {
			shopStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.mockGetOutput).AnyTimes()
			shopStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tc.expectedResp, tc.mockErr)
		}

		resp, err := shop.Create(context)
		if !reflect.DeepEqual(err, tc.mockErr) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.mockErr, err)
		}

		if !reflect.DeepEqual(tc.expectedResp, resp) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.expectedResp, resp)
		}
	}
}

func TestShop_Update(t *testing.T) {
	tests := []struct {
		callGet       bool
		id            string
		input         string
		expectedResp  interface{}
		mockErr       error
		mockGetOutput []entity.Shop
	}{
		{
			true, "3", `{ "name":  "SelectCityWalk", "location":  "tirupati", "State": "AndhraPradesh"}`,
			[]entity.Shop{{ID: 3, Name: "SelectCityWalk", Location: "tirupati", State: "AndhraPradesh"}},
			nil, []entity.Shop{{ID: 3, Name: "SelectCityWalk", Location: "tirupati", State: "AndhraPradesh"}},
		},
		{
			true, "3", `{  "location":"Dhanbad"  , "state": "Jharkhand"}`,
			[]entity.Shop{{ID: 3, Name: "SelectCityWalk", Location: "Dhanbad", State: "Jharkhand"}},
			nil, []entity.Shop{{ID: 3, Name: "SelectCityWalk", Location: "Dhanbad", State: "Jharkhand"}},
		},
		{false, "3", `{ "name":  "SkyWalkMall", "location":30, "state": "karnataka"}`, nil,
			&json.UnmarshalTypeError{Value: "number", Type: reflect.TypeOf("30"), Offset: 39, Struct: "Shop", Field: "location"},
			[]entity.Shop{}},
		{
			true, "5", `{ "name":  "Mali", "age":   40, "State": "AP"}`, nil,
			errors.EntityNotFound{Entity: "person", ID: "5"}, []entity.Shop{}},
	}

	shopStore, shop, k := initializeHandlerTest(t)

	for i, tc := range tests {
		in := strings.NewReader(tc.input)
		req := httptest.NewRequest("PUT", "/shop/"+tc.id, in)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)

		context.SetPathParams(map[string]string{
			"id": tc.id,
		})

		if tc.callGet {
			shopStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.mockGetOutput)
			shopStore.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tc.expectedResp, tc.mockErr)
		}

		resp, err := shop.Update(context)

		if !reflect.DeepEqual(tc.mockErr, err) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.mockErr, err)
		}

		if !reflect.DeepEqual(tc.expectedResp, resp) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.expectedResp, resp)
		}
	}
}

func TestPerson_Delete(t *testing.T) {
	tests := []struct {
		callDel       bool
		id            string
		expectedResp  interface{}
		expectedErr   error
		mockGetOutput []entity.Shop
	}{
		{false, "5", nil, errors.EntityNotFound{Entity: "person", ID: "5"}, nil},
		{true, "3", nil, nil, []entity.Shop{{ID: 3, Name: "Kali", Location: "HSR", State: "karnataka"}}},
	}

	personStore, person, k := initializeHandlerTest(t)

	for i, tc := range tests {
		req := httptest.NewRequest("PUT", "/persons/"+tc.id, nil)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)

		context.SetPathParams(map[string]string{
			"id": tc.id,
		})

		personStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.mockGetOutput)

		if tc.callDel == true {
			personStore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tc.expectedErr)
		}

		resp, err := person.Delete(context)

		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.expectedErr, err)
		}

		if !reflect.DeepEqual(tc.expectedResp, resp) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.expectedResp, resp)
		}
	}
}
