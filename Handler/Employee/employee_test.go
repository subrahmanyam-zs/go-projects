package Employee

import (
	"EmployeeDepartment/Handler/Entities"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var uid uuid.UUID = generateUUID()

func TestEmployeeHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          Entities.Employee
		expectedOutput Entities.Employee
	}{
		{"valid case", Entities.Employee{uid, "jason", "12-06-2002", "Bangalore", "CSE", 1}, Entities.Employee{uid, "jason", "12-06-2002", "Bangalore", "CSE", 1}},
		{"Invalid case", Entities.Employee{}, Entities.Employee{}},
	}

	for i, tc := range testcases {
		json_datareq, _ := json.Marshal(tc.input)
		json_datares, _ := json.Marshal(tc.expectedOutput)

		reader := bytes.NewReader(json_datareq)
		req := httptest.NewRequest(http.MethodPost, "/post", reader)
		res := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.PostHandler(res, req)

		if !reflect.DeepEqual(res.Body, bytes.NewBuffer(json_datares)) {
			t.Errorf("testcase %d failed got %v expected %v", i+1, res.Body.String(), string(json_datares))
		}
	}

}

type mockDatastore struct{}

func (m mockDatastore) Create(employee Entities.Employee) (Entities.Employee, error) {
	name := employee.Name
	if name == "" {
		return Entities.Employee{}, nil
	}
	return Entities.Employee{uid, "jason", "12-06-2002", "Bangalore", "CSE", 1}, nil
}
