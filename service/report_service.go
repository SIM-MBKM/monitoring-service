package service

import (
	"context"
	"mime/multipart"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"monitoring-service/repository"
	"reflect"
	"time"

	storageService "github.com/SIM-MBKM/filestorage/storage"
	"github.com/google/uuid"
)

type reportService struct {
	reportRepo  repository.ReportRepository
	fileService *FileService
}

type ReportService interface {
	Index(ctx context.Context) ([]dto.ReportResponse, error)
	Create(ctx context.Context, report dto.ReportRequest, file *multipart.FileHeader) (dto.ReportResponse, error)
	Update(ctx context.Context, id string, report dto.ReportRequest) error
	FindByID(ctx context.Context, id string) (dto.ReportResponse, error)
	Destroy(ctx context.Context, id string) error
	FindByReportScheduleID(ctx context.Context, reportScheduleID string) ([]dto.ReportResponse, error)
}

func NewReportService(reportRepo repository.ReportRepository, config *storageService.Config, tokenManager *storageService.CacheTokenManager) ReportService {
	return &reportService{
		reportRepo:  reportRepo,
		fileService: NewFileService(config, tokenManager),
	}
}

// Index retrieves all reports
func (s *reportService) Index(ctx context.Context) ([]dto.ReportResponse, error) {
	reports, err := s.reportRepo.Index(ctx, nil)
	if err != nil {
		return nil, err
	}

	var reportResponses []dto.ReportResponse
	for _, report := range reports {
		reportResponses = append(reportResponses, dto.ReportResponse{
			ID:                    report.ID.String(),
			ReportScheduleID:      report.ReportScheduleID,
			FileStorageID:         report.FileStorageID,
			Title:                 report.Title,
			Content:               report.Content,
			ReportType:            report.ReportType,
			Feedback:              report.Feedback,
			AcademicAdvisorStatus: report.AcademicAdvisorStatus,
		})
	}

	return reportResponses, nil
}

// Create creates a new report
func (s *reportService) Create(ctx context.Context, report dto.ReportRequest, file *multipart.FileHeader) (dto.ReportResponse, error) {
	// Generate new UUID for the report
	result, err := s.fileService.storage.GcsUpload(file, "sim_mbkm", "report", "report")
	if err != nil {
		return dto.ReportResponse{}, err
	}
	var reportEntity entity.Report
	reportEntity.ID = uuid.New()
	reportEntity.FileStorageID = result.FileID
	reportEntity.Title = report.Title
	reportEntity.Content = report.Content
	reportEntity.ReportType = report.ReportType
	reportEntity.Feedback = report.Feedback

	// Set timestamps
	now := time.Now()
	reportEntity.CreatedAt = &now
	reportEntity.UpdatedAt = &now

	reportResponse, err := s.reportRepo.Create(ctx, reportEntity, nil)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	return dto.ReportResponse{
		ID:                    reportResponse.ID.String(),
		ReportScheduleID:      reportResponse.ReportScheduleID,
		FileStorageID:         reportResponse.FileStorageID,
		Title:                 reportResponse.Title,
		Content:               reportResponse.Content,
		ReportType:            reportResponse.ReportType,
		Feedback:              reportResponse.Feedback,
		AcademicAdvisorStatus: reportResponse.AcademicAdvisorStatus,
	}, nil
}

// Update updates an existing report
func (s *reportService) Update(ctx context.Context, id string, subject dto.ReportRequest) error {
	res, err := s.reportRepo.FindByID(ctx, id, nil)
	if err != nil {
		return err
	}

	// Create programTypeEntity with original ID
	reportEntity := entity.Report{
		ID: res.ID,
	}

	// Get reflection values
	resValue := reflect.ValueOf(res)
	reqValue := reflect.ValueOf(subject)
	entityValue := reflect.ValueOf(&reportEntity).Elem()

	// Iterate through fields of the request type
	for i := 0; i < reqValue.Type().NumField(); i++ {
		reqField := reqValue.Type().Field(i)
		reqFieldValue := reqValue.Field(i)

		// Find corresponding field in the entity
		entityField := entityValue.FieldByName(reqField.Name)

		// Find corresponding field in the original result
		resField := resValue.FieldByName(reqField.Name)

		// Check if the field exists and can be set
		if entityField.IsValid() && entityField.CanSet() {
			// If the request field is not zero, use its value
			if !reflect.DeepEqual(reqFieldValue.Interface(), reflect.Zero(reqFieldValue.Type()).Interface()) {
				entityField.Set(reqFieldValue)
			} else {
				// Otherwise, use the original value
				entityField.Set(resField)
			}
		}
	}

	// Perform the update
	err = s.reportRepo.Update(ctx, id, reportEntity, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindByID retrieves a report by its ID
func (s *reportService) FindByID(ctx context.Context, id string) (dto.ReportResponse, error) {
	report, err := s.reportRepo.FindByID(ctx, id, nil)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	return dto.ReportResponse{
		ID:                    report.ID.String(),
		ReportScheduleID:      report.ReportScheduleID,
		FileStorageID:         report.FileStorageID,
		Title:                 report.Title,
		Content:               report.Content,
		ReportType:            report.ReportType,
		Feedback:              report.Feedback,
		AcademicAdvisorStatus: report.AcademicAdvisorStatus,
	}, nil
}

// Destroy deletes a report
func (s *reportService) Destroy(ctx context.Context, id string) error {
	err := s.reportRepo.Destroy(ctx, id, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindByReportScheduleID retrieves reports by report schedule ID
func (s *reportService) FindByReportScheduleID(ctx context.Context, reportScheduleID string) ([]dto.ReportResponse, error) {
	reports, err := s.reportRepo.FindByReportScheduleID(ctx, reportScheduleID, nil)
	if err != nil {
		return nil, err
	}

	var reportResponses []dto.ReportResponse
	for _, report := range reports {
		reportResponses = append(reportResponses, dto.ReportResponse{
			ID:                    report.ID.String(),
			ReportScheduleID:      report.ReportScheduleID,
			FileStorageID:         report.FileStorageID,
			Title:                 report.Title,
			Content:               report.Content,
			ReportType:            report.ReportType,
			Feedback:              report.Feedback,
			AcademicAdvisorStatus: report.AcademicAdvisorStatus,
		})
	}

	return reportResponses, nil
}
