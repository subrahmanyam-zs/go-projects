package cspauth

import (
	"encoding/json"
	"net/http"

	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
)

// CSPAuth middleware provides authentication using CSP
func CSPAuth(logger log.Logger, sharedKey string) func(inner http.Handler) http.Handler {
	csp := New(sharedKey)

	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if middleware.ExemptPath(req) {
				inner.ServeHTTP(w, req)

				return
			}

			appKey, err := csp.getAppKey(req)
			if err != nil {
				e := middleware.FetchErrResponseWithCode(http.StatusBadRequest, "Invalid CSP Auth Options", err.Error())
				middleware.ErrorResponse(w, req, logger, *e)

				return
			}

			csp.Set(appKey, csp.sharedKey)

			err = csp.Validate(logger, req, appKey)
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

func (c *CSP) getAppKey(req *http.Request) (string, error) {
	appKey := req.Header.Get(appKeyHeader)
	if len(appKey) < minLenAppKey {
		return "", ErrEmptyAppKey
	}

	return appKey, nil
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
func (c *CSP) Validate(logger log.Logger, r *http.Request, appKey string) error {
	ac := r.Header.Get(authContextHeader)
	if ac == "" {
		return middleware.ErrMissingHeader
	}

	authContext, err := c.getAuthContext(logger, ac, appKey)
	if err != nil {
		return middleware.ErrInvalidAuthContext
	}

	var authJSON cspAuthJSON

	err = json.Unmarshal(authContext, &authJSON)
	if err != nil {
		logger.Errorf("error while unmarshalling csp auth json, %v", err)

		return middleware.ErrInvalidAuthContext
	}

	bodyHash := GetBodyHash(r)

	// generate data and its signature for validation
	dataForSigValidation := authJSON.MachineName + authJSON.RequestDate + authJSON.IPAddress + appKey +
		c.sharedKey + authJSON.HTTPMethod + authJSON.UUID + authJSON.ClientID + bodyHash

	computedSignature := Base64Encode([]byte(HexEncode(Sha256Hash([]byte(dataForSigValidation)))))

	if computedSignature != authJSON.SignatureHash {
		return middleware.ErrInvalidAuthContext
	}

	return nil
}

func (c *CSP) getAuthContext(logger log.Logger, authContextHeader, appKey string) ([]byte, error) {
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

	keys := c.Get(appKey)

	// decrypt auth context using encryption key and initial vector
	decryptedAuthContext, err := decryptData(authContextToDecrypt, keys.encryptionKey, keys.iv)
	if err != nil {
		logger.Errorf("error occurred while decrypting auth context, %v", err)

		return nil, err
	}

	return decryptedAuthContext, nil
}
