package main

import (
	"net/http"
	"shadeless-api/main/config"
)

// Middleware function, which will be called for each request
func middlewareSetHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		whitelistUrl := config.GetInstance().GetFrontendUrl()
		w.Header().Set("Access-Control-Allow-Origin", whitelistUrl)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Vary", "Access-Control-Request-Method")
		w.Header().Set("Vary", "Access-Control-Request-Headers")

		next.ServeHTTP(w, r)
	})
}

func middlewareOptions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
