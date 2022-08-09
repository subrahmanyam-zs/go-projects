package Employee

import (
	"EmployeeDepartment/Handler/Entities"
	"EmployeeDepartment/Store"
	"encoding/json"
	"github.com/google/uuid"
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
		_, _ = res.Write([]byte("Unmarshal Error"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := e.datastore.Create(employee)
	if err != nil {
		_, _ = res.Write([]byte("Invalid Body"))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, _ = json.Marshal(resp)
	_, _ = res.Write(body)
}

func (e EmployeeHandler) GetHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[10:]

	if len(id) != 36 {
		_, _ = res.Write([]byte("Invalid Id"))
		return
	}
	uid := uuid.MustParse(id)

	resp, err := e.datastore.Read(uid)
	body, _ := json.Marshal(resp)
	if err != nil {
		_, _ = res.Write([]byte("Id not found"))
		return
	}
	_, _ = res.Write(body)

}

func (e EmployeeHandler) PutHandler(res http.ResponseWriter, req *http.Request) {
	var employee Entities.Employee
	id := req.URL.Path[10:]
	reader, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(reader, &employee)
	if err != nil {
		res.Write([]byte("Unmarshall error"))
		return
	}

	if len(id) != 36 {
		res.Write([]byte("Invalid id"))
		return
	}
	uid := uuid.MustParse(id)
	resp, err := e.datastore.Update(uid, employee)
	body, _ := json.Marshal(resp)
	if err != nil {
		_, _ = res.Write(body)
		return
	}
	res.Write(reader)
}
