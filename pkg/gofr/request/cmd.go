package request

import (
	"flag"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type CMD struct {
	params map[string]string
}

func NewCMDRequest() Request {
	c := &CMD{}

	flag.Parse()
	args := flag.Args()
	c.parseArgs(args)

	return c
}

func (c *CMD) parseArgs(args []string) {
	c.params = make(map[string]string)

	const (
		argsLen1 = 1
		argsLen2 = 2
	)

	for _, arg := range args {
		if arg[0] != '-' {
			continue
		}

		a := arg[1:]

		switch values := strings.Split(a, "="); len(values) {
		case argsLen1:
			// Support -t -a etc.
			c.params[values[0]] = "true"
		case argsLen2:
			// Support -a=b
			c.params[values[0]] = values[1]
		}
	}
}

func (c *CMD) Param(key string) string {
	return c.params[key]
}

func (c *CMD) PathParam(key string) string {
	return c.params[key]
}

// Header is the same as Param()
func (c *CMD) Header(key string) string {
	return c.Param(key)
}

func (c *CMD) Params() map[string]string {
	return c.params
}

func (c *CMD) Request() *http.Request {
	return nil
}

//nolint:gocognit // Reducing cognitive complexity will make it harder to read.
func (c *CMD) Bind(i interface{}) error {
	// pointer to struct - addressable
	ps := reflect.ValueOf(i)
	// struct
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		for k, v := range c.params {
			f := s.FieldByName(k)
			// A Value can be changed only if it is addressable and not unexported struct field
			if f.IsValid() && f.CanSet() {
				// nolint:exhaustive // no need to add other cases
				switch f.Kind() {
				case reflect.String:
					f.SetString(v)
				case reflect.Bool:
					if v == "true" {
						f.SetBool(true)
					}
				case reflect.Int:
					n, _ := strconv.Atoi(v)
					f.SetInt(int64(n))
				}
			}
		}
	}

	return nil
}

func (c *CMD) BindStrict(i interface{}) error {
	return c.Bind(i)
}

// GetClaims returns nil claim for every request
func (c *CMD) GetClaims() map[string]interface{} {
	return nil
}

// GetClaim returns nil claim value for every request
func (c *CMD) GetClaim(claimKey string) interface{} {
	return nil
}
