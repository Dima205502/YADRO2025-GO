package middleware

import (
	"net/http"
)

func Concurrency(next http.HandlerFunc, limit int) http.HandlerFunc {
	semaphore := make(chan struct{}, limit)
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case semaphore <- struct{}{}:
			defer func() { <-semaphore }()
			next(w, r)
		default:
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		}
	}
}
