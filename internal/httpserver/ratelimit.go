package httpserver

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			key := "rl:" + ip

			count, err := rdb.Incr(r.Context(), key).Result()
			if err != nil {
				http.Error(w, "rate limit error", http.StatusInternalServerError)
				return
			}

			// set TTL hanya saat pertama kali key dibuat
			if count == 1 {
				_ = rdb.Expire(r.Context(), key, window).Err()
			}

			if count > int64(limit) {
				http.Error(w, "too many requests", http.StatusTooManyRequests) // 429
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	// Kalau nanti behind proxy: X-Forwarded-For bisa dipakai.
	// Untuk local dev, ini aman juga.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// ambil IP pertama
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	// fallback: RemoteAddr (buang port)
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
