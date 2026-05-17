package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const ResetCodeKey contextKey = "reset_code"

func RateLimit(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// читаємо body
			var req struct {
				Email string `json:"email"`
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}

			// парсимо JSON
			_ = json.Unmarshal(bodyBytes, &req)

			// повертаємо body назад
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			ip := getIP(r)

			// rate limit по IP
			if !allow(r.Context(), "reset:ip:"+ip, 5, time.Minute, rdb) {
				http.Error(w, "Too many requests (IP)", http.StatusTooManyRequests)
				return
			}

			// rate limit по email
			if req.Email != "" {
				if !allow(r.Context(), "reset:email:"+req.Email, 3, time.Hour, rdb) {
					http.Error(w, "Too many requests (email)", http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func allow(ctx context.Context, key string, limit int, window time.Duration, rdb *redis.Client) bool {

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	if count == 1 {
		_ = rdb.Expire(ctx, key, window).Err()
	}

	return count <= int64(limit)
}

func getIP(r *http.Request) string {

	ip := r.Header.Get("X-Forwarded-For")

	if ip != "" {
		return strings.Split(ip, ",")[0]
	}

	return strings.Split(r.RemoteAddr, ":")[0]
}