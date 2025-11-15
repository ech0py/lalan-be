package public

import (
	"github.com/gorilla/mux"
)

// Fungsi untuk setup routing public.
func SetupPublicRoutes(router *mux.Router, h *PublicHandler) {
	// Setup group public
	admin := router.PathPrefix("/api/v1/public").Subrouter()
	admin.HandleFunc("/categories", h.GetCategories).Methods("GET")
}
