package service

import (
	"context"
	"errors"
	"log"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"monitoring-service/repository"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type reportScheduleService struct {
	reportScheduleRepo    repository.ReportScheduleReposiotry
	userManagementService *UserManagementService
}

type ReportScheduleService interface {
	Index(ctx context.Context) ([]dto.ReportScheduleResponse, error)
	Create(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (dto.ReportScheduleResponse, error)
	Update(ctx context.Context, id string, subject dto.ReportScheduleRequest) error
	FindByID(ctx context.Context, id string) (dto.ReportScheduleResponse, error)
	Destroy(ctx context.Context, id string) error
	FindByRegistrationID(ctx context.Context, registrationID string) ([]dto.ReportScheduleResponse, error)
}

func NewReportScheduleService(reportScheduleRepo repository.ReportScheduleReposiotry, userManagementbaseURI string, asyncURIs []string) ReportScheduleService {
	return &reportScheduleService{
		reportScheduleRepo:    reportScheduleRepo,
		userManagementService: NewUserManagementService(userManagementbaseURI, asyncURIs),
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
			ID:                   reportSchedule.ID.String(),
			UserID:               reportSchedule.UserID,
			RegistrationID:       reportSchedule.RegistrationID,
			AcademicAdvisorID:    reportSchedule.AcademicAdvisorID,
			AcademicAdvisorEmail: reportSchedule.AcademicAdvisorEmail,
			ReportType:           reportSchedule.ReportType,
			Week:                 reportSchedule.Week,
			StartDate:            reportSchedule.StartDate.Format(time.RFC3339),
			EndDate:              reportSchedule.EndDate.Format(time.RFC3339),
			Report: &dto.ReportResponse{
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
func (s *reportScheduleService) Create(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (dto.ReportScheduleResponse, error) {
	// first we need to get user email and user id
	user := s.userManagementService.GetUserData("GET", token)

	// convert userRole to string
	userRole, ok := user["role"].(string)
	if !ok {
		return dto.ReportScheduleResponse{}, errors.New("user role not found")
	}

	if userRole == "DOSEN PEMBIMBING" {
		userEmail, ok := user["email"].(string)
		if !ok {
			return dto.ReportScheduleResponse{}, errors.New("user email not found")
		}

		if userEmail != reportSchedule.AcademicAdvisorEmail {
			return dto.ReportScheduleResponse{}, errors.New("user email not match")
		}
	} else if userRole != "ADMIN" && userRole != "LO-MBKM" {
		return dto.ReportScheduleResponse{}, errors.New("user role not allowed")
	}

	var reportScheduleEntity entity.ReportSchedule
	reportScheduleEntity.ID = uuid.New()
	reportScheduleEntity.UserID = reportSchedule.UserID
	reportScheduleEntity.RegistrationID = reportSchedule.RegistrationID
	reportScheduleEntity.AcademicAdvisorID = reportSchedule.AcademicAdvisorID
	reportScheduleEntity.AcademicAdvisorEmail = reportSchedule.AcademicAdvisorEmail
	reportScheduleEntity.ReportType = reportSchedule.ReportType
	reportScheduleEntity.Week = reportSchedule.Week
	// convert string to time.Time
	start_date, err := time.Parse(time.RFC3339, reportSchedule.StartDate)
	if err != nil {
		log.Println("ERROR CONVERTING START DATE: ", err)
		return dto.ReportScheduleResponse{}, err
	}
	reportScheduleEntity.StartDate = &start_date
	end_date, err := time.Parse(time.RFC3339, reportSchedule.EndDate)
	if err != nil {
		log.Println("ERROR CONVERTING END DATE: ", err)
		return dto.ReportScheduleResponse{}, err
	}
	reportScheduleEntity.EndDate = &end_date

	// Set timestamps
	now := time.Now()
	reportScheduleEntity.CreatedAt = &now
	reportScheduleEntity.UpdatedAt = &now

	_, err = s.reportScheduleRepo.Create(ctx, reportScheduleEntity, nil)
	if err != nil {
		log.Println("ERROR CREATING REPORT SCHEDULE: ", err)
		return dto.ReportScheduleResponse{}, err
	}

	return dto.ReportScheduleResponse{
		ID:                   reportScheduleEntity.ID.String(),
		UserID:               reportScheduleEntity.UserID,
		RegistrationID:       reportScheduleEntity.RegistrationID,
		AcademicAdvisorID:    reportScheduleEntity.AcademicAdvisorID,
		AcademicAdvisorEmail: reportScheduleEntity.AcademicAdvisorEmail,
		ReportType:           reportScheduleEntity.ReportType,
		Week:                 reportScheduleEntity.Week,
		StartDate:            reportScheduleEntity.StartDate.Format(time.RFC3339),
		EndDate:              reportScheduleEntity.EndDate.Format(time.RFC3339),
	}, nil
}

func (s *reportScheduleService) Update(ctx context.Context, id string, subject dto.ReportScheduleRequest) error {
	res, err := s.reportScheduleRepo.FindByID(ctx, id, nil)
	if err != nil {
		return err
	}

	// Create programTypeEntity with original ID
	reportScheduleEntity := entity.ReportSchedule{
		ID: res.ID,
	}

	// Get reflection values
	resValue := reflect.ValueOf(res)
	reqValue := reflect.ValueOf(subject)
	entityValue := reflect.ValueOf(&reportScheduleEntity).Elem()

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
	err = s.reportScheduleRepo.Update(ctx, id, reportScheduleEntity, nil)
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
	reportScheduleResponse.UserID = reportSchedule.UserID
	reportScheduleResponse.RegistrationID = reportSchedule.RegistrationID
	reportScheduleResponse.AcademicAdvisorID = reportSchedule.AcademicAdvisorID
	reportScheduleResponse.AcademicAdvisorEmail = reportSchedule.AcademicAdvisorEmail
	reportScheduleResponse.ReportType = reportSchedule.ReportType
	reportScheduleResponse.Week = reportSchedule.Week
	reportScheduleResponse.StartDate = reportSchedule.StartDate.Format(time.RFC3339)
	reportScheduleResponse.EndDate = reportSchedule.EndDate.Format(time.RFC3339)

	reportScheduleResponse.Report = &dto.ReportResponse{
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
			ID:                   reportSchedule.ID.String(),
			UserID:               reportSchedule.UserID,
			RegistrationID:       reportSchedule.RegistrationID,
			AcademicAdvisorID:    reportSchedule.AcademicAdvisorID,
			AcademicAdvisorEmail: reportSchedule.AcademicAdvisorEmail,
			ReportType:           reportSchedule.ReportType,
			Week:                 reportSchedule.Week,
			StartDate:            reportSchedule.StartDate.Format(time.RFC3339),
			EndDate:              reportSchedule.EndDate.Format(time.RFC3339),
			Report: &dto.ReportResponse{
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
