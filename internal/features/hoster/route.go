package hoster

import (
	"github.com/gorilla/mux"
)

/*
Fungsi untuk mengatur rute fitur hoster.
Router dikonfigurasi dengan rute yang diperlukan.
*/
func SetupHosterRoutes(router *mux.Router, handler *HosterHandler) {
	hoster := router.PathPrefix("/api/v1/hoster").Subrouter()
	hoster.HandleFunc("/register", handler.CreateHoster).Methods("POST")
	hoster.HandleFunc("/login", handler.LoginHoster).Methods("POST")
	hoster.HandleFunc("/detail", handler.GetDetailHoster).Methods("GET")
	hoster.HandleFunc("/items", handler.CreateItem).Methods("POST")
	hoster.HandleFunc("/items/{id}", handler.GetItemByID).Methods("GET")
	hoster.HandleFunc("/items", handler.GetAllItems).Methods("GET")
	hoster.HandleFunc("/items/{id}", handler.UpdateItem).Methods("PUT")
	hoster.HandleFunc("/items/{id}", handler.DeleteItem).Methods("DELETE")
	hoster.HandleFunc("/terms", handler.CreateTermsAndConditions).Methods("POST")
	hoster.HandleFunc("/terms/{id}", handler.FindTermsAndConditionsByID).Methods("GET")
	hoster.HandleFunc("/terms", handler.GetAllTermsAndConditions).Methods("GET")
	hoster.HandleFunc("/terms", handler.UpdateTermsAndConditions).Methods("PUT")
	hoster.HandleFunc("/terms", handler.DeleteTermsAndConditions).Methods("DELETE")
}
