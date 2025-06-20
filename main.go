package main

import (
	"log"
	"monitoring-service/config"
	"monitoring-service/middleware"
	"monitoring-service/routes"
	"monitoring-service/service"
	"strconv"

	storageService "github.com/SIM-MBKM/filestorage/storage"
	baseServiceHelpers "github.com/SIM-MBKM/mod-service/src/helpers"
	securityMiddleware "github.com/SIM-MBKM/mod-service/src/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	baseServiceHelpers.LoadEnv()
	cfg := config.LoadConfig()

	// security
	securityKeyService := baseServiceHelpers.GetEnv("APP_KEY", "secret")
	expireSeconds, _ := strconv.ParseInt(baseServiceHelpers.GetEnv("APP_KEY_EXPIRE_SECONDS", "9999"), 10, 64)

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
	userManagementBaseURI := baseServiceHelpers.GetEnv("USER_MANAGEMENT_BASE_URI", "http://localhost:8080")
	registrationBaseURI := baseServiceHelpers.GetEnv("REGISTRATION_MANAGEMENT_BASE_URI", "http://localhost:8081")
	brokerBaseURI := baseServiceHelpers.GetEnv("BROKER_BASE_URI", "http://localhost:8082")
	port := baseServiceHelpers.GetEnv("GOLANG_PORT", "8088")

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
		config.BrokerbaseURI(brokerBaseURI),
		config.RegistrationManagementbaseURI(registrationBaseURI),
		[]string{"/async"},
	)
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}

	// Setup Gin router
	router := gin.Default()

	frontendConfig := securityMiddleware.FrontendConfig{
		AllowedOrigins:    cfg.FrontendAllowedOrigins,
		AllowedReferers:   cfg.FrontendAllowedReferers,
		RequireOrigin:     cfg.FrontendRequireOrigin,
		BypassForBrowsers: cfg.FrontendBypassBrowsers,
		CustomHeader:      cfg.FrontendCustomHeader,
		CustomHeaderValue: cfg.FrontendCustomHeaderValue,
	}
	// add cors
	router.Use(middleware.CORS())
	router.Use(securityMiddleware.AccessKeyMiddleware(securityKeyService, expireSeconds, &frontendConfig))

	// Setup routes for all controllers
	routes.ReportRoutes(router, app.ReportController, *userManagementService)
	routes.ReportScheduleRoutes(router, app.ReportScheduleController, *userManagementService)
	routes.TranscriptRoutes(router, app.TranscriptController, *userManagementService)
	routes.SyllabusRoutes(router, app.SyllabusController, *userManagementService)

	// Start server
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
