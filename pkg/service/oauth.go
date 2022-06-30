package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/errors"

	"go.opencensus.io/plugin/ochttp"
)

type OAuthOption struct {
	ClientID       string
	ClientSecret   string
	KeyProviderURL string
	Scope          string
	MaxSleep       int
	Audience       string
}

// nolint:gocognit // need to add new condition to check clientID and clientSecret
func (h *httpService) setClientOauthHeader(option *OAuthOption) {
	if option.ClientID == "" || option.ClientSecret == "" {
		return
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(option.ClientID+":"+option.ClientSecret))

	if option.MaxSleep <= 0 {
		option.MaxSleep = pkg.OAuthMaxSleep
	}

	retryList := oauthFibonacciRetry(option.MaxSleep)

	var (
		token   string
		expTime int
	)

	lenRetry := len(retryList)
	for i := 1; i < lenRetry; i++ {
		newToken, exp, err := getNewAccessToken(basicAuth, option)
		if err != nil {
			h.logger.Errorf("%v %v", h.url, err)
		}

		if exp <= 0 {
			time.Sleep(time.Duration(retryList[i]-retryList[i-1]) * time.Second)
			continue
		} else {
			token = newToken
			expTime = exp
			break
		}
	}

	h.mu.Lock()
	h.auth = token
	h.mu.Unlock()

	// refresh token 5 seconds before the token expires
	if expTime > pkg.OAuthExpiryBeforeTime {
		expTime -= pkg.OAuthExpiryBeforeTime
	}

	go func() {
		time.Sleep(time.Duration(expTime) * time.Second)
		h.setClientOauthHeader(option)
	}()
}

// nolint:gocognit // cognitive complexity of func is high (> 10)
func getNewAccessToken(basicAuth string, option *OAuthOption) (bearerToken string, exp int, err error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	if option.Scope != "" {
		data.Set("scope", option.Scope)
	}

	if option.Audience != "" {
		data.Set("audience", option.Audience)
	}

	reqHeaders := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	tokenService := httpService{
		Client:        &http.Client{Transport: &ochttp.Transport{}, Timeout: RetryFrequency * time.Second}, // default timeout is 5 seconds
		auth:          basicAuth,
		url:           option.KeyProviderURL,
		customHeaders: reqHeaders,
		isHealthy:     true,
	}

	var response map[string]interface{}

	resp, err := tokenService.Post(context.Background(), "", nil, []byte(data.Encode()))
	if err == nil && resp != nil {
		err = json.Unmarshal(resp.Body, &response)
	}

	if err != nil {
		return "", 0, err
	}

	var responseError error

	if !successStatusRange(resp.StatusCode) {
		responseError = &errors.Response{Reason: string(resp.Body)}
	} else {
		responseError = nil
	}

	if v, ok := response["access_token"].(string); ok {
		bearerToken = "Bearer " + v
	}

	if e, ok := response["expires_in"].(float64); ok {
		exp = int(e)
	}

	return bearerToken, exp, responseError
}

func oauthFibonacciRetry(max int) []int {
	var (
		firstElement  = 8
		secondElement = 13
		retryList     = []int{firstElement, secondElement}
	)

	for {
		if firstElement+secondElement > max {
			break
		}

		retryList = append(retryList, firstElement+secondElement)
		secondElement = firstElement + secondElement
		firstElement = secondElement - firstElement
	}

	retryList = append(retryList, max)

	return retryList
}

func successStatusRange(status int) bool {
	if status >= http.StatusOK && status <= pkg.HTTPMaxSuccessfulResponseCode {
		return true
	}

	return false
}
