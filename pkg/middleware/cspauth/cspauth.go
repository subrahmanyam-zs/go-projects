package cspauth

import (
	"encoding/json"
	"net/http"
	"strings"

	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
)

// CSPAuth middleware provides authentication using CSP
func CSPAuth(logger log.Logger, sharedKey string) func(inner http.Handler) http.Handler {
	cache := NewCache()
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if middleware.ExemptPath(req) {
				inner.ServeHTTP(w, req)

				return
			}
			opts := getOptions(sharedKey, req)

			csp, err := New(logger, opts, cache)
			if err != nil {
				e := middleware.FetchErrResponseWithCode(http.StatusBadRequest, "Invalid CSP Auth Options", err.Error())
				middleware.ErrorResponse(w, req, logger, *e)

				return
			}

			err = csp.Validate(logger, req)
			if err != nil {
				description, statusCode := middleware.GetDescription(err)
				e := middleware.FetchErrResponseWithCode(statusCode, description, err.Error())
				middleware.ErrorResponse(w, req, logger, *e)

				return
			}

			inner.ServeHTTP(w, req)
		})
	}
}

func getOptions(sharedKey string, req *http.Request) *Options {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		// ip address from the RemoteAddr
		ip = strings.Split(req.RemoteAddr, ":")[0]
	}

	return &Options{
		MachineName: req.Header.Get("User-Agent"),
		IPAddress:   ip,
		AppKey:      req.Header.Get(appKeyHeader),
		AppID:       req.Header.Get(clientIDHeader),
		SharedKey:   sharedKey,
	}
}

type cspAuthJSON struct {
	IPAddress     string `json:"IPAddress"`
	MachineName   string `json:"MachineName"`
	RequestDate   string `json:"RequestDate"`
	HTTPMethod    string `json:"HttpMethod"`
	UUID          string `json:"MsgUniqueId"`
	ClientID      string `json:"ClientId"`
	SignatureHash string `json:"SignatureHash"`
}

// Validate the csp auth headers in given request
func (c *CSP) Validate(logger log.Logger, r *http.Request) error {
	authContext, err := c.getAuthContext(logger, r.Header.Get(authContextHeader))
	if err != nil {
		return middleware.ErrInvalidAuthContext
	}

	var authJSON cspAuthJSON

	err = json.Unmarshal(authContext, &authJSON)
	if err != nil {
		logger.Errorf("error while unmarshalling csp auth json, %v", err)

		return middleware.ErrInvalidAuthContext
	}

	bodyHash := getBodyHash(r)

	// generate data and its signature for validation
	dataForSigValidation := authJSON.MachineName + authJSON.RequestDate + authJSON.IPAddress + c.options.AppKey +
		c.options.SharedKey + authJSON.HTTPMethod + authJSON.UUID + authJSON.ClientID + bodyHash

	computedSignature := base64Encode([]byte(hexEncode(sha256Hash([]byte(dataForSigValidation)))))

	if computedSignature != authJSON.SignatureHash {
		return middleware.ErrInvalidAuthContext
	}

	return nil
}

func (c *CSP) getAuthContext(logger log.Logger, authContextHeader string) ([]byte, error) {
	b64DecodeRandom, err := base64Decode(authContextHeader)
	if err != nil {
		logger.Errorf("error while base64 decoding auth context, %v", err)

		return nil, err
	}

	if len(b64DecodeRandom) <= lenRandomChars {
		return nil, middleware.ErrInvalidAuthContext
	}
	// remove random string from auth context
	authContextToDecode := b64DecodeRandom[:len(b64DecodeRandom)-lenRandomChars]

	authContextToDecrypt, err := base64Decode(string(authContextToDecode))
	if err != nil {
		logger.Errorf("error while base64 decoding auth context without random chars, %v", err)

		return nil, err
	}

	// decrypt auth context using encryption key and initial vector
	decryptedAuthContext, err := decryptData(authContextToDecrypt, c.encryptionKey, c.iv)
	if err != nil {
		logger.Errorf("error occurred while decrypting auth context, %v", err)

		return nil, err
	}

	return decryptedAuthContext, nil
}
