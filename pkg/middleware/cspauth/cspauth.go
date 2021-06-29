package cspauth

import (
	"encoding/json"
	"net/http"
	"strings"

	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
)

// CSPAuth middleware provides authentication using CSP
func CSPAuth(logger log.Logger, opts *Options) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			setOptions(opts, req)

			csp, err := New(logger, opts)
			if err != nil {
				e := middleware.FetchErrResponseWithCode(http.StatusBadRequest, "Invalid CSP Auth Options", err.Error())
				middleware.ErrorResponse(w, req, logger, *e)

				return
			}

			if ok := csp.Verify(logger, req); !ok {
				e := middleware.FetchErrResponseWithCode(http.StatusForbidden, "Invalid CSP Auth Context", "")
				middleware.ErrorResponse(w, req, logger, *e)

				return
			}

			inner.ServeHTTP(w, req)
		})
	}
}

func setOptions(opts *Options, req *http.Request) {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		// ip address from the RemoteAddr
		ip = strings.Split(req.RemoteAddr, ":")[0]
	}

	opts.MachineName = req.Header.Get("User-Agent")
	opts.IPAddress = ip
	opts.AppKey = req.Header.Get(appKeyHeader)
	opts.AppID = req.Header.Get(clientIDHeader)
}

// Verify the csp auth headers in given request
func (c *CSP) Verify(logger log.Logger, r *http.Request) bool {
	b64DecodeRandom, err := base64Decode(r.Header.Get(authContextHeader))
	if err != nil {
		logger.Errorf("error while base64 decoding auth context, %v", err)
		return false
	}

	if len(b64DecodeRandom) <= lenRandomChars {
		logger.Errorf("Invalid auth context, %v", err)
		return false
	}
	// remove random string from auth context
	authContextToDecode := b64DecodeRandom[:len(b64DecodeRandom)-lenRandomChars]

	authContextToDecrypt, err := base64Decode(string(authContextToDecode))
	if err != nil {
		logger.Errorf("error while base64 decoding auth context without random chars, %v", err)
		return false
	}

	// decrypt auth context using encryption key and initial vector
	decryptedAuthContext, err := decryptData(authContextToDecrypt, c.encryptionKey, c.iv)
	if err != nil {
		logger.Errorf("error occurred while decrypting auth context, %v", err)
		return false
	}

	var authJSON cspAuthJSON

	err = json.Unmarshal(decryptedAuthContext, &authJSON)
	if err != nil {
		logger.Errorf("error while unmarshalling csp auth json, %v", err)
		return false
	}

	httpBody := getBody(r)

	var bodyHash string
	if len(httpBody) > 0 {
		bodyHash = hexEncode(sha256Hash(httpBody))
	}

	dataForSigValidation := authJSON.MachineName + authJSON.RequestDate + authJSON.IPAddress + c.options.AppKey +
		c.options.SharedKey + authJSON.HTTPMethod + authJSON.UUID + authJSON.ClientID + bodyHash

	// compute signature with data for signature validation
	computedSignature := base64Encode([]byte(hexEncode(sha256Hash([]byte(dataForSigValidation)))))

	return computedSignature == authJSON.SignatureHash
}
