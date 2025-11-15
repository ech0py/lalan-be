package hoster

import (
	"lalan-be/internal/middleware"

	"github.com/gorilla/mux"
)

/*
	Mengatur routing untuk fitur hoster.

Router dikonfigurasi dengan endpoint publik dan terproteksi.
*/
func SetupHosterRoutes(router *mux.Router, h *HosterHandler) {
	hoster := router.PathPrefix("/api/v1/hoster").Subrouter()
	// Setup public routes
	hoster.HandleFunc("/register", h.CreateHoster).Methods("POST")
	hoster.HandleFunc("/login", h.LoginHoster).Methods("POST")
	protected := hoster.PathPrefix("").Subrouter()

	// Middleware JWT
	protected.Use(middleware.JWTMiddleware)

	// Middleware hoster only
	protected.Use(middleware.Hoster)
	// Setup protected routes

	protected.HandleFunc("/detail", h.GetDetailHoster).Methods("GET")

}
