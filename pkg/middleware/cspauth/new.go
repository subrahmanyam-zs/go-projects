package cspauth

import "developer.zopsmart.com/go/gofr/pkg/log"

// CSP generates and validates csp auth headers
type CSP struct {
	options       *Options
	encryptionKey []byte // encryptionKey to be used for aes encryption/decryption
	iv            []byte // initial vector(iv) to be used for aes encryption/decryption
}

// New validates the options and creates new instance of CSP
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

// Options used to initialize CSP
type Options struct {
	MachineName string
	IPAddress   string
	AppKey      string
	SharedKey   string
	AppID       string
}

func (o *Options) validate() error {
	if o.SharedKey == "" {
		return ErrEmptySharedKey
	}

	if len(o.AppKey) < minAppKeyLen {
		return ErrEmptyAppKey
	}

	if o.AppID == "" {
		return ErrEmptyAppID
	}

	return nil
}
