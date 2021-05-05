package gofr

import (
	"strings"

	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/log"
)

type Gofr struct {
	datastore.DataStore

	cmd         *cmdApp
	Config      Config
	Server      *server
	TemplateDir string
	Logger      log.Logger

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
	k.addRoute("GET", path, handler)
}

func (k *Gofr) PUT(path string, handler Handler) {
	k.addRoute("PUT", path, handler)
}

func (k *Gofr) POST(path string, handler Handler) {
	k.addRoute("POST", path, handler)
}

func (k *Gofr) DELETE(path string, handler Handler) {
	k.addRoute("DELETE", path, handler)
}

func (k *Gofr) PATCH(path string, handler Handler) {
	k.addRoute("PATCH", path, handler)
}

func (k *Gofr) EnableSwaggerUI() {
	k.addRoute("GET", "/swagger", SwaggerUIHandler)
	k.addRoute("GET", "/swagger/{name}", SwaggerUIHandler)
}
