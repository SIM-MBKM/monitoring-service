package service

import (
	"context"
	"errors"
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
	reportRepo            repository.ReportRepository
	reportScheduleRepo    repository.ReportScheduleReposiotry
	fileService           *FileService
	userManagementService *UserManagementService
}

type ReportService interface {
	Index(ctx context.Context) ([]dto.ReportResponse, error)
	Create(ctx context.Context, report dto.ReportRequest, file *multipart.FileHeader, token string) (dto.ReportResponse, error)
	Update(ctx context.Context, id string, report dto.ReportRequest) error
	FindByID(ctx context.Context, id string, token string) (dto.ReportResponse, error)
	Destroy(ctx context.Context, id string) error
	FindByReportScheduleID(ctx context.Context, reportScheduleID string) ([]dto.ReportResponse, error)
	Approval(ctx context.Context, id string, token string, report dto.ReportApprovalRequest) error
}

func NewReportService(reportRepo repository.ReportRepository, reportScheduleRepo repository.ReportScheduleReposiotry, userManagementBaseURI string, asyncURIs []string, config *storageService.Config, tokenManager *storageService.CacheTokenManager) ReportService {
	return &reportService{
		reportRepo:            reportRepo,
		reportScheduleRepo:    reportScheduleRepo,
		fileService:           NewFileService(config, tokenManager),
		userManagementService: NewUserManagementService(userManagementBaseURI, asyncURIs),
	}
}

func (s *reportService) Approval(ctx context.Context, id string, token string, report dto.ReportApprovalRequest) error {

	reportEntity, err := s.reportRepo.FindByID(ctx, id, nil)
	if err != nil {
		return err
	}

	reportSchedule, err := s.reportScheduleRepo.FindByID(ctx, reportEntity.ReportScheduleID, nil)
	if err != nil {
		return err
	}

	advisor := s.userManagementService.GetUserData("GET", token)
	advisorEmail, ok := advisor["email"].(string)
	if !ok {
		return errors.New("advisor email not found")
	}

	if advisorEmail != reportSchedule.AcademicAdvisorEmail {
		return errors.New("unauthorized")
	}
	reportEntity.AcademicAdvisorStatus = report.Status
	reportEntity.Feedback = report.Feedback

	err = s.reportRepo.Approval(ctx, id, reportEntity, nil)
	if err != nil {
		return err
	}

	return nil
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
func (s *reportService) Create(ctx context.Context, report dto.ReportRequest, file *multipart.FileHeader, token string) (dto.ReportResponse, error) {
	// Generate new UUID for the report
	var result *storageService.FileResponse
	if file != nil {
		var err error
		result, err = s.fileService.storage.GcsUpload(file, "sim_mbkm", "", "")
		if err != nil {
			return dto.ReportResponse{}, err
		}
	}

	reportSchedule, err := s.reportScheduleRepo.FindByID(ctx, report.ReportScheduleID, nil)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	user := s.userManagementService.GetUserData("GET", token)

	if user["id"] != reportSchedule.UserID {
		return dto.ReportResponse{}, errors.New("unauthorized")
	}

	var reportEntity entity.Report
	reportEntity.ID = uuid.New()
	reportEntity.ReportScheduleID = report.ReportScheduleID
	if result != nil {
		reportEntity.FileStorageID = result.FileID
	}
	reportEntity.Title = report.Title
	reportEntity.Content = report.Content
	reportEntity.ReportType = report.ReportType
	reportEntity.AcademicAdvisorStatus = "PENDING"

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

func (s *reportService) ReportAccess(ctx context.Context, report dto.ReportRequest, token string) (bool, error) {
	user := s.userManagementService.GetUserData("GET", token)

	// get report schedule
	reportSchedule, err := s.reportScheduleRepo.FindByID(ctx, report.ReportScheduleID, nil)
	if err != nil {
		return false, err
	}

	// convert userRole to string
	userRole, ok := user["role"].(string)
	if !ok {
		return false, errors.New("user role not found")
	}

	if userRole == "DOSEN PEMBIMBING" {
		userEmail, ok := user["email"].(string)
		if !ok {
			return false, errors.New("user email not found")
		}

		if userEmail != reportSchedule.AcademicAdvisorEmail {
			return false, errors.New("user email not match")
		}
	} else if userRole == "MAHASISWA" {
		userId, ok := user["id"].(string)
		if !ok {
			return false, errors.New("user id not found")
		}

		if userId != reportSchedule.UserID {
			return false, errors.New("user id not match")
		}
	} else if userRole != "ADMIN" && userRole != "LO-MBKM" {
		return false, errors.New("user role not allowed")
	}

	return true, nil
}

// FindByID retrieves a report by its ID
func (s *reportService) FindByID(ctx context.Context, id string, token string) (dto.ReportResponse, error) {

	report, err := s.reportRepo.FindByID(ctx, id, nil)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	access, err := s.ReportAccess(ctx, dto.ReportRequest{ReportScheduleID: report.ReportScheduleID}, token)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	if !access {
		return dto.ReportResponse{}, errors.New("unauthorized")
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
