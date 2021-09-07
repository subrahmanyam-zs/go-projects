package oauth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bmizerany/assert"

	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
)

func TestLDAPOAuth(t *testing.T) {
	basicHeader := "Basic Wk9QT1JERVJTVkM6ZHNmaDcyaHM2QWhzbmE="
	// nolint
	oAuthHeader := "bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjIwMTEtMDQtMjk9PSJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.B5C9tz71T-PjyoMH-gv198iNFguDZ5SpVcwrgdLxU83A92o1tsJWh8_7Zm6ulMUupNEAzGD69DB077j01nXz6ut5XtnXWE50HNTxlS_19zndpPxqFcKnWyoArip5A1MCgQjKQ3exwZc7aFQwgBXvJMNk-5N4od_bUMGvOb0q3ApbfzbwIt94daToPjhfLy4xf8UoNhh_Lq14CNHCZXNgGeter5TvnHnDBN4oDfw6nziKdJnslNkUJ2hHsqp8VObUK57C8aS51x2UiOwTJ1NqDv0PFVgRbC7ncFZG6M87x9BGTwB0XvraXYU7Zimewp4plzdIMnjIXXp8kuviYl7feA"
	testCases := []struct {
		endPoint     string
		ldapGroup    string
		authHeader   string
		expectedCode int
	}{
		{"/hello", "order-service", basicHeader, 200},
		{"/hello", "order-service", oAuthHeader, 200},
		{"/hello", "", oAuthHeader, 200},
		{"/hello", "", basicHeader, 200},
		{"/hello", "", "", 200},
		{"/.well-known/heartbeat", "", "", 200},
	}

	b := new(bytes.Buffer)
	mockLogger := log.NewMockLogger(b)
	oAuthOptions := Options{
		ValidityFrequency: 0,
		JWKPath:           getTestServerURL(),
	}

	for _, testCase := range testCases {
		ldapOptions := middleware.LDAPOptions{
			RegexToMethodGroup: map[string][]middleware.MethodGroup{
				"hello": {{
					Method: "GET,POST",
					Group:  testCase.ldapGroup,
				}},
			},
			Addr:                       "ldapstage.zopsmart.com:636",
			CacheInvalidationFrequency: 10,
			InsecureSkipVerify:         true,
		}

		handler := LDAPOAuth(mockLogger, &ldapOptions, oAuthOptions)(&MockHandlerForLDAPOAuth{})
		w := httptest.NewRecorder()
		request := httptest.NewRequest("GET", testCase.endPoint, nil)
		request.Header.Set("Authorization", testCase.authHeader)
		handler.ServeHTTP(w, request)
		assert.Equal(t, testCase.expectedCode, w.Code)
	}
}

func TestVerifyOptions(t *testing.T) {
	testCases := []struct {
		ldapAddress  string
		JWKAddress   string
		ldapGroups   bool
		errorMessage string
	}{
		{"ldapstage.zopsmart.com:636", getTestServerURL(), true, ""},
		{"ldapstage.zopsmart.com:636", getTestServerURL(), false, "no mappings defined for LDAP groups"},
		{"ldapstage.zopsmart.com:636", "", true, "empty oAuth options/ missing JWK End point."},
		{"ldapstage.zopsmart.com:636", "", false, "empty oAuth options/ missing JWK End point."},
		{"", getTestServerURL(), true, "empty LDAP options/ missing LDAP Address"},
		{"", getTestServerURL(), false, "empty LDAP options/ missing LDAP Address"},
		{"", "", true, "empty LDAP options/ missing LDAP Address"},
		{"", "", false, "empty LDAP options/ missing LDAP Address"},
	}

	oAuthOptions := Options{
		ValidityFrequency: 0,
	}
	regexMapping := map[string][]middleware.MethodGroup{
		"hello": {{
			Method: "GET,POST",
			Group:  "random-group",
		}},
	}
	ldapOptions := middleware.LDAPOptions{
		CacheInvalidationFrequency: 10,
		InsecureSkipVerify:         true,
	}

	for i, testCase := range testCases {
		oAuthOptions.JWKPath = testCase.JWKAddress
		ldapOptions.Addr = testCase.ldapAddress

		if testCase.ldapGroups {
			ldapOptions.RegexToMethodGroup = regexMapping
		} else {
			ldapOptions.RegexToMethodGroup = nil
		}

		expectedResult := testCase.errorMessage == ""

		b := new(bytes.Buffer)
		mockLogger := log.NewMockLogger(b)

		result := verifyOptions(mockLogger, &ldapOptions, &oAuthOptions)
		assert.Equal(t, expectedResult, result, i)
		assert.T(t, strings.Contains(b.String(), testCase.errorMessage), testCase.errorMessage, b.String(), i)
	}
}

type MockHandlerForLDAPOAuth struct{}

func (r *MockHandlerForLDAPOAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}
