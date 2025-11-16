package public

import (
	"github.com/gorilla/mux"
)

/*
Fungsi untuk mengatur rute fitur publik.
Router dikonfigurasi dengan rute publik.
*/
func SetupPublicRoutes(router *mux.Router, h *PublicHandler) {
	public := router.PathPrefix("/api/v1/public").Subrouter()
	public.HandleFunc("/category", h.GetAllCategories).Methods("GET")
	public.HandleFunc("/item", h.GetAllItems).Methods("GET")
	public.HandleFunc("/tnc", h.GetAllTermsAndConditions).Methods("GET")
}
