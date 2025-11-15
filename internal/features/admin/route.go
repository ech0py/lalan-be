package admin

import (
	"github.com/gorilla/mux"

	"lalan-be/internal/middleware"
)

/*
	Mengatur routing untuk fitur admin.

Router dikonfigurasi dengan endpoint publik dan terproteksi.
*/
func SetupAdminRoutes(router *mux.Router, h *AdminHandler) {
	// Setup group admin
	admin := router.PathPrefix("/api/v1/admin").Subrouter()

	// Setup public routes
	admin.HandleFunc("/register", h.CreateAdmin).Methods("POST")
	admin.HandleFunc("/login", h.LoginAdmin).Methods("POST")

	// Setup protected routes
	protected := admin.PathPrefix("").Subrouter()

	// Middleware JWT
	protected.Use(middleware.JWTMiddleware)

	// Middleware admin only
	protected.Use(middleware.Admin)

	// Endpoint protected
	protected.HandleFunc("/category/create", h.CreateCategory).Methods("POST")
	protected.HandleFunc("/category/update", h.UpdateCategory).Methods("PUT")
	protected.HandleFunc("/category/delete", h.DeleteCategory).Methods("DELETE")
}
