package cspauth

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

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

func (o *Options) generateAuthJSON(method string, body []byte) []byte {
	guid := uuid.NewV4()
	msgUniqueID := hexEncode(guid[:])

	var bodyHash string

	if len(body) > 0 {
		bodyHash = hexEncode(sha256Hash(body))
	}
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
	bytes, _ := json.Marshal(authJSON)

	return bytes
}
