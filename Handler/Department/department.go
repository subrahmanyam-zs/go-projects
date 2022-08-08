package Department

import (
	"EmployeeDepartment/Handler/Entities"
	"EmployeeDepartment/Store"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DepartmentHandler struct {
	datastore Store.Department
}

func New(department Store.Department) DepartmentHandler {
	return DepartmentHandler{datastore: department}
}

func (e DepartmentHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	var department Entities.Department

	body, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(body, &department)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("invalid body"))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resp, err := e.datastore.Create(department)
	if err != nil {
		_, _ = w.Write([]byte("could not create department"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
