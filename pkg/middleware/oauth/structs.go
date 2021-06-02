package oauth

import (
	"crypto/rsa"
	"sync"
)

type OAuth struct {
	options Options
	cache   PublicKeyCache
}

type Options struct {
	// Set validity frequency in seconds
	ValidityFrequency int
	JWKPath           string
}

type header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
	URL       string `json:"jku"`
	KeyID     string `json:"kid"`
}

type JWT struct {
	payload   string
	header    header
	signature string
	token     string
}

type PublicKey struct {
	ID         string   `json:"kid"`
	Alg        string   `json:"alg"`
	Type       string   `json:"kty"`
	Use        string   `json:"use"`
	Operations []string `json:"key_ops"`

	// rsa fields
	Modulus         string `json:"n"`
	PublicExponent  string `json:"e"`
	PrivateExponent string `json:"d"`

	rsaPublicKey rsa.PublicKey
}

type PublicKeys struct {
	Keys []PublicKey `json:"keys"`
}

type PublicKeyCache struct {
	publicKeys PublicKeys
	mu         sync.RWMutex
}
