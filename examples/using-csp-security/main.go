package main

import (
	"github-lvs.corpzone.internalzone.com/mcafee/cnsr-gofr-csp-auth/validator"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	opts := validator.Options{
		Keys: map[string]string{
			app.Config.Get("CSP_APP_KEY_CATALOG"):     app.Config.Get("CSP_SHARED_KEY_CATALOG"),
			app.Config.Get("CSP_APP_KEY_USER_ASSETS"): app.Config.Get("CSP_SHARED_KEY_USER_ASSETS"),
		},
	}

	app.Server.UseMiddleware(validator.CSPAuth(app.Logger, opts))

	app.GET("/hello", func(ctx *gofr.Context) (interface{}, error) {
		name := ctx.Param("name")
		return "Hello " + name, nil
	})

	app.POST("/greet", func(ctx *gofr.Context) (interface{}, error) {
		var data struct {
			Name string `json:"name"`
		}

		err := ctx.Bind(&data)
		if err != nil {
			return nil, &errors.InvalidParam{Param: []string{"body"}}
		}

		return "Hello " + data.Name, nil
	})

	app.Start()
}
