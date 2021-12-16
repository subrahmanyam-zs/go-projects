package gofr

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func healthCheckHandlerServer(ctx *Context, port int) *http.Server {
	r := mux.NewRouter()
	healthResp, err := HealthHandler(ctx)

	r.HandleFunc("/.well-known/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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

	return &http.Server{Addr: ":" + strconv.Itoa(port), Handler: r}
}
