package admin

import (
	"github.com/gorilla/mux"

	"lalan-be/internal/middleware"
)

// Fungsi untuk setup routing admin.
func SetupAdminRoutes(router *mux.Router, h *AdminHandler) {
	// Setup group admin
	admin := router.PathPrefix("/api/v1/admin").Subrouter()

	// Setup public routes
	admin.HandleFunc("/create", h.CreateAdmin).Methods("POST")
	admin.HandleFunc("/login", h.LoginAdmin).Methods("POST")

	// Setup protected routes
	protected := admin.PathPrefix("").Subrouter()

	// Middleware JWT
	protected.Use(middleware.JWTMiddleware)

	// Middleware admin only
	protected.Use(middleware.AdminOnly)

	// Endpoint protected
	protected.HandleFunc("/category/create", h.CreateCategory).Methods("POST")
	protected.HandleFunc("/category/update", h.UpdateCategory).Methods("PUT")
	protected.HandleFunc("/category/delete", h.DeleteCategory).Methods("DELETE")
}
