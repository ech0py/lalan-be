package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
	Variabel untuk menandai status pemuatan environment.

Menunjukkan apakah environment sudah dimuat.
*/
var envLoaded bool

/*
	Mengambil rahasia JWT dari environment.

Rahasia JWT dikembalikan sebagai byte slice.
*/
func GetJWTSecret() []byte {
	secret := GetEnv("JWT_SECRET", "tesingdev")
	return []byte(secret)
}

/*
	Memuat file environment jika belum dimuat.

Environment dimuat dari file .env.dev jika bukan produksi.
*/
func LoadEnv() {
	if envLoaded {
		return
	}
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load(".env.dev")
	}
	envLoaded = true
}

/*
	Mengambil nilai environment dengan fallback.

Nilai environment atau fallback dikembalikan.
*/
func GetEnv(key, fallback string) string {
	LoadEnv()
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

/*
	Mengambil nilai environment yang wajib.

Nilai environment dikembalikan atau program dihentikan.
*/
func MustGetEnv(key string) string {
	LoadEnv()
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return v
}
