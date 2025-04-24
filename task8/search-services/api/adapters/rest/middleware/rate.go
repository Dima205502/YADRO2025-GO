package middleware

import (
	"net/http"
	"time"

	"yadro.com/course/api/core"
)

func Rate(next http.HandlerFunc, algo core.LimiterAlgo, limit int, interval time.Duration) http.HandlerFunc {
	rl := core.GetRateLimiter(algo, limit, interval)
	return func(w http.ResponseWriter, r *http.Request) {
		rl.Wait()
		next(w, r)
	}
}
