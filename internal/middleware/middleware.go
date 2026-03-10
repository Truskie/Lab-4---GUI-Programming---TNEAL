package middleware

import (
	"log"
	"net/http"
	"time"
)

// logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//getting timestamp
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		//request details
		log.Printf("[%s] %s %s", timestamp, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

// Custom timing middlware
func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		//checking req duration
		log.Printf("[Timing] %s %s took %v", r.Method, r.URL.Path, duration)
	})
}
