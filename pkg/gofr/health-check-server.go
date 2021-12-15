package gofr

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func healthCheckHandlerServer(ctx *Context, port int) *http.Server {
	mux := http.NewServeMux()
	healthResp, err := HealthHandler(ctx)

	mux.HandleFunc(defaultHealthCheckRoute, func(w http.ResponseWriter, r *http.Request) {
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

	srv := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: mux}

	return srv
}
