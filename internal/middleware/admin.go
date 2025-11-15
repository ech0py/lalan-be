package middleware

import (
	"net/http"

	"lalan-be/internal/response"
)

// Fungsi middleware untuk akses admin only.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek role admin
		if GetUserRole(r) != "admin" {
			response.Forbidden(w, "Admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
