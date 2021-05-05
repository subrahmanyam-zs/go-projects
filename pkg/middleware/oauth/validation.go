package oauth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zopsmart/gofr/pkg/middleware"

	"github.com/zopsmart/gofr/pkg/log"
)

func getJWT(logger log.Logger, r *http.Request) (JWT, error) {
	token := r.Header.Get("Authorization")

	jwtVal := strings.Fields(token)
	if token == "" || len(jwtVal) != 2 || !strings.EqualFold(jwtVal[0], "bearer") {
		return JWT{}, middleware.ErrInvalidRequest
	}

	// Checking if incoming token string conforms to the predefined jwt structure
	jwtParts := strings.Split(jwtVal[1], ".")

	const JWTPartsLen = 3
	if len(jwtParts) != JWTPartsLen {
		logger.Error("jwt token is not of the format hhh.ppp.sss")
		return JWT{}, middleware.ErrInvalidToken
	}

	var h header

	decodedHeader, err := base64.RawStdEncoding.DecodeString(jwtParts[0])
	if err != nil {
		return JWT{}, middleware.ErrInvalidToken
	}

	err = json.Unmarshal(decodedHeader, &h)
	if err != nil {
		return JWT{}, middleware.ErrInvalidToken
	}

	return JWT{payload: jwtParts[1], header: h, signature: jwtParts[2], token: jwtVal[1]}, nil
}
