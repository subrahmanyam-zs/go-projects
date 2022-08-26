package department

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"EmployeeDepartment/entities"
	"EmployeeDepartment/errorsHandler"
	"EmployeeDepartment/handler"
	"EmployeeDepartment/service"
)

type Handler struct {
	service service.Department
}

func New(department service.Department) Handler {
	return Handler{service: department}
}

func (h Handler) PostHandler(w http.ResponseWriter, req *http.Request) {
	var department entities.Department

	reader, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(reader, &department)
	if err != nil {
		handler.SetStatusCode(w, req.Method, nil, &errorsHandler.InvalidDetails{Msg: "body"})
		return
	}

	resp, err := h.service.Create(req.Context(), department)
	if err != nil {
		handler.SetStatusCode(w, req.Method, nil, err)
		return
	}

	handler.SetStatusCode(w, req.Method, resp, nil)
}

func (h Handler) PutHandler(res http.ResponseWriter, req *http.Request) {
	var department entities.Department

	sid := req.URL.Path[12:]

	id, err := strconv.Atoi(sid)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, &errorsHandler.InvalidDetails{Msg: "id"})
		return
	}

	dept, err := h.service.GetDepartment(req.Context(), id)
	if err != nil {
		handler.SetStatusCode(res, req.Method, dept, err)
		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(reqBody, &department)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, &errorsHandler.InvalidDetails{Msg: "body"})
		return
	}

	resp, err := h.service.Update(req.Context(), id, department)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, err)
		return
	}

	handler.SetStatusCode(res, req.Method, resp, nil)
}

func (h Handler) DeleteHandler(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Path[12:])
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, &errorsHandler.InvalidDetails{Msg: "id"})
		return
	}

	resp, err := h.service.Delete(req.Context(), id)
	if err != nil {
		handler.SetStatusCode(res, req.Method, nil, err)
		return
	}

	handler.SetStatusCode(res, req.Method, resp, nil)
}
