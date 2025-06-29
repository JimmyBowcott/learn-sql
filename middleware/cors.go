package middleware

import (
	"net/http"
	"os"
)

func CorsMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != os.Getenv("FRONTEND_HOSTNAME") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_HOSTNAME"))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h(w, r)
	}
}
