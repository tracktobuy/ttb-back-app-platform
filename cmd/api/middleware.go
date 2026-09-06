package main

import (
	"net/http"
	"strings"
)

// enableCORS lets browser clients served from the configured origins call the
// API. Preflight requests are answered here instead of by the router, which
// only registers the concrete methods each route uses and would reject OPTIONS.
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// The response depends on the request origin, so caches must not serve
		// one origin's response to another.
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" && app.isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "300")

				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) isAllowedOrigin(origin string) bool {

	for _, allowed := range app.config.CorsAllowedOrigins {
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}

	return false
}
