package gofr

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func healthCheckHandlerServer(ctx *Context, port int, route string) *http.Server {
	mux := http.NewServeMux()
	healthResp, err := HealthHandler(ctx)

	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			ctx.Logger.Error(err)

			data, _ := json.Marshal(err)

			_, err := w.Write(data)
			if err != nil {
				ctx.Logger.Error(err)

				return
			}

			return
		}

		data, _ := json.Marshal(healthResp)

		_, err := w.Write(data)
		if err != nil {
			ctx.Logger.Error(err)

			return
		}
	})

	srv := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: mux}

	return srv
}
