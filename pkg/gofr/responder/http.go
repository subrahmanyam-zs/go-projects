package responder

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
)

type responseType int

const (
	JSON responseType = iota
	XML
	TEXT
)

type HTTP struct {
	path          string
	method        string
	w             http.ResponseWriter
	resType       responseType
	correlationID string
}

// NewContextualResponder creates a HTTP responder which gives JSON/XML response based on context
func NewContextualResponder(w http.ResponseWriter, r *http.Request) Responder {
	route := mux.CurrentRoute(r)

	var path string
	if route != nil {
		path, _ = route.GetPathTemplate()
		// remove the trailing slash
		path = strings.TrimSuffix(path, "/")
	}

	responder := &HTTP{
		w:             w,
		method:        r.Method,
		path:          path,
		correlationID: middleware.GetCorrelationID(r),
	}

	// set correlation id in response
	w.Header().Set("X-Correlation-Id", responder.correlationID)

	cType := r.Header.Get("Content-type")
	switch cType {
	case "text/xml", "application/xml":
		responder.resType = XML
	case "text/plain":
		responder.resType = TEXT
	default:
		responder.resType = JSON
	}

	return responder
}

func (h HTTP) Respond(data interface{}, err error) {
	// if template is returned then everything is dictated by template
	if d, ok := data.(template.Template); ok {
		var b []byte
		b, err = d.Render()

		if err != nil {
			h.processTemplateError(err)
			return
		}

		h.w.Header().Set("Content-Type", d.ContentType())
		h.w.WriteHeader(http.StatusOK)
		_, _ = h.w.Write(b)

		return
	}

	if f, ok := data.(template.File); ok {
		h.w.Header().Set("Content-Type", f.ContentType)
		_, _ = h.w.Write(f.Content)

		return
	}

	var (
		response   interface{}
		statusCode int
	)

	res, okay := data.(*types.Response)
	if res == nil {
		res = &types.Response{}
	}

	if !okay {
		response = data
		statusCode = getStatusCode(h.method, data, err)
	} else {
		response = getResponse(res, err)
		statusCode = getStatusCode(h.method, res.Data, err)
	}

	h.processResponse(statusCode, response)
}

func (h HTTP) processResponse(statusCode int, response interface{}) {
	switch h.resType {
	case JSON:
		h.w.Header().Set("Content-type", "application/json")
		h.w.WriteHeader(statusCode)

		if response != nil {
			_ = json.NewEncoder(h.w).Encode(response)
		}

	case XML:
		h.w.Header().Set("Content-type", "application/xml")
		h.w.WriteHeader(statusCode)

		if response != nil {
			_ = xml.NewEncoder(h.w).Encode(response)
		}
	case TEXT:
		h.w.Header().Set("Content-type", "text/plain")
		h.w.WriteHeader(statusCode)

		if response != nil {
			_, _ = h.w.Write([]byte(fmt.Sprintf("%s", response)))
		}
	}
}

func getStatusCode(method string, data interface{}, err error) int {
	statusCode := 200

	if err == nil {
		if method == http.MethodPost {
			statusCode = 201
		} else if method == http.MethodDelete {
			statusCode = 204
		}

		return statusCode
	}

	if e, ok := err.(errors.MultipleErrors); ok {
		if data != nil {
			return http.StatusPartialContent
		}

		statusCode = e.StatusCode
		if e.StatusCode == 0 {
			statusCode = http.StatusInternalServerError
		}

		return statusCode
	}

	return statusCode
}

func getResponse(res *types.Response, err error) interface{} {
	// Response error should be of MultipleErrors type
	em, ok := err.(errors.MultipleErrors)

	if res == nil || !ok {
		return res
	}
	// If data and error both are present (Partial Content)
	if res.Data != nil {
		dataMap := make(map[string]interface{})
		dataMap["errors"] = em.Errors

		b := new(bytes.Buffer)
		_ = json.NewEncoder(b).Encode(res.Data)
		_ = json.NewDecoder(b).Decode(&dataMap)

		// To handle the case of interface having nullable type and its value is nil
		if dataMap == nil {
			res.Data = nil // Ensuring response is not partial content
			return em
		}

		return types.Response{Data: dataMap, Meta: res.Meta}
	}

	// error is present but only status code is needed to be set and no body
	if em.Error() == "" {
		return nil
	}
	// error is set and returned in the body
	return em
}

func (h HTTP) processTemplateError(err error) {
	errorData := &errors.Response{}
	errorData.Reason = err.Error()

	switch err.(type) {
	case errors.FileNotFound:
		errorData.Code = "File Not Found"
		errorData.StatusCode = http.StatusNotFound
	default:
		errorData.StatusCode = http.StatusInternalServerError
		errorData.Code = "Internal Server Error"
		// pushing error type to prometheus
		middleware.ErrorTypesStats.With(prometheus.Labels{"type": "UnknownError", "path": h.path, "method": h.method}).Inc()
	}

	errMultiple := errors.MultipleErrors{
		StatusCode: errorData.StatusCode,
		Errors:     []error{errorData},
	}

	h.processResponse(errMultiple.StatusCode, errMultiple)
}
