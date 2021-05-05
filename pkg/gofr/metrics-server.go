package gofr

import (
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zopsmart/gofr/pkg/log"
)

func metricsServer(logger log.Logger, port int, route string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(route, promhttp.Handler())

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	var srv = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	logger.Infof("Starting metrics server at :%v", port)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Errorf("error in metrics server %v", err)
		}
	}()

	return srv
}
