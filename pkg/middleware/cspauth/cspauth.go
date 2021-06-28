package cspauth

import (
	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
	"encoding/json"
	"net/http"
)

const (
	securityType             = "1"
	securityVersion          = "V1"
	encryptionBlockSizeBytes = 16
	cspEncryptionIterations  = 1000
	lenRandomChars           = 6
	appKeyHeader             = "ak"
	clientIDHeader           = "cd"
	authContextHeader        = "ac"
	securityTypeHeader       = "st"
	securityVersionHeader    = "sv"
)

// CSP generates and validates csp auth headers
type CSP struct {
	options       *Options
	encryptionKey []byte //encryption key to be used for aes encryption/decryption
	iv            []byte //iv to be used for aes encryption/decryption
}

// New creates new instance of CSP
func New(logger log.Logger, opts *Options) (*CSP, error) {
	if err := opts.validate(); err != nil {
		logger.Warnf("Invalid Options, %v", err)
		return nil, err
	}

	return &CSP{
		options:       opts,
		encryptionKey: createKey([]byte(opts.AppKey), []byte(opts.AppKey[:12])),
		iv:            createKey([]byte(opts.SharedKey), []byte(opts.AppKey[:12])),
	}, nil
}

// CSPAuth middleware
func CSPAuth(logger log.Logger, opts Options) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			setOptions(&opts, req)

			csp, err := New(logger, &opts)
			if err != nil {
				e := middleware.FetchErrResponseWithCode(http.StatusBadRequest, "Invalid CSP Auth Options", err.Error())
				middleware.ErrorResponse(w, req, logger, *e)

				return
			}

			if ok := csp.Verify(logger, req); !ok {
				csp.Set(req)
			}

			inner.ServeHTTP(w, req)
		})
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

// Set the csp auth headers in http request
func (c *CSP) Set(r *http.Request) {
	cspAuthJSONBytes := c.options.generateAuthJSON(r.Method, getBody(r))
	cipherText := encryptData(cspAuthJSONBytes, c.encryptionKey, c.iv)
	b64CtBeforeRand := base64Encode(cipherText)
	x := getRandomChars()
	cipherTextWithRand := append([]byte(b64CtBeforeRand), x...)
	b64ct := base64Encode(cipherTextWithRand)

	r.Header.Set(appKeyHeader, c.options.AppKey)
	r.Header.Set(clientIDHeader, c.options.AppID)
	r.Header.Set(securityVersionHeader, securityVersion)
	r.Header.Set(securityTypeHeader, securityType)
	r.Header.Set(authContextHeader, b64ct)
}

// Verify the csp auth headers in given request
func (c *CSP) Verify(logger log.Logger, r *http.Request) bool {
	//base64 decoding the auth context.
	b64DecodeRandom, err := base64Decode(r.Header.Get(authContextHeader))
	if err != nil {
		logger.Errorf("error while base64 decoding auth context, %v", err)
		return false
	}

	if len(b64DecodeRandom) <= lenRandomChars {
		logger.Errorf("Invalid auth context, %v", err)
		return false
	}
	//remove random string from auth context
	authContextToDecode := b64DecodeRandom[:len(b64DecodeRandom)-lenRandomChars]
	//base64 decode auth context
	authContextToDecrypt, err := base64Decode(string(authContextToDecode))
	if err != nil {
		logger.Errorf("error while base64 decoding auth context without random chars, %v", err)
		return false
	}

	//now decrypt auth context.
	decryptedAuthContext, err := decryptData(authContextToDecrypt, c.encryptionKey, c.iv)
	if err != nil {
		logger.Errorf("error occurred while decrypting auth context, %v", err)
		return false
	}
	//unmarshal the decrypted data into cspAuthJson object
	var authJSON cspAuthJSON

	err = json.Unmarshal(decryptedAuthContext, &authJSON)
	if err != nil {
		logger.Errorf("error while unmarshalling csp auth json, %v", err)
		return false
	}

	httpBody := getBody(r)

	//generate requestData string and take base64encoded(hash(requestdata)) and compare with signaturehash in json
	var bodyHash string
	if len(httpBody) > 0 {
		bodyHash = hexEncode(sha256Hash(httpBody))
	}
	dataForSigValidation := authJSON.MachineName + authJSON.RequestDate + authJSON.IPAddress + c.options.AppKey + c.options.SharedKey + authJSON.HTTPMethod + authJSON.UUID + authJSON.ClientID + bodyHash

	//compute signature with data for signature validation
	computedSignature := base64Encode([]byte(hexEncode(sha256Hash([]byte(dataForSigValidation)))))
	//match signature check
	return computedSignature == authJSON.SignatureHash
}
