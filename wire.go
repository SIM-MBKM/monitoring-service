//go:build wireinject
// +build wireinject

package main

import (
	"monitoring-service/controller"
	"monitoring-service/repository"
	"monitoring-service/service"

	storageService "github.com/SIM-MBKM/filestorage/storage"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// Application struct to hold all controllers
type Application struct {
	ReportController         controller.ReportController
	ReportScheduleController controller.ReportScheduleController
}

func newApplication(
	reportController controller.ReportController,
	reportScheduleController controller.ReportScheduleController,
) *Application {
	return &Application{
		ReportController:         reportController,
		ReportScheduleController: reportScheduleController,
	}
}

// Repository providers
func ProvideBaseRepository(db *gorm.DB) repository.BaseRepository {
	return repository.NewBaseRepository(db)
}

func ProvideReportRepository(db *gorm.DB) repository.ReportRepository {
	return repository.NewReportRepository(db)
}

func ProvideReportScheduleRepository(db *gorm.DB) repository.ReportScheduleReposiotry {
	return repository.NewReportScheduleRepository(db)
}

// Service providers
func ProvideFileService(config *storageService.Config, tokenManager *storageService.CacheTokenManager) *service.FileService {
	return service.NewFileService(config, tokenManager)
}

func ProvideUserManagementService(
	userManagementBaseURI string,
	asyncURIs []string,
) *service.UserManagementService {
	return service.NewUserManagementService(userManagementBaseURI, asyncURIs)
}

func ProvideReportService(
	reportRepo repository.ReportRepository,
	reportScheduleRepo repository.ReportScheduleReposiotry,
	userManagementBaseURI string,
	asyncURIs []string,
	config *storageService.Config,
	tokenManager *storageService.CacheTokenManager,
) service.ReportService {
	return service.NewReportService(
		reportRepo,
		reportScheduleRepo,
		userManagementBaseURI,
		asyncURIs,
		config,
		tokenManager,
	)
}

func ProvideReportScheduleService(
	reportScheduleRepo repository.ReportScheduleReposiotry,
	userManagementBaseURI string,
	asyncURIs []string,
) service.ReportScheduleService {
	return service.NewReportScheduleService(reportScheduleRepo, userManagementBaseURI, asyncURIs)
}

// Controller providers
func ProvideReportController(reportService service.ReportService) controller.ReportController {
	return *controller.NewReportController(reportService)
}

func ProvideReportScheduleController(reportScheduleService service.ReportScheduleService) controller.ReportScheduleController {
	return *controller.NewReportScheduleController(reportScheduleService)
}

// Provider sets
var (
	RepositorySet = wire.NewSet(
		ProvideBaseRepository,
		ProvideReportRepository,
		ProvideReportScheduleRepository,
	)

	ServiceSet = wire.NewSet(
		ProvideFileService,
		ProvideUserManagementService,
		ProvideReportService,
		ProvideReportScheduleService,
	)

	ControllerSet = wire.NewSet(
		ProvideReportController,
		ProvideReportScheduleController,
	)

	AllSet = wire.NewSet(
		RepositorySet,
		ServiceSet,
		ControllerSet,
	)
)

// InitializeAPI creates the full application with all dependencies
func InitializeAPI(
	db *gorm.DB,
	config *storageService.Config,
	tokenManager *storageService.CacheTokenManager,
	userManagementBaseURI string,
	asyncURIs []string,
) (*Application, error) {
	wire.Build(
		AllSet,
		newApplication,
	)

	return nil, nil
}
