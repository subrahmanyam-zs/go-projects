package Department

import (
	"EmployeeDepartment/Handler/Entities"
	"EmployeeDepartment/Service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type DepartmentHandler struct {
	datastore Service.Department
}

func New(department Service.Department) DepartmentHandler {
	return DepartmentHandler{datastore: department}
}

func (e DepartmentHandler) PostHandler(w http.ResponseWriter, req *http.Request) {
	var department Entities.Department

	reader, _ := io.ReadAll(req.Body)

	err := json.Unmarshal(reader, &department)
	if err != nil {
		//fmt.Println(err)
		_, _ = w.Write([]byte("Unmarshal Error"))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resp, err := e.datastore.Create(department)
	if err != nil {
		_, _ = w.Write([]byte("Invalid id"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (e DepartmentHandler) PutHandler(res http.ResponseWriter, req *http.Request) {
	var department Entities.Department
	sid := req.URL.Path[12:]
	reqBody, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(reqBody, &department)
	if err != nil {
		res.Write([]byte("Unmarshal Error"))
		return
	}
	id, _ := strconv.Atoi(sid)
	resp, err := e.datastore.Update(id, department)
	fmt.Println(id)
	if err != nil {
		res.Write([]byte("Id not found"))
		return
	}
	body, _ := json.Marshal(resp)
	res.Write(body)
}

func (d DepartmentHandler) DeleteHandler(res http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(req.URL.Path[12:])
	resp, err := d.datastore.Delete(id)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(resp)

}
