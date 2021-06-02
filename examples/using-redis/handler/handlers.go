package handler

import (
	"developer.zopsmart.com/go/gofr/examples/using-redis/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Config struct {
	store store.Store
}

func New(c store.Store) *Config {
	return &Config{
		store: c,
	}
}

// SetKey is a handler function of type gofr.Handler, it sets keys
func (m Config) SetKey(c *gofr.Context) (interface{}, error) {
	input := make(map[string]string)

	if err := c.Bind(&input); err != nil {
		return nil, invalidBodyErr{}
	}

	for key, value := range input {
		if err := m.store.Set(c, key, value, 0); err != nil {
			c.Logger.Error("got error: ", err)
			return nil, invalidInputErr{}
		}
	}

	return "Successful", nil
}

// GetKey is a handler function of type gofr.Handler, it fetches keys
func (m Config) GetKey(c *gofr.Context) (interface{}, error) {
	// fetch the path parameter as specified in the route
	key := c.PathParam("key")
	if key == "" {
		return nil, errors.MissingParam{Param: []string{"key"}}
	}

	value, err := m.store.Get(c, key)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]string)
	resp[key] = value

	return resp, nil
}

// DeleteKey is a handler function of type gofr.Handler, it deletes keys
func (m Config) DeleteKey(c *gofr.Context) (interface{}, error) {
	// fetch the path parameter as specified in the route
	key := c.PathParam("key")
	if key == "" {
		return nil, errors.MissingParam{Param: []string{"key"}}
	}

	if err := m.store.Delete(c, key); err != nil {
		c.Logger.Errorf("err: ", err)
		return nil, deleteErr{}
	}

	return "Deleted successfully", nil
}

type (
	deleteErr       struct{}
	invalidInputErr struct{}
	invalidBodyErr  struct{}
)

func (d deleteErr) Error() string {
	return "error: failed to delete"
}

func (i invalidInputErr) Error() string {
	return "error: invalid input"
}

func (i invalidBodyErr) Error() string {
	return "error: invalid body"
}
