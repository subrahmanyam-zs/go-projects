package cspauth

import "developer.zopsmart.com/go/gofr/pkg/log"

// CSP generates and validates csp auth headers
type CSP struct {
	options       *Options
	encryptionKey []byte // encryptionKey to be used for aes encryption/decryption
	iv            []byte // initial vector(iv) to be used for aes encryption/decryption
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

type cspAuthJSON struct {
	IPAddress     string `json:"IPAddress"`
	MachineName   string `json:"MachineName"`
	RequestDate   string `json:"RequestDate"`
	HTTPMethod    string `json:"HttpMethod"`
	UUID          string `json:"MsgUniqueId"`
	ClientID      string `json:"ClientId"`
	SignatureHash string `json:"SignatureHash"`
}
