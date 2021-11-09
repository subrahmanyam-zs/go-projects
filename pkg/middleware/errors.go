package middleware

import "net/http"

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidToken    = Error("invalid_token")
	ErrInvalidRequest  = Error("invalid_request")
	ErrServiceDown     = Error("service_unavailable")
	ErrInvalidHeader   = Error("invalid_header")
	ErrMissingHeader   = Error("missing_header")
	ErrUnauthorised    = Error("missing_permission")
	ErrUnauthenticated = Error("failed_auth")
)

func GetDescription(err error) (description string, statusCode int) {
	var authErr = "Authorization error"

	switch err {
	case ErrInvalidToken:
		description = "The access token is invalid or has expired"
		statusCode = http.StatusUnauthorized
	case ErrInvalidRequest:
		description = "The access token is missing"
		statusCode = http.StatusUnauthorized
	case ErrServiceDown:
		description = "Unable to validate the token"
		statusCode = http.StatusInternalServerError
	case ErrMissingHeader:
		description = "Missing Authorization header"
		statusCode = http.StatusUnauthorized
	case ErrInvalidHeader:
		description = "Invalid Authorization header"
		statusCode = http.StatusBadRequest
	case ErrUnauthorised:
		description = authErr
		statusCode = http.StatusForbidden
	case ErrUnauthenticated:
		description = authErr
		statusCode = http.StatusUnauthorized
	}

	return description, statusCode
}
