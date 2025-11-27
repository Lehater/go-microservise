// utils/rate_limiter.go
package utils

import (
	"net/http"

	"golang.org/x/time/rate"
)

// 1000 req/s и burst 5000 — стартовые настройки под wrk-тест.
// При необходимости burst можно увеличить и описать это в отчёте.
var limiter = rate.NewLimiter(rate.Limit(1000), 5000)

//	func RateLimitMiddleware(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			if !limiter.Allow() {
//				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
//				return
//			}
//			next.ServeHTTP(w, r)
//		})
//	}
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Чтение не ограничиваем: важно сохранить максимально возможный RPS и минимальную latency.
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		// Изменяющие операции (POST/PUT/DELETE) не должны превышать 1000 req/s.
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
