package route

import (
	"net/http"

	"lalan-be/internal/handler"
	"lalan-be/internal/middleware"
)

/*
Mendaftarkan route autentikasi hoster.
Menyiapkan endpoint untuk registrasi, login, profil, dan test terproteksi dengan middleware JWT.
*/
func HosterRoutes(h *handler.HosterHandler) {
	v1 := "/v1"

	detailHandler := middleware.JWTMiddleware(http.HandlerFunc(h.GetHosterProfile))

	http.HandleFunc(v1+"/auth/register", h.Register)
	http.HandleFunc(v1+"/auth/login", h.LoginHoster)

	http.Handle(v1+"/auth/detail", detailHandler)

	testHandler := middleware.JWTMiddleware(http.HandlerFunc(h.TestProtected))
	http.Handle(v1+"/test-protected", testHandler)
}
