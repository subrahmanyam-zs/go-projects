package responder

import (
	"fmt"

	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
)

type CMD struct{}

func (c *CMD) Respond(data interface{}, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

	if d, ok := data.(template.Template); ok {
		var b []byte
		b, err = d.Render()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))

		return
	}

	if f, ok := data.(template.File); ok {
		fmt.Println(string(f.Content))

		return
	}

	fmt.Println(data)
}
