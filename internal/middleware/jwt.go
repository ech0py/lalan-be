package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"lalan-be/internal/config"
	"lalan-be/internal/response"
)

/*
Konstanta untuk kunci konteks.
Konstanta ini mendefinisikan kunci untuk menyimpan user ID dan role dalam konteks.
*/
const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

/*
Type untuk kunci konteks.
Type ini digunakan sebagai kunci untuk nilai konteks.
*/
type contextKey string

/*
Struktur untuk claims JWT.
Struktur ini berisi claims JWT standar dan role pengguna.
*/
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

/*
Fungsi untuk middleware JWT.
Middleware ini memvalidasi token JWT dan memperbarui konteks.
*/
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek header authorization
		auth := r.Header.Get("Authorization")
		if auth == "" {
			response.Unauthorized(w, "Token required")
			return
		}

		parts := strings.Split(auth, " ")
		// Cek format Bearer
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(w, "Invalid token format")
			return
		}

		claims := &Claims{}
		// Parse token
		token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
			return config.GetJWTSecret(), nil
		})

		// Cek validitas token
		if err != nil || !token.Valid {
			response.Unauthorized(w, "Invalid or expired token")
			return
		}

		// Set context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*
Fungsi untuk mendapatkan user ID dari konteks.
User ID dikembalikan sebagai string.
*/
func GetUserID(r *http.Request) string {
	id, _ := r.Context().Value(UserIDKey).(string)
	return id
}

/*
Fungsi untuk mendapatkan user role dari konteks.
User role dikembalikan sebagai string.
*/
func GetUserRole(r *http.Request) string {
	role, _ := r.Context().Value(UserRoleKey).(string)
	return role
}
