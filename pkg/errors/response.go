package errors

import (
	"fmt"
)

type Response struct {
	StatusCode int         `json:"-"`
	Code       string      `json:"code"`
	Reason     string      `json:"reason"`
	ResourceID string      `json:"resourceId,omitempty"`
	Detail     interface{} `json:"detail,omitempty"`
	Path       string      `json:"path,omitempty"`
	RootCauses []RootCause `json:"rootCauses,omitempty"`
	DateTime   `json:"datetime"`
}

type RootCause map[string]interface{}

type DateTime struct {
	Value    string `json:"value"`
	TimeZone string `json:"timezone"`
}

func (k *Response) Error() string {
	return fmt.Sprint(k.Reason)
}

type Error string

func (e Error) Error() string {
	return string(e)
}
