package Employee

import (
	"EmployeeDepartment/Handler/Entities"
	"EmployeeDepartment/Store"
	"encoding/json"
	"io"
	"net/http"
)

type EmployeeHandler struct {
	datastore Store.Employee
}

func New(emp Store.Employee) EmployeeHandler {
	return EmployeeHandler{datastore: emp}
}
func (e EmployeeHandler) PostHandler(res http.ResponseWriter, req *http.Request) {
	var employee Entities.Employee
	body, _ := io.ReadAll(req.Body)

	err := json.Unmarshal(body, &employee)

	if err != nil {
		_, _ = res.Write([]byte("invalid body"))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	resp, err := e.datastore.Create(employee)
	if err != nil {
		_, _ = res.Write([]byte("could not create employee"))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, _ = json.Marshal(resp)
	_, _ = res.Write(body)
}
