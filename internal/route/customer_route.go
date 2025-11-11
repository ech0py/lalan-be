package route

import (
	"lalan-be/internal/handler"
	"lalan-be/internal/middleware"
	"net/http"
)

/*
Mendaftarkan route autentikasi customer.
Menyiapkan endpoint untuk registrasi, login, profil, dan test terproteksi dengan middleware JWT.
*/
func CustomerRoutes(h *handler.CustomerHandler) {
	v1 := "/v1"

	detailHandler := middleware.JWTMiddleware(http.HandlerFunc(h.GetCustomerProfile))

	http.HandleFunc(v1+"/customer/register", h.RegisterCustomer)
	http.HandleFunc(v1+"/customer/login", h.LoginCustomer)

	http.Handle(v1+"/customer/detail", detailHandler)

	testHandler := middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		if userID == "" {
			http.Error(w, "Token gagal diverifikasi", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Token VALID! Kamu login sebagai customer:", "user_id": "` + userID + `"}`))
	}))
	http.Handle(v1+"/customer/test-protected", testHandler)
}
