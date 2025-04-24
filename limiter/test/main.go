package main

import (
	"context"
	"limiter/bucket"
	"net/http"
)

func NewLimiter(ctx context.Context, rps, burst int, next http.HandlerFunc) http.HandlerFunc {
	b := bucket.New(ctx, rps, burst)
	return func(w http.ResponseWriter, r *http.Request) {
		if err := b.Wait(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		next(w, r)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.RemoteAddr))
}

func main() {
	http.HandleFunc("/hello", NewLimiter(context.Background(), 100, 10, hello))
	http.ListenAndServe(":9999", nil)
}
