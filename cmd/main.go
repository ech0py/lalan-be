package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lalan-be/internal/config"
	"lalan-be/internal/features/admin"
	"lalan-be/internal/features/public"
	"lalan-be/internal/middleware"

	"github.com/gorilla/mux"
)

/*
Fungsi utama aplikasi.
Menginisialisasi konfigurasi, database, dependency injection,
routing, dan server HTTP dengan graceful shutdown.
*/
func main() {
	// Muat variabel lingkungan
	config.LoadEnv()

	// Inisialisasi koneksi database
	cfg, err := config.DatabaseConfig()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	db := cfg.DB
	defer db.Close()

	log.Printf(
		"Database connected → host=%s port=%d db=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Dependency injection untuk admin
	aRepo := admin.NewAdminRepository(db)
	aService := admin.NewAdminService(aRepo)
	aHandler := admin.NewAdminHandler(aService)

	// Dependency injection untuk public
	pRepo := public.NewPublicRepository(db)
	pService := public.NewPublicService(pRepo)
	pHandler := public.NewPublicHandler(pService)

	// Setup routing
	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)

	admin.SetupAdminRoutes(router, aHandler)
	public.SetupPublicRoutes(router, pHandler)

	// Buat server dengan graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Channel sinyal untuk shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server running at http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-c
	log.Println("Shutting down server...")

	// Kontekst graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
