// Package middlewares provides HTTP middleware functions for request handling.
package middlewares

import (
	"log"
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimit defines an interface for retrieving rate limiters by IP.
type RateLimit interface {
	GetLimiter(ip string) *rate.Limit
}

// IPRateLimiter manages rate limiters for IP addresses in a thread-safe manner.
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.Mutex
	r   rate.Limit
	b   int
}

// GetLimiter returns the rate.Limiter for a given IP address, creating one if necessary.
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}
	return limiter
}

// NewIPRateLimiter constructs a new IPRateLimiter with the specified rate and burst limits.
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
}

// PerClientRateLimit returns an HTTP handler middleware that enforces rate limits per client IP.
func PerClientRateLimit(ipLimiter *IPRateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("[middleware/ratelimit] error: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		limiter := ipLimiter.GetLimiter(ip)
		if !limiter.Allow() {
			log.Println("[middleware/ratelimit] limit exceed")
			http.Error(w, "limit exceed", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
