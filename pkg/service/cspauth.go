package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
	"developer.zopsmart.com/go/gofr/pkg/middleware/cspauth"
)

const (
	lenRandomChars    = 6
	minLenAppKey      = 12
	nanoToMicroSecond = 1000
	securityType      = "1"
	securityVersion   = "V1"
)

var (
	// ErrEmptySharedKey is raised when shared key is empty
	ErrEmptySharedKey = errors.Error("shared key cannot be empty")
)

type csp struct {
	options       *CSPOption
	encryptionKey []byte // encryptionKey to be used for aes encryption/decryption
	iv            []byte // initial vector(iv) to be used for aes encryption/decryption
}

// NewCSP validates the options and creates new instance of csp
func NewCSP(logger log.Logger, opts *CSPOption) (*csp, error) {
	if err := opts.validate(); err != nil {
		logger.Warnf("Invalid Options, %v", err)
		return nil, err
	}

	return &csp{
		options:       opts,
		encryptionKey: cspauth.CreateKey([]byte(opts.AppKey), []byte(opts.AppKey[:12]), 32),
		iv:            cspauth.CreateKey([]byte(opts.SharedKey), []byte(opts.AppKey[:12]), 16),
	}, nil
}

// getAuthContext generates the auth context
func (c *csp) getAuthContext(req *http.Request) string {
	cspAuthJSONBytes := c.options.generateAuthJSON(req.Method, cspauth.GetBodyHash(req))
	cipherText := encryptData(cspAuthJSONBytes, c.encryptionKey, c.iv)
	b64CtBeforeRand := cspauth.Base64Encode(cipherText)
	x := getRandomChars()
	cipherTextWithRand := append([]byte(b64CtBeforeRand), x...)

	return cspauth.Base64Encode(cipherTextWithRand)
}

type CSPOption struct {
	AppKey      string
	ClientID    string
	SharedKey   string
	MachineName string
	IPAddress   string
}

func (o *CSPOption) validate() error {
	if o.SharedKey == "" {
		return ErrEmptySharedKey
	}

	if len(o.AppKey) < minLenAppKey {
		return middleware.ErrInvalidAppKey
	}

	return nil
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

func (o *CSPOption) generateAuthJSON(method, bodyHash string) []byte {
	guid := uuid.NewV4()
	msgUniqueID := cspauth.HexEncode(guid[:])

	// take hash of machineName+requestDate+ip+appKey+sharedKey+httpMethod+guid+clientId+bodyhash
	requestTime := genTimestamp()
	requestData := o.MachineName + requestTime + o.IPAddress + o.AppKey + o.SharedKey + method + msgUniqueID + o.ClientID + bodyHash
	signatureHash := cspauth.Base64Encode([]byte(cspauth.HexEncode(cspauth.Sha256Hash([]byte(requestData)))))
	authJSON := cspAuthJSON{
		IPAddress:     o.IPAddress,
		MachineName:   o.MachineName,
		ClientID:      o.ClientID,
		HTTPMethod:    method,
		RequestDate:   requestTime,
		SignatureHash: signatureHash,
		UUID:          msgUniqueID,
	}

	authJSONBytes, _ := json.Marshal(authJSON)

	return authJSONBytes
}

func encryptData(plaintext, key, iv []byte) []byte {
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

// generate timestamp in format "YYYY-MM-DD hh:mm:ss.uuuuuu" (microseconds)
func genTimestamp() string {
	t := time.Now().UTC()
	ts := t.Format("2006-01-02 15:04:05.") + strconv.Itoa(t.Nanosecond()/nanoToMicroSecond)

	return ts
}

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, cspauth.ErrInvalidBlockSize
	}

	bLen := len(b)
	if bLen == 0 {
		return nil, cspauth.ErrInvalidPKCS7Data
	}

	n := blocksize - (bLen % blocksize)
	pb := make([]byte, bLen+n)
	copy(pb, b)
	copy(pb[bLen:], bytes.Repeat([]byte{byte(n)}, n))

	return pb, nil
}
