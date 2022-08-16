package employee

import (
	"EmployeeDepartment/entities"
	"EmployeeDepartment/service"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type EmployeeHandler struct {
	validate service.Employee
}

func New(emp service.Employee) EmployeeHandler {
	return EmployeeHandler{validate: emp}
}

func (e EmployeeHandler) PostHandler(res http.ResponseWriter, req *http.Request) {
	var employee entities.Employee
	body, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(body, &employee)
	if err != nil {
		_, _ = res.Write([]byte("Unmarshal Error"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err1 := e.validate.Create(employee)
	if err1 != nil {
		_, _ = res.Write([]byte("Invalid Body"))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	res.Write(body)
}

func (e EmployeeHandler) GetHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[10:]

	if len(id) != 36 {
		_, _ = res.Write([]byte("Invalid Id"))
		return
	}
	uid := uuid.MustParse(id)
	resp, err := e.validate.Read(uid)
	body, _ := json.Marshal(resp)
	if err != nil {
		_, _ = res.Write([]byte("Id not found"))
		return
	}
	res.Write(body)

}

func (e EmployeeHandler) PutHandler(res http.ResponseWriter, req *http.Request) {
	parameter := mux.Vars(req)
	id := parameter["id"]
	var employee entities.Employee
	reader, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(reader, &employee)
	if err != nil {
		res.Write([]byte("Unmarshall error"))
		return
	}
	uid := uuid.MustParse(id)
	resp, err := e.validate.Update(uid, employee)
	if err != nil {
		_, _ = res.Write([]byte("Id not found"))
		return
	}
	body, _ := json.Marshal(resp)
	res.Write(body)
}

func (e EmployeeHandler) DeleteHandler(res http.ResponseWriter, req *http.Request) {
	id := uuid.MustParse(req.URL.Path[10:])
	resp, err := e.validate.Delete(id)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(resp)

}

func (e EmployeeHandler) GetAll(res http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	includeDepartment := req.URL.Query().Get("includeDepartment")
	fmt.Println(name, includeDepartment)
	b, err := strconv.ParseBool(includeDepartment)
	resp, err := e.validate.ReadAll(name, b)
	if err != nil {
		res.Write([]byte("Unmarshal Error"))
		return
	}
	data, err := json.Marshal(resp)
	res.Write(data)
}
