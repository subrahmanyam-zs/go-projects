package service

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/zopsmart/gofr/pkg/log"
	"go.opencensus.io/plugin/ochttp"
)

type soapService struct {
	httpService
}

//nolint:golint // this type cannot be exported since we don't want the user to have access to the members
func NewSOAPClient(resourceURL string, logger log.Logger, user, pass string) *soapService {
	auth := ""
	if user != "" {
		auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
	}

	transport := &ochttp.Transport{}

	return &soapService{
		httpService{
			url:    resourceURL,
			logger: logger,
			Client: &http.Client{Transport: transport, Timeout: RetryFrequency * time.Second}, // default timeout is 5 seconds
			sp: surgeProtector{
				isEnabled:             false,
				customHeartbeatURL:    "/.well-known/heartbeat",
				retryFrequencySeconds: RetryFrequency,
			},
			auth:      auth,
			isHealthy: true,
			healthCh:  make(chan bool),
		}}
}

// Call is a soap call for the given SOAP Action and body. The only allowed method in SOAP is POST
func (s *soapService) Call(ctx context.Context, action string, body []byte) (*Response, error) {
	return s.call(ctx, "POST", "", nil, body, map[string]string{"SOAPAction": action, "Content-Type": "text/xml"})
}

// CallWithHeaders is a soap call for the given SOAP Action and body. The only allowed method in SOAP is POST
func (s *soapService) CallWithHeaders(ctx context.Context, action string, body []byte, headers map[string]string) (*Response, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	headers["SOAPAction"] = action
	headers["Content-Type"] = "text/xml"

	return s.call(ctx, "POST", "", nil, body, headers)
}

func (s *soapService) Bind(resp []byte, i interface{}) error {
	s.httpService.contentType = XML

	return s.httpService.Bind(resp, i)
}

func (s *soapService) BindStrict(resp []byte, i interface{}) error {
	s.httpService.contentType = XML

	return s.httpService.Bind(resp, i)
}
