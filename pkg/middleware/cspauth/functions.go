package cspauth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

func createKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, cspEncryptionIterations, encryptionBlockSizeBytes, sha1.New)
}

func getBody(r *http.Request) []byte {
	if r.Body == nil {
		return []byte{}
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}

// Generate sha256 hash
func sha256Hash(body []byte) []byte {
	hash := sha256.New()
	_, _ = hash.Write(body)

	return hash.Sum(nil)
}

// Generate base64 encoded string
func base64Encode(body []byte) string {
	return base64.StdEncoding.EncodeToString(body)
}

// Generate hex encoded string
func hexEncode(body []byte) string {
	return strings.ToUpper(hex.EncodeToString(body))
}

// Decode base64 string
func base64Decode(body string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(body)
}

func decryptData(ciphertext, key, iv []byte) (plaintext []byte, err error) {
	aesCipher, _ := aes.NewCipher(key)
	aesDecrypter := cipher.NewCBCDecrypter(aesCipher, iv)
	plaintext = make([]byte, len(ciphertext))
	aesDecrypter.CryptBlocks(plaintext, ciphertext)
	plaintext, err = pkcs7Unpad(plaintext, aesCipher.BlockSize())

	return
}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errInvalidBlockSize
	}

	if len(b) == 0 {
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
