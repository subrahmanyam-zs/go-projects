package cspauth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/pbkdf2"
)

var (
	errInvalidBlockSize    = errors.Error("invalid blocksize")
	errInvalidPKCS7Data    = errors.Error("invalid PKCS7 data (empty or not padded)")
	errInvalidPKCS7Padding = errors.Error("invalid padding on input")
)

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
