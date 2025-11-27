// utils/rate_limiter.go
package utils

import (
	"net/http"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Limit(1000), 5000)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
