package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/pbkdf2"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

const (
	encryptionBlockSizeBytes = 16
	cspEncryptionIterations  = 1000
	securityType             = "1"
	securityVersion          = "V1"
	securityTypeHeader       = "jst"
	securityVersionHeader    = "sv"
	appKeyHeader             = "ak"
	clientIDHeader           = "cd"
	authContextHeader        = "ac"
	lenRandomChars           = 6
	minLenAppKey             = 12
	nanoToMicroSecond        = 1000
)

var (
	// ErrEmptySharedKey is raised when shared key is empty
	ErrEmptySharedKey = errors.Error("shared key cannot be empty")
	// ErrEmptyAppKey is raised when app key is is not more than 12 bytes
	ErrEmptyAppKey = errors.Error("app key should be more than 12 bytes")
	// ErrEmptyAppID is raised when app id is empty
	ErrEmptyAppID = errors.Error("app id cannot be empty")

	errInvalidBlockSize = errors.Error("invalid block size")
	errInvalidPKCS7Data = errors.Error("invalid PKCS7 data (empty or not padded)")
)

type csp struct {
	options       *CSPOption
	encryptionKey []byte // encryptionKey to be used for aes encryption/decryption
	iv            []byte // initial vector(iv) to be used for aes encryption/decryption
}

// New validates the options and creates new instance of csp
func New(logger log.Logger, opts *CSPOption) (*csp, error) {
	if err := opts.validate(); err != nil {
		logger.Warnf("Invalid Options, %v", err)
		return nil, err
	}

	return &csp{
		options:       opts,
		encryptionKey: createKey([]byte(opts.AppKey), []byte(opts.AppKey[:12])),
		iv:            createKey([]byte(opts.SharedKey), []byte(opts.AppKey[:12])),
	}, nil
}

// getAuthContext generates the auth context
func (c *csp) getAuthContext(method string, body io.Reader) string {
	cspAuthJSONBytes := c.options.generateAuthJSON(method, getBodyHash(body))
	cipherText := encryptData(cspAuthJSONBytes, c.encryptionKey, c.iv)
	b64CtBeforeRand := base64Encode(cipherText)
	x := getRandomChars()
	cipherTextWithRand := append([]byte(b64CtBeforeRand), x...)

	return base64Encode(cipherTextWithRand)
}

type CSPOption struct {
	AppKey      string
	AppID       string
	SharedKey   string
	MachineName string
	IPAddress   string
}

func (o *CSPOption) validate() error {
	if o.SharedKey == "" {
		return ErrEmptySharedKey
	}

	if len(o.AppKey) < minLenAppKey {
		return ErrEmptyAppKey
	}

	if o.AppID == "" {
		return ErrEmptyAppID
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
	msgUniqueID := hexEncode(guid[:])

	// take hash of machineName+requestDate+ip+appKey+sharedKey+httpMethod+guid+clientId+bodyhash
	requestTime := genTimestamp()
	requestData := o.MachineName + requestTime + o.IPAddress + o.AppKey + o.SharedKey + method + msgUniqueID + o.AppID + bodyHash
	signatureHash := base64Encode([]byte(hexEncode(sha256Hash([]byte(requestData)))))
	authJSON := cspAuthJSON{
		IPAddress:     o.IPAddress,
		MachineName:   o.MachineName,
		ClientID:      o.AppID,
		HTTPMethod:    method,
		RequestDate:   requestTime,
		SignatureHash: signatureHash,
		UUID:          msgUniqueID,
	}

	authJSONBytes, _ := json.Marshal(authJSON)

	return authJSONBytes
}

func createKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, cspEncryptionIterations, encryptionBlockSizeBytes, sha1.New)
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
		return nil, errInvalidBlockSize
	}

	bLen := len(b)
	if bLen == 0 {
		return nil, errInvalidPKCS7Data
	}

	n := blocksize - (bLen % blocksize)
	pb := make([]byte, bLen+n)
	copy(pb, b)
	copy(pb[bLen:], bytes.Repeat([]byte{byte(n)}, n))

	return pb, nil
}

// Generate base64 encoded string
func base64Encode(body []byte) string {
	return base64.StdEncoding.EncodeToString(body)
}

// Generate hex encoded string
func hexEncode(body []byte) string {
	return strings.ToUpper(hex.EncodeToString(body))
}

func getBodyHash(body io.Reader) string {
	if body == nil {
		return ""
	}

	bodyBytes, _ := ioutil.ReadAll(body)

	return hexEncode(sha256Hash(bodyBytes))
}

// Generate sha256 hash
func sha256Hash(body []byte) []byte {
	hash := sha256.New()
	_, _ = hash.Write(body)

	return hash.Sum(nil)
}
