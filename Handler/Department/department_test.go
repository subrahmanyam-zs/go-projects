package Department

import (
	"EmployeeDepartment/Handler/Entities"
	"bytes"
	"errors"
	"github.com/google/uuid"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestEmployeePost(t *testing.T) {

	testCases := []struct {
		desc           string
		input          []byte
		expectedOutput []byte
	}{
		{"Valid input", []byte(`{"id":1, "name":"HR","floorNo": 1}`), []byte(`{"Id":1,"Name":"HR","FloorNo":1}`)},
		{"Invalid input", []byte(`{"id":0,"name":"Tech","floorNo":2}`), []byte("Invalid id")},
		{"for Unmarshal error", []byte(`{"id":"2","name":"hr","floorNo":2}`), []byte("Unmarshal Error")},
	}

	for i, tc := range testCases {
		req := httptest.NewRequest("POST", "/employee", bytes.NewReader(tc.input))
		w := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.PostHandler(w, req)

		if !reflect.DeepEqual(w.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Body.String(), string(tc.expectedOutput))
		}
	}
}

type mockDatastore struct {
	id uuid.UUID
}

func (m mockDatastore) Create(department Entities.Department) (Entities.Department, error) {
	
	if department.Id == 0 {
		return Entities.Department{}, errors.New("error")
	}

	return Entities.Department{1, "HR", 1}, nil
}
