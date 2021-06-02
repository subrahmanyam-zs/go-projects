package handler

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
)

type Person struct {
	Username string
	Password string
}

// HomeHandler renders home.html file
// Uses server push to give out the app.css file required by the home.html file
func HomeHandler(c *gofr.Context) (interface{}, error) {
	if c.ServerPush != nil {
		c.Logger.Info("gofr.png required by home.html pushed using server push")

		err := c.ServerPush.Push("/static/gofr.png", nil)
		if err != nil {
			return nil, err
		}
	}

	return template.Template{File: "home.html", Type: template.HTML}, nil
}

// ServeStatic Renders the file present in /static folder
func ServeStatic(c *gofr.Context) (interface{}, error) {
	fileName := c.PathParam("name")

	return template.Template{File: fileName}, nil
}
