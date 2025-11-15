package middleware

import (
	"net/http"

	"lalan-be/internal/response"
)

/*
	Memeriksa akses admin.

Request diteruskan jika role admin.
*/
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek role admin
		if GetUserRole(r) != "admin" {
			response.Forbidden(w, "Admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
	Memeriksa akses hoster.

Request diteruskan jika role hoster.
*/
func Hoster(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek role hoster
		if GetUserRole(r) != "hoster" {
			response.Forbidden(w, "Hoster access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
