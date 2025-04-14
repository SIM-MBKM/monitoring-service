package service

import (
	"context"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"monitoring-service/repository"
	"time"

	"github.com/google/uuid"
)

type reportScheduleService struct {
	reportScheduleRepo repository.ReportScheduleReposiotry
}

type ReportScheduleService interface {
	Index(ctx context.Context) ([]dto.ReportScheduleResponse, error)
	Create(ctx context.Context, reportSchedule entity.ReportSchedule) (dto.ReportScheduleResponse, error)
	Update(ctx context.Context, reportSchedule entity.ReportSchedule) error
	FindByID(ctx context.Context, id string) (dto.ReportScheduleResponse, error)
	Destroy(ctx context.Context, id string) error
	FindByRegistrationID(ctx context.Context, registrationID string) ([]dto.ReportScheduleResponse, error)
}

func NewReportScheduleService(reportScheduleRepo repository.ReportScheduleReposiotry) ReportScheduleService {
	return &reportScheduleService{
		reportScheduleRepo: reportScheduleRepo,
	}
}

// Index retrieves all report schedules
func (s *reportScheduleService) Index(ctx context.Context) ([]dto.ReportScheduleResponse, error) {
	reportSchedules, err := s.reportScheduleRepo.Index(ctx, nil)
	if err != nil {
		return nil, err
	}

	var reportScheduleResponses []dto.ReportScheduleResponse
	for _, reportSchedule := range reportSchedules {
		reportScheduleResponses = append(reportScheduleResponses, dto.ReportScheduleResponse{
			ID:                reportSchedule.ID.String(),
			RegistrationID:    reportSchedule.RegistrationID,
			AcademicAdvisorID: reportSchedule.AcademicAdvisorID,
			ReportType:        reportSchedule.ReportType,
			Week:              reportSchedule.Week,
			StartDate:         reportSchedule.StartDate.Format("2006-01-02"),
			EndDate:           reportSchedule.EndDate.Format("2006-01-02"),
			Report: dto.ReportResponse{
				ID:                    reportSchedule.Report[0].ID.String(),
				ReportScheduleID:      reportSchedule.ID.String(),
				Title:                 reportSchedule.Report[0].Title,
				Content:               reportSchedule.Report[0].Content,
				ReportType:            reportSchedule.Report[0].ReportType,
				Feedback:              reportSchedule.Report[0].Feedback,
				AcademicAdvisorStatus: reportSchedule.Report[0].AcademicAdvisorStatus,
				FileStorageID:         reportSchedule.Report[0].FileStorageID,
			},
		})
	}

	return reportScheduleResponses, nil
}

// Create creates a new report schedule
func (s *reportScheduleService) Create(ctx context.Context, reportSchedule entity.ReportSchedule) (dto.ReportScheduleResponse, error) {
	// Generate new UUID for the report schedule
	reportSchedule.ID = uuid.New()

	// Set timestamps
	now := time.Now()
	reportSchedule.CreatedAt = &now
	reportSchedule.UpdatedAt = &now

	reportSchedule, err := s.reportScheduleRepo.Create(ctx, reportSchedule, nil)
	if err != nil {
		return dto.ReportScheduleResponse{}, err
	}

	var reportScheduleResponse dto.ReportScheduleResponse
	reportScheduleResponse.ID = reportSchedule.ID.String()
	reportScheduleResponse.RegistrationID = reportSchedule.RegistrationID
	reportScheduleResponse.AcademicAdvisorID = reportSchedule.AcademicAdvisorID
	reportScheduleResponse.ReportType = reportSchedule.ReportType
	reportScheduleResponse.Week = reportSchedule.Week
	reportScheduleResponse.StartDate = reportSchedule.StartDate.Format("2006-01-02")
	reportScheduleResponse.EndDate = reportSchedule.EndDate.Format("2006-01-02")

	return reportScheduleResponse, nil
}

// Update updates an existing report schedule
func (s *reportScheduleService) Update(ctx context.Context, reportSchedule entity.ReportSchedule) error {
	// Update timestamp
	now := time.Now()
	reportSchedule.UpdatedAt = &now

	err := s.reportScheduleRepo.Update(ctx, reportSchedule, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindByID retrieves a report schedule by its ID
func (s *reportScheduleService) FindByID(ctx context.Context, id string) (dto.ReportScheduleResponse, error) {
	reportSchedule, err := s.reportScheduleRepo.FindByID(ctx, id, nil)
	if err != nil {
		return dto.ReportScheduleResponse{}, err
	}

	var reportScheduleResponse dto.ReportScheduleResponse
	reportScheduleResponse.ID = reportSchedule.ID.String()
	reportScheduleResponse.RegistrationID = reportSchedule.RegistrationID
	reportScheduleResponse.AcademicAdvisorID = reportSchedule.AcademicAdvisorID
	reportScheduleResponse.ReportType = reportSchedule.ReportType
	reportScheduleResponse.Week = reportSchedule.Week
	reportScheduleResponse.StartDate = reportSchedule.StartDate.Format("2006-01-02")
	reportScheduleResponse.EndDate = reportSchedule.EndDate.Format("2006-01-02")

	reportScheduleResponse.Report = dto.ReportResponse{
		ID:                    reportSchedule.Report[0].ID.String(),
		ReportScheduleID:      reportSchedule.ID.String(),
		FileStorageID:         reportSchedule.Report[0].FileStorageID,
		Title:                 reportSchedule.Report[0].Title,
		Content:               reportSchedule.Report[0].Content,
		ReportType:            reportSchedule.Report[0].ReportType,
		Feedback:              reportSchedule.Report[0].Feedback,
		AcademicAdvisorStatus: reportSchedule.Report[0].AcademicAdvisorStatus,
	}

	return reportScheduleResponse, nil
}

// Destroy deletes a report schedule
func (s *reportScheduleService) Destroy(ctx context.Context, id string) error {
	err := s.reportScheduleRepo.Destroy(ctx, id, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindByRegistrationID retrieves report schedules by registration ID
func (s *reportScheduleService) FindByRegistrationID(ctx context.Context, registrationID string) ([]dto.ReportScheduleResponse, error) {
	reportSchedules, err := s.reportScheduleRepo.FindByRegistrationID(ctx, registrationID, nil)
	if err != nil {
		return nil, err
	}

	var reportScheduleResponses []dto.ReportScheduleResponse
	for _, reportSchedule := range reportSchedules {
		reportScheduleResponses = append(reportScheduleResponses, dto.ReportScheduleResponse{
			ID:                reportSchedule.ID.String(),
			RegistrationID:    reportSchedule.RegistrationID,
			AcademicAdvisorID: reportSchedule.AcademicAdvisorID,
			ReportType:        reportSchedule.ReportType,
			Week:              reportSchedule.Week,
			StartDate:         reportSchedule.StartDate.Format("2006-01-02"),
			EndDate:           reportSchedule.EndDate.Format("2006-01-02"),
			Report: dto.ReportResponse{
				ID:                    reportSchedule.Report[0].ID.String(),
				ReportScheduleID:      reportSchedule.ID.String(),
				FileStorageID:         reportSchedule.Report[0].FileStorageID,
				Title:                 reportSchedule.Report[0].Title,
				Content:               reportSchedule.Report[0].Content,
				ReportType:            reportSchedule.Report[0].ReportType,
				Feedback:              reportSchedule.Report[0].Feedback,
				AcademicAdvisorStatus: reportSchedule.Report[0].AcademicAdvisorStatus,
			},
		})
	}

	return reportScheduleResponses, nil
}
