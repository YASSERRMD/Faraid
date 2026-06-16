package api

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// secureHeaders adds defensive HTTP security headers to every response.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Content-Security-Policy", "default-src 'self'")
		h.Set("Referrer-Policy", "same-origin")
		h.Set("Permissions-Policy", "interest-cohort=()")
		next.ServeHTTP(w, r)
	})
}

// requestSizeLimit limits each request body to max bytes. Requests with larger
// bodies receive 413 and the oversized body is discarded.
func requestSizeLimit(max int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, max)
			next.ServeHTTP(w, r)
		})
	}
}

// ipRateLimiter implements a fixed-window per-IP rate limiter. The window map
// is replaced atomically each minute, bounding memory to one entry per active
// IP address.
type ipRateLimiter struct {
	mu          sync.Mutex
	counts      map[string]int
	windowStart time.Time
	limit       int
}

func newIPRateLimiter(rpm int) *ipRateLimiter {
	return &ipRateLimiter{
		counts:      make(map[string]int),
		windowStart: time.Now(),
		limit:       rpm,
	}
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.windowStart) >= time.Minute {
		l.counts = make(map[string]int)
		l.windowStart = time.Now()
	}
	l.counts[ip]++
	return l.counts[ip] <= l.limit
}

// rateLimit returns a middleware that limits each client IP to rpm requests per
// minute. The client IP is taken from X-Forwarded-For if present, otherwise
// from RemoteAddr.
func rateLimit(rpm int) func(http.Handler) http.Handler {
	limiter := newIPRateLimiter(rpm)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				if i := strings.Index(xff, ","); i >= 0 {
					ip = strings.TrimSpace(xff[:i])
				} else {
					ip = strings.TrimSpace(xff)
				}
			}
			if !limiter.allow(ip) {
				writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
