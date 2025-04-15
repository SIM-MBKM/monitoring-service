package main

import (
	"log"
	"monitoring-service/config"
	"monitoring-service/routes"
	"monitoring-service/service"

	storageService "github.com/SIM-MBKM/filestorage/storage"
	"github.com/SIM-MBKM/mod-service/src/helpers"
	baseServiceHelpers "github.com/SIM-MBKM/mod-service/src/helpers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	baseServiceHelpers.LoadEnv()

	// Initialize database connection
	db := config.SetupDatabaseConnection()

	// Initialize file storage
	storageConfig, err := storageService.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to initialize file storage: %v", err)
	}

	cache := storageService.NewMemoryCache()
	tokenManager := storageService.NewCacheTokenManager(storageConfig, cache)

	// Get environment variables
	userManagementBaseURI := helpers.GetEnv("USER_MANAGEMENT_BASE_URI", "http://localhost:8080")
	port := helpers.GetEnv("GOLANG_PORT", "8088")

	// Initialize user management service for routes
	userManagementService := service.NewUserManagementService(
		userManagementBaseURI,
		[]string{"/async"},
	)

	// Initialize application with Wire dependency injection
	app, err := InitializeAPI(
		db,
		storageConfig,
		tokenManager,
		userManagementBaseURI,
		[]string{"/async"},
	)
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}

	// Setup Gin router
	router := gin.Default()

	// Setup routes for all controllers
	routes.ReportRoutes(router, app.ReportController, *userManagementService)
	routes.ReportScheduleRoutes(router, app.ReportScheduleController, *userManagementService)

	// Start server
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
