// Package middlewares provides HTTP middleware functions for request handling.
package middlewares

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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
func PerClientRateLimit(ipLimiter *IPRateLimiter, next http.Handler, httpRequest metric.Int64Counter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		limiter := ipLimiter.GetLimiter(ip)
		if !limiter.Allow() {
			log.Println("[middleware/ratelimit] limit exceed")
			http.Error(w, "limit exceed", http.StatusTooManyRequests)
			return
		}
		httpRequest.Add(r.Context(), 1, metric.WithAttributes(
			attribute.String("route", r.URL.Path),
			attribute.String("method", r.Method),
		))
		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the client IP address from request headers or remote address.
func getClientIP(r *http.Request) (ip string) {
	for name, values := range r.Header {
		if strings.ToLower(name) == "cf-connecting-ip" {
			ip = searchIp(values, "cloudflare")
		}
		if strings.ToLower(name) == "x-forwarded-for" {
			ip = searchIp(values, "caddy")
		}
	}
	// TODO: move to searchIp
	if ip == "" {
		ipLocal, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		} else {
			ip = ipLocal
		}
	}
	log.Printf("[api/middleware] get ip for dev %s\n", ip)
	return ip
}

// searchIp parses and returns the IP address from header values based on the source platform.
func searchIp(values []string, target string) string {
	switch target {
	case "cloudflare":
		for _, ip := range values {
			if ip != "" {
				return strings.TrimSpace(ip)
			}
		}
	case "caddy":
		for _, xff := range values {
			if xff != "" {
				ips := strings.Split(xff, ",")
				return strings.TrimSpace(ips[0])
			}
		}

	}
	return ""
}
