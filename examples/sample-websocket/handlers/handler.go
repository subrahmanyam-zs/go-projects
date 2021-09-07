package handlers

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
)

func WSHandler(c *gofr.Context) (interface{}, error) {
	var (
		mt      int
		message []byte
		err     error
	)

	if c.WebSocketConnection != nil {
		for {
			mt, message, err = c.WebSocketConnection.ReadMessage()
			if err != nil {
				c.Logger.Error("read:", err)
				break
			}

			c.Logger.Logf("recv: %v", string(message))

			err = c.WebSocketConnection.WriteMessage(mt, message)
			if err != nil {
				c.Logger.Error("write:", err)
				break
			}
		}
	}

	return nil, err
}

func HomeHandler(c *gofr.Context) (interface{}, error) {
	return template.Template{File: "home.html", Type: template.HTML}, nil
}
