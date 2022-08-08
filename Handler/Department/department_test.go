package Department

import (
	"EmployeeDepartment/Handler/Entities"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestEmployeePost(t *testing.T) {

	testCases := []struct {
		reqBody  Entities.Department
		respBody Entities.Department
	}{
		{Entities.Department{1, "HR", 1}, Entities.Department{1, "HR", 1}},
		{Entities.Department{}, Entities.Department{}},
	}

	for i, v := range testCases {
		jsondatareq, _ := json.Marshal(v.reqBody)
		jsondataresp, _ := json.Marshal(v.respBody)
		req := httptest.NewRequest("POST", "/employee", bytes.NewReader(jsondatareq))
		w := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.PostHandler(w, req)

		if !reflect.DeepEqual(w.Body, bytes.NewBuffer(jsondataresp)) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Body.String(), string(jsondataresp))
		}
	}
}

type mockDatastore struct {
	id uuid.UUID
}

func (m mockDatastore) Create(department Entities.Department) (Entities.Department, error) {
	id := department.Id
	if id == 0 {
		return Entities.Department{}, nil
	}

	return Entities.Department{1, "HR", 1}, nil
}
