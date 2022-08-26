package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"EmployeeDepartment/errorsHandler"
)

func SetStatusCode(w http.ResponseWriter, method string, data interface{}, err error) {
	switch err.(type) {
	case errorsHandler.AlreadyExists:
		w.WriteHeader(http.StatusConflict)
		WriteToBody(w, err)
	case *errorsHandler.InvalidDetails:
		w.WriteHeader(http.StatusBadRequest)
		WriteToBody(w, err)
	case *errorsHandler.IDNotFound:
		w.WriteHeader(http.StatusNotFound)
		WriteToBody(w, err)
	case *errorsHandler.NoData:
		WriteToBody(w, err)
	case nil:
		WriteSuccessResponse(method, w, data)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		WriteToBody(w, err)
	}
}

func WriteSuccessResponse(method string, w http.ResponseWriter, data interface{}) {
	switch method {
	case http.MethodPost:
		writeResponseBody(w, http.StatusCreated, data)
	case http.MethodGet:
		writeResponseBody(w, http.StatusOK, data)
	case http.MethodPut:
		writeResponseBody(w, http.StatusOK, data)
	case http.MethodDelete:
		writeResponseBody(w, http.StatusNoContent, data)
	}
}

func writeResponseBody(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data == nil {
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println("error in writing response")
		return
	}
}

func WriteToBody(res http.ResponseWriter, err error) {
	_, err = res.Write([]byte(err.Error()))
	if err != nil {
		log.Println(err)
	}
}
