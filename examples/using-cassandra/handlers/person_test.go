package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"developer.zopsmart.com/go/gofr/examples/using-cassandra/entity"
	"developer.zopsmart.com/go/gofr/examples/using-cassandra/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func initializeHandlerTest(t *testing.T) (*store.MockPerson, Person, *gofr.Gofr) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	personStore := store.NewMockPerson(ctrl)
	person := New(personStore)
	k := gofr.New()

	return personStore, person, k
}

func TestPerson_Get(t *testing.T) {
	tests := []struct {
		queryParams  string
		expectedResp []*entity.Person
		mockErr      error
	}{
		{"id=1", []*entity.Person{{ID: 1, Name: "Aakash", Age: 25, State: "Bihar"}}, nil},
		{"name=Aakash&age=25", []*entity.Person{{ID: 1, Name: "Aakash", Age: 25, State: "Bihar"}}, nil},
		{"", []*entity.Person{
			{ID: 1, Name: "Aakash", Age: 25, State: "Bihar"},
			{ID: 2, Name: "himari", Age: 30, State: "bihar"},
		}, nil},
	}

	personStore, person, k := initializeHandlerTest(t)

	for i, tc := range tests {
		req := httptest.NewRequest("GET", "/persons?"+tc.queryParams, nil)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)

		personStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.expectedResp)

		resp, err := person.Get(context)

		if err != tc.mockErr {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.mockErr, err)
		}

		if !reflect.DeepEqual(tc.expectedResp, resp) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.expectedResp, resp)
		}
	}
}

func TestPerson_Create_InvalidInsertionIDAndJSONError(t *testing.T) {
	tests := []struct {
		callGet       bool
		input         string
		expectedResp  interface{}
		mockGetOutput []*entity.Person
		mockErr       error
	}{
		{
			false, `{"id":    3, "name":  "Kali", "age":   "40", "State": "karnataka"}`,
			nil, nil,
			&json.UnmarshalTypeError{Value: "string", Type: reflect.TypeOf(40), Offset: 43, Struct: "Person", Field: "age"},
		},
		{
			true, `{"id":    3, "name":  "Kali", "age":   40, "State": "karnataka"}`,
			nil, []*entity.Person{{ID: 3, Name: "Kali", Age: 40, State: "karnataka"}},
			errors.EntityAlreadyExists{},
		},
	}

	personStore, person, k := initializeHandlerTest(t)

	for i, tc := range tests {
		in := strings.NewReader(tc.input)
		req := httptest.NewRequest("POST", "/dummy", in)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)

		if tc.callGet == true {
			personStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.mockGetOutput).AnyTimes()
		}

		resp, err := person.Create(context)

		if !reflect.DeepEqual(tc.mockErr, err) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.mockErr, err)
		}

		if !reflect.DeepEqual(tc.expectedResp, resp) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.expectedResp, resp)
		}
	}
}

func TestPerson_Create(t *testing.T) {
	tests := []struct {
		input        string
		expectedResp interface{}
		mockErr      error
	}{
		{`{"id":4, "name":"Kali", "age":40, "State":"karnataka"}`,
			[]*entity.Person{{ID: 4, Name: "Kali", Age: 40, State: "karnataka"}}, nil},
	}

	personStore, person, k := initializeHandlerTest(t)

	for i, tc := range tests {
		in := strings.NewReader(tc.input)
		req := httptest.NewRequest("POST", "/dummy", in)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)

		personStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil)
		personStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tc.expectedResp, tc.mockErr)

		resp, err := person.Create(context)

		if err != tc.mockErr {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.mockErr, err)
		}

		if !reflect.DeepEqual(tc.expectedResp, resp) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.expectedResp, resp)
		}
	}
}

func TestPerson_InvalidUpdateIDAndJSONError(t *testing.T) {
	tests := []struct {
		callGet     bool
		id          string
		input       string
		expectedErr error
	}{
		{false, "3", `{ "name":  "Mali", "age":   "40", "State": "AP"}`,
			&json.UnmarshalTypeError{Value: "string", Type: reflect.TypeOf(40), Offset: 32, Struct: "Person", Field: "age"},
		},
		{true, "5", `{ "name":  "Mali", "age":   40, "State": "AP"}`, errors.EntityNotFound{Entity: "person", ID: "5"}},
	}

	personStore, person, k := initializeHandlerTest(t)

	for i, tc := range tests {
		in := strings.NewReader(tc.input)
		req := httptest.NewRequest("PUT", "/dummy/"+tc.id, in)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)

		context.SetPathParams(map[string]string{
			"id": tc.id,
		})

		if tc.callGet == true {
			personStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil)
		}

		_, err := person.Update(context)

		if !reflect.DeepEqual(tc.expectedErr, err) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.expectedErr, err)
		}
	}
}

func TestPerson_Update(t *testing.T) {
	tests := []struct {
		id            string
		input         string
		expectedResp  interface{}
		mockErr       error
		mockGetOutput []*entity.Person
	}{
		{
			"3", `{ "name":  "Mali", "age":   40, "State": "AP"}`, []*entity.Person{{ID: 3, Name: "Mali", Age: 40, State: "AP"}},
			nil, []*entity.Person{{ID: 3, Name: "Kali", Age: 40, State: "karnataka"}},
		},
		{
			"3", `{  "age":   35, "State": "AP"}`, []*entity.Person{{ID: 3, Name: "Kali", Age: 35, State: "AP"}},
			nil, []*entity.Person{{ID: 3, Name: "Kali", Age: 40, State: "karnataka"}},
		},
	}

	personStore, person, k := initializeHandlerTest(t)

	for i, tc := range tests {
		in := strings.NewReader(tc.input)
		req := httptest.NewRequest("PUT", "/persons/"+tc.id, in)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)

		context.SetPathParams(map[string]string{
			"id": tc.id,
		})

		personStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.mockGetOutput)
		personStore.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tc.expectedResp, tc.mockErr)

		resp, err := person.Update(context)

		if err != tc.mockErr {
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
		mockGetOutput []*entity.Person
	}{
		{false, "5", nil, errors.EntityNotFound{Entity: "person", ID: "5"}, nil},
		{true, "3", nil, nil, []*entity.Person{{ID: 3, Name: "Kali", Age: 40, State: "karnataka"}}},
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
