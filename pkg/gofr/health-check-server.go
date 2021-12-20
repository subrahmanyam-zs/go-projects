package gofr

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"go.opencensus.io/trace"

	"github.com/gorilla/mux"

	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
)

func healthCheckHandlerServer(app *cmdApp) *http.Server {
	r := mux.NewRouter()

	r.Use(validateRoutes(app.context.Logger), app.contextInjector)

	r.HandleFunc(defaultHealthCheckRoute, app.healthCheck)

	// handles 404
	r.NotFoundHandler = r.NewRoute().HandlerFunc(http.NotFound).GetHandler()

	return &http.Server{Addr: ":" + strconv.Itoa(app.healthCheckSvr.port), Handler: r}
}

// healthCheck calls HealthHandler and returns the response for healthCheck.
func (app *cmdApp) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data []byte

	healthResp, err := HealthHandler(app.context)
	if err != nil {
		app.context.Logger.Error(err)

		data, err = json.Marshal(err)
		if err != nil {
			app.context.Logger.Error(err)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, _ = w.Write(data)

		return
	}

	data, err = json.Marshal(healthResp)
	if err != nil {
		app.context.Logger.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, _ = w.Write(data)
}

func validateRoutes(l log.Logger) func(http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			path := req.URL.Path

			if !strings.Contains(defaultHealthCheckRoute, path) {
				err := middleware.FetchErrResponseWithCode(http.StatusNotFound,
					fmt.Sprintf("Route %v not found", req.URL), "Invalid Route")

				middleware.ErrorResponse(w, req, l, *err)

				return
			}

			if req.Method != http.MethodGet {
				err := middleware.FetchErrResponseWithCode(http.StatusMethodNotAllowed,
					fmt.Sprintf("%v method not allowed for Route %v", req.Method, req.URL), "Invalid Method")

				middleware.ErrorResponse(w, req, l, *err)

				return
			}

			inner.ServeHTTP(w, req)
		})
	}
}

// nolint:dupl // contextInjector is used only for health-check in GOFR CMD
func (app *cmdApp) contextInjector(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := app.contextPool.Get().(*Context)
		c.reset(responder.NewContextualResponder(w, r), request.NewHTTPRequest(r))
		*r = *r.WithContext(ctx.WithValue(r.Context(), appData, &sync.Map{}))
		c.Context = r.Context()
		*r = *r.WithContext(ctx.WithValue(c.Context, gofrContextkey, c))

		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = r.Header.Get("X-B3-TraceID")
		}
		if correlationID == "" {
			correlationID = trace.FromContext(r.Context()).SpanContext().TraceID.String()
		}

		c.Logger = log.NewCorrelationLogger(correlationID)

		inner.ServeHTTP(w, r)

		app.contextPool.Put(c)
	})
}
