package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
Variabel untuk status pemuatan environment.
Variabel ini menandai apakah environment sudah dimuat.
*/
var envLoaded bool

/*
Fungsi untuk mendapatkan rahasia JWT.
Rahasia JWT dikembalikan sebagai byte slice.
*/
func GetJWTSecret() []byte {
	secret := GetEnv("JWT_SECRET", "tesingdev")
	return []byte(secret)
}

/*
Fungsi untuk memuat environment.
Environment dimuat dari file jika belum dimuat.
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
Fungsi untuk mendapatkan nilai environment dengan fallback.
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
Fungsi untuk mendapatkan nilai environment yang wajib.
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
