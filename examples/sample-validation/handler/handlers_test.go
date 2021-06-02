package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestValidateEntry(t *testing.T) {
	k := gofr.New()

	testcases := []struct {
		method        string
		target        string
		body          Details
		expectedError error
	}{
		{"POST", "http://localhost:9010/phone", Details{"+912123456789098", "c.r@yahoo.com"}, nil},
		{"POST", "http://localhost:9010/phone", Details{}, errors.InvalidParam{Param: []string{"Phone Number length"}}},
	}

	for index, tc := range testcases {
		tempBody, _ := json.Marshal(tc.body)
		body := bytes.NewReader(tempBody)
		r := httptest.NewRequest(tc.method, tc.target, body)
		req := request.NewHTTPRequest(r)

		c := gofr.NewContext(nil, req, k)

		_, err := ValidateEntry(c)
		if !(reflect.DeepEqual(err, tc.expectedError)) {
			t.Errorf("Test FAILED for %v, got error: %v, expected error : %v", index+1, err, tc.expectedError)
		}
	}
}
