package service

import (
	"context"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"monitoring-service/repository"
	"time"

	"github.com/google/uuid"
)

type reportService struct {
	reportRepo repository.ReportRepository
}

type ReportService interface {
	Index(ctx context.Context) ([]dto.ReportResponse, error)
	Create(ctx context.Context, report entity.Report) (dto.ReportResponse, error)
	Update(ctx context.Context, report entity.Report) error
	FindByID(ctx context.Context, id string) (dto.ReportResponse, error)
	Destroy(ctx context.Context, id string) error
	FindByReportScheduleID(ctx context.Context, reportScheduleID string) ([]dto.ReportResponse, error)
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{
		reportRepo: reportRepo,
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
func (s *reportService) Create(ctx context.Context, report entity.Report) (dto.ReportResponse, error) {
	// Generate new UUID for the report
	report.ID = uuid.New()

	// Set timestamps
	now := time.Now()
	report.CreatedAt = &now
	report.UpdatedAt = &now

	report, err := s.reportRepo.Create(ctx, report, nil)
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

// Update updates an existing report
func (s *reportService) Update(ctx context.Context, report entity.Report) error {
	// Update timestamp
	now := time.Now()
	report.UpdatedAt = &now

	err := s.reportRepo.Update(ctx, report, nil)
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
