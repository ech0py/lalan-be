package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"lalan-be/internal/config"
	"lalan-be/internal/response"
)

// Type untuk context key.
type contextKey string

// Struct untuk claims JWT.
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// Konstanta untuk context keys.
const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

// Fungsi middleware untuk validasi JWT.
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

// Fungsi untuk dapatkan user ID dari context.
func GetUserID(r *http.Request) string {
	id, _ := r.Context().Value(UserIDKey).(string)
	return id
}

// Fungsi untuk dapatkan user role dari context.
func GetUserRole(r *http.Request) string {
	role, _ := r.Context().Value(UserRoleKey).(string)
	return role
}
