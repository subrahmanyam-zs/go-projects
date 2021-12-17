package gofr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
)

func healthCheckHandlerServer(ctx *Context, port int) *http.Server {
	r := mux.NewRouter()

	r.Use(validateRoutes(ctx.Logger))

	r.HandleFunc("/.well-known/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		healthResp, err := HealthHandler(ctx)
		if err != nil {
			ctx.Logger.Error(err)

			data, err := json.Marshal(err)
			if err != nil {
				ctx.Logger.Error(err)

				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, _ = w.Write(data)

			return
		}

		data, err := json.Marshal(healthResp)
		if err != nil {
			ctx.Logger.Error(err)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, _ = w.Write(data)
	})

	// handles 404
	r.NotFoundHandler = r.NewRoute().HandlerFunc(http.NotFound).GetHandler()

	return &http.Server{Addr: ":" + strconv.Itoa(port), Handler: r}
}

func validateRoutes(log log.Logger) func(http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			path := fmt.Sprintf("%s", req.URL.Path)

			if !strings.Contains("/.well-known/health-check", path) {
				err := middleware.FetchErrResponseWithCode(http.StatusNotFound,
					fmt.Sprintf("Route %v not found", req.URL), "Invalid Route")

				middleware.ErrorResponse(w, req, log, *err)

				return
			}

			if req.Method != http.MethodGet {
				err := middleware.FetchErrResponseWithCode(http.StatusMethodNotAllowed,
					fmt.Sprintf("%v method not allowed for Route %v", req.Method, req.URL), "Invalid Method")

				middleware.ErrorResponse(w, req, log, *err)

				return
			}

			inner.ServeHTTP(w, req)
		})
	}
}
