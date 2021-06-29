package cspauth

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
