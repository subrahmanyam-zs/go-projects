package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
)

func TestHelloWorldHandler(t *testing.T) {
	c := gofr.NewContext(nil, nil, nil)

	resp, err := HelloWorld(c)
	if err != nil {
		t.Errorf("FAILED, got error: %v", err)
	}

	expected := "Hello World!"

	got := fmt.Sprintf("%v", resp)
	if got != expected {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, got)
	}
}

func TestHelloNameHandler(t *testing.T) {
	tests := []struct {
		name     string
		expected types.Response
	}{
		{
			name:     "SomeName",
			expected: types.Response{Data: "Hello SomeName", Meta: map[string]interface{}{"page": 1, "offset": 0}},
		},
		{
			name:     "Firstname Lastname",
			expected: types.Response{Data: "Hello Firstname Lastname", Meta: map[string]interface{}{"page": 1, "offset": 0}},
		},
	}

	for _, test := range tests {
		r := httptest.NewRequest("GET", "http://dummy/hello?name="+url.QueryEscape(test.name), nil)
		req := request.NewHTTPRequest(r)
		c := gofr.NewContext(nil, req, nil)

		resp, err := HelloName(c)
		if err != nil {
			t.Errorf("FAILED, got error: %v", err)
		}

		result, _ := resp.(types.Response)

		if result.Data != test.expected.Data {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expected, resp)
		}
	}
}

func TestPostNameHandler(t *testing.T) {
	var jsonStr = []byte(`{"Username":"username"}`)
	r := httptest.NewRequest("POST", "/post", bytes.NewBuffer(jsonStr))
	req := request.NewHTTPRequest(r)

	r.Header.Set("Content-Type", "application/json")

	k := gofr.New()
	c := gofr.NewContext(nil, req, k)

	resp, err := PostName(c)
	if err != nil {
		t.Errorf("FAILED, got error: %v", err)
	}

	var got Person

	respBytes, _ := json.Marshal(resp)
	_ = json.Unmarshal(respBytes, &got)

	expected := Person{Username: "username"}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, got)
	}
}

func TestPostNameHandlerfail(t *testing.T) {
	// invalid JSON passed
	var jsonStr = []byte(`{"Username":}`)

	r := httptest.NewRequest("POST", "/post", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	req := request.NewHTTPRequest(r)
	c := gofr.NewContext(nil, req, gofr.New())

	resp, err := PostName(c)
	if err == nil {
		t.Errorf("FAILED, got error: %v", err)
	}

	if resp != nil {
		t.Errorf("FAILED, Expected: nil, Got: %v", resp)
	}
}

func TestMultipleErrorHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "http://dummy/multiple-errors", nil)
	req := request.NewHTTPRequest(r)
	c := gofr.NewContext(nil, req, nil)

	_, err := MultipleErrorHandler(c)
	if err == nil {
		t.Errorf("FAILED, got nil expectedErr")
		return
	}

	expectedErr := `Incorrect value for parameter: EmailAddress
Parameter Address is required for this request`

	got := err.Error()
	if got != expectedErr {
		t.Errorf("FAILED, Expected: %v, Got: %v", expectedErr, got)
	}
}
