package gofr

import (
	ctx "context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	gofrErrors "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
)

// routeKeySetter is used to set the routKey in the request context
func routeKeySetter(w http.ResponseWriter, r *http.Request) *http.Request {
	// dummy handler for setting routeKey
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r = req
	})

	muxRouter := mux.NewRouter()
	muxRouter.NewRoute().Path(r.URL.Path).Methods(r.Method).Handler(handler)
	muxRouter.ServeHTTP(w, r)

	return r
}

// TestHandler_ServeHTTP_StatusCode tests the different combination of statusCode and errors.
func TestHandler_ServeHTTP_StatusCode(t *testing.T) {
	testCases := []struct {
		error      error
		statusCode int
		code       string
		data       interface{}
	}{
		{gofrErrors.InvalidParam{Param: []string{"organizationId"}}, http.StatusBadRequest, "Invalid Parameter", nil},
		{gofrErrors.EntityAlreadyExists{}, http.StatusOK, "", "some data"},
		{gofrErrors.EntityNotFound{Entity: "user", ID: "2"}, http.StatusNotFound, "Entity Not Found", nil},
		{errors.New("unexpected response from internal dependency"), http.StatusInternalServerError, "Internal Server Error", nil},
		{nil, http.StatusOK, "", nil},
		{gofrErrors.MissingParam{Param: []string{"organizationId"}}, http.StatusBadRequest, "Missing Parameter", nil},
		{gofrErrors.MissingParam{Param: []string{"organizationId"}}, http.StatusPartialContent, "Missing Parameter",
			map[string]interface{}{"name": "Alice"}},
		{gofrErrors.MissingParam{Param: []string{"organizationId"}}, http.StatusBadRequest, "Missing Parameter", types.Response{}},
		{nil, http.StatusOK, "", &types.Response{Data: map[string]interface{}{"name": "Alice"}}},
	}

	for i, tc := range testCases {
		tc := tc
		k := New()
		w := newCustomWriter()
		r := httptest.NewRequest("GET", "/Dummy", nil)
		r = routeKeySetter(w, r)
		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)
		*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

		Handler(func(c *Context) (interface{}, error) {
			return tc.data, tc.error
		}).ServeHTTP(w, r)

		if w.Status != tc.statusCode {
			t.Errorf("TestCase[%v]\nIncorrect status code: \nGot\n%v\nExpected\n%v\n", i, w.Status, tc.statusCode)
		}

		if tc.code != "" && !strings.Contains(w.Body, tc.code) {
			t.Errorf("FAILED, Expected %v in response body", tc.code)
		}
	}
}

// TestHandler_ServeHTTP_ErrorFormat
func TestHandler_ServeHTTP_ErrorFormat(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest("GET", "/Dummy", nil)
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

	Handler(func(c *Context) (interface{}, error) {
		return nil, &gofrErrors.Response{StatusCode: 400, Code: "Invalid name", Reason: "The name in the parameter is incorrect"}
	}).ServeHTTP(w, r)

	e := struct {
		Errors []gofrErrors.Response `json:"errors"`
	}{[]gofrErrors.Response{}}

	_ = json.Unmarshal([]byte(w.Body), &e)

	if len(e.Errors) != 1 && e.Errors[0].Code != "Invalid name" &&
		e.Errors[0].Reason != "The name in the parameter is incorrect" && e.Errors[0].Detail != 1 {
		t.Errorf("Error formating failed.")
	}
}

// TestHandler_ServeHTTP_Content_Type tests the JSON content type for a response
func TestHandler_ServeHTTP_Content_Type(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest("GET", "/Dummy", nil)
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

	Handler(func(c *Context) (interface{}, error) {
		return "hi", nil
	}).ServeHTTP(w, r)

	contentType := "application/json"
	if w.Headers.Get("Content-Type") != contentType {
		t.Errorf("Response and Request Content Type does not match")
	}
}

// TestHandler_ServeHTTP
func TestHandler_ServeHTTP(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest("GET", "/Dummy", nil)
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))
	Handler(func(c *Context) (interface{}, error) {
		p := product{Name: "Orange", CategoryID: 1}
		data := struct {
			Product product `json:"product"`
		}{p}

		return data, nil
	}).ServeHTTP(w, r)

	expectedResponse := []byte(`{"data":{"product":{"name":"Orange","categoryId":1}}}`)

	if !reflect.DeepEqual(string(expectedResponse), strings.TrimSpace(w.Body)) {
		t.Errorf("Failed. Incorrect Response format. Expected %v\tGot %v\n", string(expectedResponse), w.Body)
	}
}

type product struct {
	Name       string `json:"name"`
	CategoryID int    `json:"categoryId"`
}

// TestHandler_ServeHTTP_XML tests the XML content type for a response
func TestHandler_ServeHTTP_XML(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest("GET", "/Dummy", nil)
	r.Header.Add("Content-type", "application/xml")
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))
	expectedError := gofrErrors.Response{StatusCode: http.StatusInternalServerError, Reason: "something unexpected occurred"}

	Handler(func(c *Context) (interface{}, error) {
		return nil, errors.New("something unexpected occurred")
	}).ServeHTTP(w, r)

	if !strings.Contains(w.Body, expectedError.Reason) || !strings.Contains(w.Body, strconv.Itoa(expectedError.StatusCode)) {
		t.Errorf("Error formating failed for xml.")
	}
}

// TestHandler_ServeHTTP_Text tests the TEXT content type for a response
func TestHandler_ServeHTTP_Text(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest("GET", "/Dummy", nil)
	r.Header.Add("Content-type", "text/plain")
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

	Handler(func(c *Context) (interface{}, error) {
		return nil, errors.New("something unexpected occurred")
	}).ServeHTTP(w, r)

	if w.Body != "something unexpected occurred" {
		t.Errorf("Error formating failed")
	}
}

// TestHandler_ServeHTTP_PartialContent
func TestHandler_ServeHTTP_PartialContent(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest("GET", "/Dummy", nil)
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))
	Handler(func(c *Context) (interface{}, error) {
		p := product{Name: "Orange", CategoryID: 1}
		data := struct {
			Product product `json:"product"`
		}{p}

		return data, &gofrErrors.Response{Reason: "test error", DateTime: gofrErrors.DateTime{Value: "2020-07-01T14:54:41Z", TimeZone: "IST"}}
	}).ServeHTTP(w, r)

	// nolint:lll // response should be of this type
	expectedResponse := []byte(`{"data":{"errors":[{"code":"","reason":"test error","datetime":{"value":"2020-07-01T14:54:41Z","timezone":"IST"}}],"product":{"categoryId":1,"name":"Orange"}}}`)

	if !reflect.DeepEqual(string(expectedResponse), strings.TrimSpace(w.Body)) {
		t.Errorf("Failed. Incorrect Response format. Expected %v\tGot %v\n", string(expectedResponse), w.Body)
		return
	}

	if w.Status != 206 {
		t.Errorf("Failed. StatusCode expected: 206, got: %v", w.Status)
	}
}

// TestHandler_ServeHTTP_EntityAlreadyExists
func TestHandler_ServeHTTP_EntityAlreadyExists(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest(http.MethodPost, "/Dummy", nil)
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

	Handler(func(c *Context) (interface{}, error) {
		p := product{Name: "Orange", CategoryID: 1}
		data := struct {
			Product product `json:"product"`
		}{p}

		return data, gofrErrors.EntityAlreadyExists{}
	}).ServeHTTP(w, r)

	expectedResponse := []byte(`{"data":{"product":{"name":"Orange","categoryId":1}}}`)

	if !reflect.DeepEqual(string(expectedResponse), strings.TrimSpace(w.Body)) {
		t.Errorf("Failed. Incorrect Response format. Expected %v\tGot %v\n", string(expectedResponse), w.Body)
		return
	}

	if w.Status != 200 {
		t.Errorf("Failed. StatusCode expected: 200, got: %v", w.Status)
	}
}

// Test_HealthInvalidMethod checks the health for method.
func Test_HealthInvalidMethod(t *testing.T) {
	testCases := []struct {
		error      error
		statusCode int
		code       string
		data       interface{}
		method     string
	}{
		{gofrErrors.MethodMissing{}, http.StatusMethodNotAllowed, "", nil, http.MethodPost},
		{nil, http.StatusOK, "", nil, "GET"},
	}

	for _, tc := range testCases {
		tc := tc
		k := New()
		w := newCustomWriter()
		r := httptest.NewRequest(tc.method, "/.well-known/health-check", nil)
		r = routeKeySetter(w, r)
		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)
		*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

		Handler(func(c *Context) (interface{}, error) {
			return tc.data, tc.error
		}).ServeHTTP(w, r)

		if w.Status != tc.statusCode {
			t.Errorf("\nIncorrect status code: \nGot\n%v\nExpected\n%v\n", w.Status, tc.statusCode)
		}
	}
}

func testNil() (*types.Response, error) {
	return nil, gofrErrors.MissingParam{}
}

//nolint:unparam //there is only one case to test
func testEmptyStruct() (*product, error) {
	return nil, gofrErrors.InvalidParam{Param: []string{"filter"}}
}

// TestHTTP_Respond_Nil tests nil and empty struct response
func TestHTTP_Respond_Nil(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest("GET", "/Dummy", nil)
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))
	{
		// Test for nil types.Response
		expectedError := gofrErrors.Response{StatusCode: http.StatusBadRequest, Code: "Missing Parameter",
			Reason: "This request is missing parameters"}

		Handler(func(c *Context) (interface{}, error) {
			return testNil()
		}).ServeHTTP(w, r)

		if !strings.Contains(w.Body, expectedError.Reason) || !strings.Contains(w.Body, expectedError.Code) {
			t.Errorf("Error formating failed.")
		}

		if w.Status != expectedError.StatusCode {
			t.Errorf("Failed. StatusCode expected: 400, got: %v", w.Status)
		}
	}

	{
		// Test for empty struct, where type is non nil but value is nil
		expectedError := gofrErrors.Response{StatusCode: http.StatusBadRequest, Code: "Invalid Parameter",
			Reason: "Incorrect value for parameter: filter"}

		Handler(func(c *Context) (interface{}, error) {
			return testEmptyStruct()
		}).ServeHTTP(w, r)

		if !strings.Contains(w.Body, expectedError.Reason) || !strings.Contains(w.Body, expectedError.Code) {
			t.Errorf("Error formating failed.")
		}

		if w.Status != expectedError.StatusCode {
			t.Errorf("Failed. StatusCode expected: 400, got: %v", w.Status)
		}
	}
}

// TestHTTP_Respond_Delete tests the DELETE method
func TestHTTP_Respond_Delete(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

	Handler(func(c *Context) (interface{}, error) {
		return nil, nil
	}).ServeHTTP(w, r)

	expectedResponse := []byte(`{"data":null}`)

	if !reflect.DeepEqual(string(expectedResponse), strings.TrimSpace(w.Body)) {
		t.Errorf("Failed. Incorrect Response format. Expected %v\tGot %v\n", string(expectedResponse), w.Body)
		return
	}

	if w.Status != http.StatusNoContent {
		t.Errorf("Failed. StatusCode expected: %v, got: %v", http.StatusNoContent, w.Status)
	}
}

// TestHandler_ServeHTTP_Error tests the different error cases returned from Respond
func TestHandler_ServeHTTP_Error(t *testing.T) {
	k := New()
	w := newCustomWriter()
	r := httptest.NewRequest("GET", "/error", nil)
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

	{
		// Error is present but only status code is set and no body
		expectedError := gofrErrors.Response{StatusCode: http.StatusInternalServerError}
		Handler(func(c *Context) (interface{}, error) {
			return nil, &gofrErrors.Response{StatusCode: http.StatusInternalServerError}
		}).ServeHTTP(w, r)

		if w.Status != expectedError.StatusCode {
			t.Errorf("Failed StatusCode expected: %v, got: %v", expectedError.StatusCode, w.Status)
		}
	}

	{
		// Error is set and returned in the body
		expectedError := gofrErrors.Response{StatusCode: http.StatusInternalServerError, Code: "TEST_ERROR",
			Reason: "test error occurred"}
		Handler(func(c *Context) (interface{}, error) {
			return nil, &gofrErrors.Response{StatusCode: http.StatusInternalServerError,
				Code:   "TEST_ERROR",
				Reason: "test error occurred",
			}
		}).ServeHTTP(w, r)

		if w.Status != expectedError.StatusCode {
			t.Errorf("Failed StatusCode expected: %v, got: %v", expectedError.StatusCode, w.Status)
		}
		if !strings.Contains(w.Body, expectedError.Reason) || !strings.Contains(w.Body, expectedError.Code) {
			t.Errorf("Error formating failed.")
		}
	}
}

// Test_Head tests if HEAD request for an endpoint returns the correct content length in the response header
func Test_Head(t *testing.T) {
	k := New()
	w := newCustomWriter()
	// making the get request
	r := httptest.NewRequest(http.MethodGet, "/Dummy", nil)
	r = routeKeySetter(w, r)
	req := request.NewHTTPRequest(r)
	resp := responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

	h := Handler(func(c *Context) (interface{}, error) {
		return "hello", nil
	})

	h.ServeHTTP(w, r)
	// content length of GET response should be equal to HEAD response
	expected := w.Headers.Get("content-length")

	// making the HEAD request for the same endpoint
	r = httptest.NewRequest(http.MethodHead, "/Dummy", nil)
	r = routeKeySetter(w, r)
	req = request.NewHTTPRequest(r)
	resp = responder.NewContextualResponder(w, r)
	*r = *r.WithContext(ctx.WithValue(r.Context(), gofrContextkey, NewContext(resp, req, k)))

	got := w.Header().Get("content-length")
	if got != expected {
		t.Errorf("got %v\n expected %v\n", got, expected)
	}
}
