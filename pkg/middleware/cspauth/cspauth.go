package cspauth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/log"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/pbkdf2"
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

var (
	// ErrEmptySharedKey is raised when shared key is empty
	ErrEmptySharedKey = errors.Error("shared key cannot be empty")
	// ErrEmptyAppKey is raised when app key is is not more than 12 bytes
	ErrEmptyAppKey = errors.Error("app key should be more than 12 bytes for successful key generation")
	// ErrEmptyAppID is raised when app id is empty
	ErrEmptyAppID = errors.Error("app id cannot be empty")

	errInvalidBlockSize    = errors.Error("invalid blocksize")
	errInvalidPKCS7Data    = errors.Error("invalid PKCS7 data (empty or not padded)")
	errInvalidPKCS7Padding = errors.Error("invalid padding on input")
)

// CSP generates and validates csp auth headers
type CSP struct {
	options       *Options
	encryptionKey []byte //encryption key to be used for aes encryption/decryption
	iv            []byte //iv to be used for aes encryption/decryption
}

// New creates new instance of CSP
func New(logger log.Logger, opts *Options) *CSP {
	if err := opts.validate(); nil != err {
		logger.Warnf("Invalid Options, %v", err)
		return nil
	}
	return &CSP{
		options:       opts,
		encryptionKey: createKey([]byte(opts.AppKey), []byte(opts.AppKey[:12])),
		iv:            createKey([]byte(opts.SharedKey), []byte(opts.AppKey[:12])),
	}
}

// CSPAuth middleware
func CSPAuth(logger log.Logger, opts Options) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if opts.SharedKey == "" {
				inner.ServeHTTP(w, req)

				return
			}

			setOptions(&opts, req)
			csp := New(logger, &opts)

			if csp == nil {
				logger.Warnf("request doesn't contains valid headers for CSP Auth")
				w.WriteHeader(http.StatusBadRequest)

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
	decryptedAuthContext, err := decryptData([]byte(authContextToDecrypt), c.encryptionKey, c.iv)
	if err != nil {
		logger.Errorf("error occurred while decrypting auth context, %v", err)
		return false
	}
	//unmarshal the decrypted data into cspAuthJson object
	var authJSON cspAuthJSON

	err = json.Unmarshal([]byte(decryptedAuthContext), &authJSON)
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

func createKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, cspEncryptionIterations, encryptionBlockSizeBytes, sha1.New)
}

func getBody(r *http.Request) []byte {
	if r.Body == nil {
		return []byte{}
	}
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

//Generate sha 256 hash
func sha256Hash(body []byte) []byte {
	hash := sha256.New()
	hash.Write(body)
	return hash.Sum(nil)
}

//Generate base 64 encoded string
func base64Encode(body []byte) string {
	return base64.StdEncoding.EncodeToString(body)
}

//Generate hex encoded string
func hexEncode(body []byte) string {
	return strings.ToUpper(hex.EncodeToString(body))
}

//Decode base 64 string
func base64Decode(body string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(body)
}

//generate timestamp in format "YYYY-MM-DD hh:mm:ss.uuuuuu" (microseconds)
func genTimestamp() string {
	t := time.Now().UTC()
	ts := t.Format("2006-01-02 15:04:05.") + strconv.Itoa(t.Nanosecond()/1000)
	return ts
}

func encryptData(plaintext []byte, key []byte, iv []byte) []byte {
	aesCipher, _ := aes.NewCipher(key)
	aesEncrypter := cipher.NewCBCEncrypter(aesCipher, iv)
	plaintext, _ = pkcs7Pad(plaintext, aesCipher.BlockSize())

	ciphertext := make([]byte, len(plaintext))
	aesEncrypter.CryptBlocks(ciphertext, plaintext)
	return ciphertext
}

func getRandomChars() []byte {
	uu := uuid.NewV4()
	ux := uu.String()
	return []byte(ux[:lenRandomChars])
}

func decryptData(ciphertext []byte, key []byte, iv []byte) (plaintext []byte, err error) {
	aesCipher, _ := aes.NewCipher(key)
	aesDecrypter := cipher.NewCBCDecrypter(aesCipher, iv)
	plaintext = make([]byte, len(ciphertext))
	aesDecrypter.CryptBlocks(plaintext, ciphertext)
	plaintext, err = pkcs7Unpad(plaintext, aesCipher.BlockSize())

	return
}

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, errInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, errInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, errInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, errInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, errInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}
