package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// API Gateway configuration
	routes := map[string]string{
		"/api/users/":         "http://user-service:8081",
		"/api/auth/":          "http://auth-service:8082",
		"/api/payments/":      "http://payment-service:8083",
		"/api/notifications/": "http://notification-service:8084",
		"/api/analytics/":     "http://analytics-service:8085",
	}

	// Set up handlers for each route
	for path, targetURL := range routes {
		target, err := url.Parse(targetURL)
		if err != nil {
			log.Fatal(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Proxying request: %s -> %s", r.URL.Path, targetURL)
			proxy.ServeHTTP(w, r)
		})
	}

	// Serve static frontend content
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Start server
	log.Println("API Gateway starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
