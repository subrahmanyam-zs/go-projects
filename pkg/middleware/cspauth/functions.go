package cspauth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	//nolint:gosec //existing functionality cannot be deprecated
	"crypto/sha1"

	"golang.org/x/crypto/pbkdf2"
)

// CreateKey generates key of length keyLen from the password, salt and iteration count
func CreateKey(password, salt []byte, keyLen int) []byte {
	return pbkdf2.Key(password, salt, cspEncryptionIterations, keyLen, sha1.New)
}

// GetBodyHash returns the encoded hash of the body
func GetBodyHash(r *http.Request) string {
	if r.Body == nil {
		return ""
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	lenBody := len(bodyBytes)
	if lenBody == 0 {
		return ""
	}

	return HexEncode(Sha256Hash(bodyBytes))
}

// Sha256Hash generate sha256 hash
func Sha256Hash(body []byte) []byte {
	hash := sha256.Sum256(body)

	return hash[:]
}

// Base64Encode generate base64 encoded string
func Base64Encode(body []byte) string {
	return base64.StdEncoding.EncodeToString(body)
}

// HexEncode generate hex encoded string
func HexEncode(body []byte) string {
	return strings.ToUpper(hex.EncodeToString(body))
}

// Decode base 64 string
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
		return nil, ErrInvalidBlockSize
	}

	bLen := len(b)

	if bLen == 0 {
		return nil, ErrInvalidPKCS7Data
	}

	if bLen%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}

	c := b[bLen-1]
	n := int(c)

	if n == 0 || n > bLen {
		return nil, ErrInvalidPKCS7Padding
	}

	for i := 0; i < n; i++ {
		if b[bLen-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}

	return b[:bLen-n], nil
}
