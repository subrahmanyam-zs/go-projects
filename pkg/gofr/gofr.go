package gofr

import (
	"net/http"
	"strings"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr/metrics"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/notifier"
)

type Gofr struct {
	datastore.DataStore

	cmd         *cmdApp
	Config      Config
	Server      *server
	TemplateDir string
	Logger      log.Logger
	Metric      metrics.Metric
	Notifier    notifier.Notifier

	ResourceMap          map[string][]string
	ResourceCustomShapes map[string][]string

	ServiceHealth  []HealthCheck
	DatabaseHealth []HealthCheck
}

func (k *Gofr) Start() {
	if k.cmd != nil {
		k.cmd.Start(k.Logger)
	} else {
		k.Server.Start(k.Logger)
	}
}

func (k *Gofr) addRoute(method, path string, handler Handler) {
	if k.cmd != nil {
		k.cmd.Router.AddRoute(path, handler) // Ignoring method in CMD App.
	} else {
		if path != "/" {
			path = strings.TrimSuffix(path, "/")
			k.Server.Router.Route(method, path+"/", handler)
		}
		k.Server.Router.Route(method, path, handler)
	}
}

func (k *Gofr) GET(path string, handler Handler) {
	k.addRoute(http.MethodGet, path, handler)
}

func (k *Gofr) PUT(path string, handler Handler) {
	k.addRoute(http.MethodPut, path, handler)
}

func (k *Gofr) POST(path string, handler Handler) {
	k.addRoute(http.MethodPost, path, handler)
}

func (k *Gofr) DELETE(path string, handler Handler) {
	k.addRoute(http.MethodDelete, path, handler)
}

func (k *Gofr) PATCH(path string, handler Handler) {
	k.addRoute(http.MethodPatch, path, handler)
}

// Deprecated: EnableSwaggerUI is deprecated. Auto enabled swagger-endpoints.
func (k *Gofr) EnableSwaggerUI() {
	k.addRoute(http.MethodGet, "/swagger", SwaggerUIHandler)
	k.addRoute(http.MethodGet, "/swagger/{name}", SwaggerUIHandler)

	k.Logger.Warn("Usage of EnableSwaggerUI is deprecated. Swagger Endpoints are auto-enabled")
}
