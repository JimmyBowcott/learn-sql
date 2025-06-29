package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	rateLimit      = 30
	rateInterval   = time.Minute
	clients        = make(map[string]*clientData)
	clientsMu      sync.Mutex
	cleanupTimeout = 10 * time.Minute
)

type clientData struct {
	tokens       int
	lastRefill   time.Time
	lastActivity time.Time
}

func RateLimitMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		clientsMu.Lock()
		client, exists := clients[ip]
		if !exists {
			client = &clientData{tokens: rateLimit, lastRefill: time.Now(), lastActivity: time.Now()}
			clients[ip] = client
		}
		now := time.Now()
		elapsed := now.Sub(client.lastRefill)
		refill := int(elapsed.Minutes() * float64(rateLimit))
		if refill > 0 {
			client.tokens = min(rateLimit, client.tokens+refill)
			client.lastRefill = now
		}
		if client.tokens <= 0 {
			clientsMu.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		client.tokens--
		client.lastActivity = now
		clientsMu.Unlock()

		h(w, r)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func CleanupInactiveClients() {
	for {
		time.Sleep(time.Minute)
		clientsMu.Lock()
		now := time.Now()
		for ip, client := range clients {
			if now.Sub(client.lastActivity) > cleanupTimeout {
				delete(clients, ip)
			}
		}
		clientsMu.Unlock()
	}
}

