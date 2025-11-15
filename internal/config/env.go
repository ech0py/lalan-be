package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Fungsi untuk mengambil rahasia JWT.
func GetJWTSecret() []byte {
	secret := GetEnv("JWT_SECRET", "tesingdev")
	return []byte(secret)
}

// Flag untuk status pemuatan environment.
var envLoaded bool

// Fungsi untuk memuat environment.
func LoadEnv() {
	// Cek status pemuatan
	if envLoaded {
		return
	}

	// Muat .env.dev jika bukan produksi
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load(".env.dev")
	}

	envLoaded = true
}

// Fungsi untuk mengambil environment dengan fallback.
func GetEnv(key, fallback string) string {
	LoadEnv()

	// Return nilai jika ada
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Fungsi untuk mengambil environment wajib.
func MustGetEnv(key string) string {
	LoadEnv()

	v := os.Getenv(key)
	// Fatal jika kosong
	if v == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return v
}
