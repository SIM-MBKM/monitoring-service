//go:build wireinject
// +build wireinject

package main

import (
	"monitoring-service/config"
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
	TranscriptController     controller.TranscriptController
	SyllabusController       controller.SyllabusController
}

func newApplication(
	reportController controller.ReportController,
	reportScheduleController controller.ReportScheduleController,
	transcriptController controller.TranscriptController,
	syllabusController controller.SyllabusController,
) *Application {
	return &Application{
		ReportController:         reportController,
		ReportScheduleController: reportScheduleController,
		TranscriptController:     transcriptController,
		SyllabusController:       syllabusController,
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

func ProvideTranscriptRepository(db *gorm.DB) repository.TranscriptRepository {
	return repository.NewTranscriptRepository(db)
}

func ProvideSyllabusRepository(db *gorm.DB) repository.SyllabusRepository {
	return repository.NewSyllabusRepository(db)
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
	registrationBaseURI config.RegistrationManagementbaseURI,
	asyncURIs []string,
) service.ReportScheduleService {
	return service.NewReportScheduleService(reportScheduleRepo, userManagementBaseURI, string(registrationBaseURI), asyncURIs)
}

func ProvideTranscriptService(
	transcriptRepo repository.TranscriptRepository,
	userManagementBaseURI string,
	registrationBaseURI config.RegistrationManagementbaseURI,
	asyncURIs []string,
	config *storageService.Config,
	tokenManager *storageService.CacheTokenManager,
) service.TranscriptService {
	return service.NewTranscriptService(
		transcriptRepo,
		userManagementBaseURI,
		string(registrationBaseURI),
		asyncURIs,
		config,
		tokenManager,
	)
}

func ProvideSyllabusService(
	syllabusRepo repository.SyllabusRepository,
	userManagementBaseURI string,
	registrationBaseURI config.RegistrationManagementbaseURI,
	asyncURIs []string,
	config *storageService.Config,
	tokenManager *storageService.CacheTokenManager,
) service.SyllabusService {
	return service.NewSyllabusService(
		syllabusRepo,
		userManagementBaseURI,
		string(registrationBaseURI),
		asyncURIs,
		config,
		tokenManager,
	)
}

// Controller providers
func ProvideReportController(reportService service.ReportService) controller.ReportController {
	return *controller.NewReportController(reportService)
}

func ProvideReportScheduleController(reportScheduleService service.ReportScheduleService) controller.ReportScheduleController {
	return *controller.NewReportScheduleController(reportScheduleService)
}

func ProvideTranscriptController(transcriptService service.TranscriptService) controller.TranscriptController {
	return *controller.NewTranscriptController(transcriptService)
}

func ProvideSyllabusController(syllabusService service.SyllabusService) controller.SyllabusController {
	return *controller.NewSyllabusController(syllabusService)
}

// Provider sets
var (
	RepositorySet = wire.NewSet(
		ProvideBaseRepository,
		ProvideReportRepository,
		ProvideReportScheduleRepository,
		ProvideTranscriptRepository,
		ProvideSyllabusRepository,
	)

	ServiceSet = wire.NewSet(
		ProvideFileService,
		ProvideUserManagementService,
		ProvideReportService,
		ProvideReportScheduleService,
		ProvideTranscriptService,
		ProvideSyllabusService,
	)

	ControllerSet = wire.NewSet(
		ProvideReportController,
		ProvideReportScheduleController,
		ProvideTranscriptController,
		ProvideSyllabusController,
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
	registrationBaseURI config.RegistrationManagementbaseURI,
	asyncURIs []string,
) (*Application, error) {
	wire.Build(
		AllSet,
		newApplication,
	)

	return nil, nil
}
