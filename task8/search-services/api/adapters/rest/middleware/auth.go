package middleware

import (
	"net/http"
	"strings"
)

type TokenVerifier interface {
	Verify(token string) error
}

func Auth(next http.HandlerFunc, verifier TokenVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		authType, token, found := strings.Cut(authHeader, " ")
		if !found || authType != "Token" {
			http.Error(w, "invalid authorization format", http.StatusUnauthorized)
			return
		}

		if err := verifier.Verify(token); err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
