package cspauth

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

// Options used to initialize CSP
type Options struct {
	MachineName string
	IPAddress   string
	AppKey      string
	SharedKey   string
	AppID       string
}

func setOptions(opts *Options, req *http.Request) {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		// get the ip address from the RemoteAddr
		ip = strings.Split(req.RemoteAddr, ":")[0]
	}

	opts.MachineName = req.Header.Get("User-Agent")
	opts.IPAddress = ip
	opts.AppKey = req.Header.Get(appKeyHeader)
	opts.AppID = req.Header.Get(clientIDHeader)
}

func (o *Options) validate() error {
	if 0 == len(o.SharedKey) {
		return ErrEmptySharedKey
	}
	if 12 > len(o.AppKey) {
		return ErrEmptyAppKey
	}
	if 0 == len(o.AppID) {
		return ErrEmptyAppID
	}
	return nil
}

func (o *Options) generateAuthJSON(method string, body []byte) []byte {
	guid := uuid.NewV4()
	msgUniqueID := hexEncode(guid[:])
	var bodyHash string
	if 0 < len(body) {
		bodyHash = hexEncode(sha256Hash(body))
	}
	//take hash of machineName+requestDate+ip+appKey+sharedKey+httpMethod+guid+clientId+bodyhash
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
