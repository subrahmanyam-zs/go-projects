package employee

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"EmployeeDepartment/entities"
	"EmployeeDepartment/errorsHandler"
	"EmployeeDepartment/handler"
	"EmployeeDepartment/service"
	"EmployeeDepartment/store"
)

type Handler struct {
	service service.Employee
}

func New(emp service.Employee) Handler {
	return Handler{service: emp}
}

func (h Handler) PostHandler(res http.ResponseWriter, req *http.Request) {
	var employee entities.Employee

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(body, &employee)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, &errorsHandler.InvalidDetails{Msg: "body"})
		return
	}

	resp, err := h.service.Create(req.Context(), &employee)
	if err != nil {
		handler.SetStatusCode(res, req.Method, resp, err)
		return
	}

	handler.SetStatusCode(res, req.Method, resp, err)
}

func (h Handler) GetHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[10:]

	thirtySix := 36
	if len(id) != thirtySix {
		handler.WriteToBody(res, &errorsHandler.InvalidDetails{Msg: "ID"})
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, &errorsHandler.InvalidDetails{Msg: "id"})
		return
	}

	resp, err := h.service.Read(req.Context(), uid)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, err)
		return
	}

	handler.SetStatusCode(res, req.Method, resp, err)
}

func (h Handler) PutHandler(res http.ResponseWriter, req *http.Request) {
	var employee entities.Employee

	id := req.URL.Path[10:]

	uid, err := uuid.Parse(id)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, &errorsHandler.InvalidDetails{Msg: "id"})
		return
	}

	empDept, err := h.service.Read(req.Context(), uid)
	if err != nil {
		handler.SetStatusCode(res, req.Method, empDept, &errorsHandler.IDNotFound{Msg: "Id not found"})
		return
	}

	reader, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(reader, &employee)
	if err != nil {
		handler.SetStatusCode(res, req.Method, employee, &errorsHandler.InvalidDetails{Msg: "body"})
		return
	}

	resp, err := h.service.Update(req.Context(), uid, &employee)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, err)
		return
	}

	handler.SetStatusCode(res, req.Method, resp, err)
}

func (h Handler) DeleteHandler(res http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.URL.Path[10:])
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, &errorsHandler.InvalidDetails{Msg: "id"})
		return
	}

	resp, err := h.service.Delete(req.Context(), id)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, err)
		return
	}

	handler.SetStatusCode(res, req.Method, resp, err)
}

func (h Handler) GetAll(res http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	includeDepartment := req.URL.Query().Get("includeDepartment")

	b, err := strconv.ParseBool(includeDepartment)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, &errorsHandler.InvalidDetails{Msg: "value for includeDept"})
		return
	}

	resp, err := h.service.ReadAll(store.Parameters{Ctx: context.TODO(), Name: name, IncludeDepartment: b})
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, err)
		return
	}

	handler.SetStatusCode(res, req.Method, resp, err)
}
