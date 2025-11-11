package main

import (
	"lalan-be/internal/config"
	"lalan-be/internal/handler"
	"lalan-be/internal/repository"
	"lalan-be/internal/route"
	"lalan-be/internal/service"
	"log"
	"net/http"
	"os"
)

/*
Menginisialisasi konfigurasi database, repositori, layanan, handler, dan rute untuk aplikasi.
Memulai server HTTP di port yang ditentukan.
*/
func main() {
	cfg, err := config.DatabaseConfig()
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer cfg.DB.Close()

	authRepo := repository.NewAuthRepository(cfg.DB)
	authService := service.NewHosterService(authRepo)
	authHandler := handler.NewHosterHandler(authService)
	route.HosterRoutes(authHandler)

	customerRepo := repository.NewCustomerRepository(cfg.DB)
	customerService := service.NewCustomerService(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerService)
	route.CustomerRoutes(customerHandler)

	categoryRepo := repository.NewCategoryRepository(cfg.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	route.CategoryRoutes(categoryHandler)

	itemRepo := repository.NewItemRepository(cfg.DB)
	itemService := service.NewItemService(itemRepo)
	itemHandler := handler.NewItemHandler(itemService)
	route.ItemRoutes(itemHandler)

	tncRepo := repository.NewTermsAndConditionsRepository(cfg.DB)
	tncService := service.NewTermsAndConditionsService(tncRepo)
	tncHandler := handler.NewTermsAndConditionsHandler(tncService)
	route.TermsAndConditionsRoutes(tncHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
