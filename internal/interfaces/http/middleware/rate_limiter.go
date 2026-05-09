package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimit struct {
	Requests int
	Window   time.Duration
	Burst    int
}

type RouteRateLimit struct {
	Prefix string
	Limit  RateLimit
}

type rateLimitConfig struct {
	defaultLimit RateLimit
	routeLimits  []RouteRateLimit
}

type rateLimitOption func(*rateLimitConfig)

func WithRouteRateLimit(prefix string, limit RateLimit) rateLimitOption {
	return func(config *rateLimitConfig) {
		if prefix == "" {
			return
		}
		config.routeLimits = append(config.routeLimits, RouteRateLimit{
			Prefix: prefix,
			Limit:  limit,
		})
	}
}

func RateLimiter(defaultLimit RateLimit, options ...rateLimitOption) func(http.Handler) http.Handler {
	config := rateLimitConfig{defaultLimit: defaultLimit}
	for _, option := range options {
		option(&config)
	}

	clients := make(map[string]*rate.Limiter)
	var mu sync.Mutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeLimit := config.limitForPath(r.URL.Path)
			if routeLimit.Requests <= 0 || routeLimit.Window <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			key := clientKey(clientIP(r), r.URL.Path, config.routeLimits)

			mu.Lock()
			limiter, ok := clients[key]
			if !ok {
				limiter = newLimiter(routeLimit)
				clients[key] = limiter
			}
			allowed := limiter.Allow()
			mu.Unlock()

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "rate limit exceeded",
					"code":  "RATE_LIMIT_EXCEEDED",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (config rateLimitConfig) limitForPath(path string) RateLimit {
	for _, routeLimit := range config.routeLimits {
		if strings.HasPrefix(path, routeLimit.Prefix) {
			return routeLimit.Limit
		}
	}
	return config.defaultLimit
}

func clientKey(ip, path string, routeLimits []RouteRateLimit) string {
	for _, routeLimit := range routeLimits {
		if strings.HasPrefix(path, routeLimit.Prefix) {
			return ip + "|" + routeLimit.Prefix
		}
	}
	return ip + "|default"
}

func newLimiter(limit RateLimit) *rate.Limiter {
	burst := limit.Burst
	if burst <= 0 {
		burst = limit.Requests
	}
	return rate.NewLimiter(rate.Every(limit.Window/time.Duration(limit.Requests)), burst)
}

func clientIP(r *http.Request) string {
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		if ip := strings.TrimSpace(strings.Split(forwardedFor, ",")[0]); ip != "" {
			return ip
		}
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
