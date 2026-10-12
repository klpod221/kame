package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type clientBucket struct {
	tokens     float64
	lastRefill time.Time
}

// IPRateLimiter implements a thread-safe token bucket rate limiter per client IP.
type IPRateLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientBucket
	rate    float64 // tokens added per second
	burst   float64 // max token capacity
}

// NewRateLimiter creates a new rate limiter allowing requestsPerMinute with specified burst capacity.
func NewRateLimiter(requestsPerMinute, burst int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		clients: make(map[string]*clientBucket),
		rate:    float64(requestsPerMinute) / 60.0,
		burst:   float64(burst),
	}

	// Background routine to evict inactive IP buckets every 5 minutes
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			limiter.cleanup(10 * time.Minute)
		}
	}()

	return limiter
}

func (l *IPRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	bucket, exists := l.clients[ip]
	if !exists {
		l.clients[ip] = &clientBucket{
			tokens:     l.burst - 1,
			lastRefill: now,
		}
		return true
	}

	// Refill tokens based on elapsed duration
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * l.rate
	if bucket.tokens > l.burst {
		bucket.tokens = l.burst
	}
	bucket.lastRefill = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

func (l *IPRateLimiter) cleanup(maxIdle time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := time.Now().Add(-maxIdle)
	for ip, b := range l.clients {
		if b.lastRefill.Before(cutoff) {
			delete(l.clients, ip)
		}
	}
}

// Handler returns a Chi middleware handler for the rate limiter.
func (l *IPRateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		if !l.allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "30")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Rate limit exceeded. Please slow down.",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header from reverse proxy
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		clientIP := strings.TrimSpace(parts[0])
		if clientIP != "" {
			return clientIP
		}
	}

	// Check X-Real-IP header
	xrip := r.Header.Get("X-Real-IP")
	if xrip != "" {
		return strings.TrimSpace(xrip)
	}

	// Fallback to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
