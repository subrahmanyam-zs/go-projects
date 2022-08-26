package middleware

import (
	"log"
	"net/http"

	slice "golang.org/x/exp/slices"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-api-key")
		apiKeys := []string{"jason470", "jason573", "jason"}
		if !slice.Contains(apiKeys, token) {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unauthorized Key"))
			if err != nil {
				log.Println(err)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
