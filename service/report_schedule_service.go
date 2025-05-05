package service

import (
	"context"
	"errors"
	"log"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"monitoring-service/helper"
	"monitoring-service/repository"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type reportScheduleService struct {
	reportScheduleRepo    repository.ReportScheduleReposiotry
	userManagementService *UserManagementService
	registrationService   *RegistrationManagementService
}

type ReportScheduleService interface {
	Index(ctx context.Context, token string) (dto.ReportScheduleByAdvisorResponse, error)
	Create(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (dto.ReportScheduleResponse, error)
	Update(ctx context.Context, id string, subject dto.ReportScheduleRequest, token string) error
	FindByID(ctx context.Context, id string, token string) (dto.ReportScheduleResponse, error)
	Destroy(ctx context.Context, id string, token string) error
	FindByRegistrationID(ctx context.Context, registrationID string) ([]dto.ReportScheduleResponse, error)
	ReportScheduleAccess(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (bool, error)
	FindByUserID(ctx context.Context, token string) ([]dto.ReportScheduleResponse, error)
	FindByAdvisorEmail(ctx context.Context, token string, pagReq dto.PaginationRequest, reportScheduleRequest dto.ReportScheduleAdvisorRequest) (dto.ReportScheduleByAdvisorResponse, dto.PaginationResponse, error)
	FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.ReportScheduleByStudentResponse, error)
}

func NewReportScheduleService(reportScheduleRepo repository.ReportScheduleReposiotry, userManagementbaseURI string, registrationManagementbaseURI string, asyncURIs []string) ReportScheduleService {
	return &reportScheduleService{
		reportScheduleRepo:    reportScheduleRepo,
		userManagementService: NewUserManagementService(userManagementbaseURI, asyncURIs),
		registrationService:   NewRegistrationManagementService(registrationManagementbaseURI, asyncURIs),
	}
}

func (s *reportScheduleService) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.ReportScheduleByStudentResponse, error) {
	user := s.userManagementService.GetUserData("GET", token)

	userNRP, ok := user["nrp"].(string)
	if !ok {
		return dto.ReportScheduleByStudentResponse{}, errors.New("user NRP not found")
	}

	reportSchedules, err := s.reportScheduleRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)
	if err != nil {
		log.Println("ERROR FINDING REPORT SCHEDULE BY USER NRP AND GROUP BY REGISTRATION ID: ", err)
		return dto.ReportScheduleByStudentResponse{}, err
	}

	reportScheduleResponses := make(map[string][]dto.ReportScheduleResponse)

	for registrationID, reportSchedules := range reportSchedules {
		registration := s.registrationService.GetRegistrationByID("GET", registrationID, token)
		registrationActivityName, ok := registration["activity_name"].(string)
		if !ok {
			return dto.ReportScheduleByStudentResponse{}, errors.New("registration activity name not found")
		}
		var reportScheduleResponse []dto.ReportScheduleResponse
		for _, reportSchedule := range reportSchedules {

			response := dto.ReportScheduleResponse{
				ID:                   reportSchedule.ID.String(),
				UserID:               reportSchedule.UserID,
				UserNRP:              reportSchedule.UserNRP,
				ActivityName:         registrationActivityName,
				RegistrationID:       reportSchedule.RegistrationID,
				AcademicAdvisorID:    reportSchedule.AcademicAdvisorID,
				AcademicAdvisorEmail: reportSchedule.AcademicAdvisorEmail,
				ReportType:           reportSchedule.ReportType,
				Week:                 reportSchedule.Week,
				StartDate:            reportSchedule.StartDate.Format(time.RFC3339),
				EndDate:              reportSchedule.EndDate.Format(time.RFC3339),
			}

			if len(reportSchedule.Report) > 0 {
				response.Report = &dto.ReportResponse{
					ID:                    reportSchedule.Report[0].ID.String(),
					ReportScheduleID:      reportSchedule.ID.String(),
					Title:                 reportSchedule.Report[0].Title,
					Content:               reportSchedule.Report[0].Content,
					ReportType:            reportSchedule.Report[0].ReportType,
					Feedback:              reportSchedule.Report[0].Feedback,
					AcademicAdvisorStatus: reportSchedule.Report[0].AcademicAdvisorStatus,
					FileStorageID:         reportSchedule.Report[0].FileStorageID,
				}
			}

			reportScheduleResponse = append(reportScheduleResponse, response)
		}

		reportScheduleResponses[registrationActivityName] = reportScheduleResponse
	}

	return dto.ReportScheduleByStudentResponse{
		Reports: reportScheduleResponses,
	}, nil
}

func (s *reportScheduleService) FindByAdvisorEmail(ctx context.Context, token string, pagReq dto.PaginationRequest, reportScheduleRequest dto.ReportScheduleAdvisorRequest) (dto.ReportScheduleByAdvisorResponse, dto.PaginationResponse, error) {
	user := s.userManagementService.GetUserData("GET", token)
	advisorEmail, ok := user["email"].(string)
	if !ok {
		return dto.ReportScheduleByAdvisorResponse{}, dto.PaginationResponse{}, errors.New("advisor email not found")
	}

	// log.Println("Getting report schedules for advisor email:", advisorEmail)

	reportSchedules, totalCount, err := s.reportScheduleRepo.FindByAdvisorEmailAndGroupByUserID(ctx, advisorEmail, nil, &pagReq, reportScheduleRequest.UserNRP)
	if err != nil {
		return dto.ReportScheduleByAdvisorResponse{}, dto.PaginationResponse{}, err
	}

	// log.Println("Found report schedules, total count:", totalCount)
	// log.Println("Map keys in reportSchedules:", getMapKeys(reportSchedules))

	reportScheduleResponses := make(map[string][]dto.ReportScheduleResponse)

	for userNRP, reportScheduleAdvisors := range reportSchedules {
		// log.Printf("Processing user NRP: %s with %d schedules", userNRP, len(reportScheduleAdvisors))

		var reportSchedule []dto.ReportScheduleResponse
		for _, reportScheduleAdvisor := range reportScheduleAdvisors {
			// log.Printf("  Processing schedule %d for NRP %s with actual UserNRP %s",
			// 	i, userNRP, reportScheduleAdvisor.UserNRP)

			// Create the base response
			// get activity name
			registration := s.registrationService.GetRegistrationByID("GET", reportScheduleAdvisor.RegistrationID, token)
			activityName, ok := registration["activity_name"].(string)
			if !ok {
				return dto.ReportScheduleByAdvisorResponse{}, dto.PaginationResponse{}, errors.New("activity name not found")
			}

			response := dto.ReportScheduleResponse{
				ID:                   reportScheduleAdvisor.ID.String(),
				UserID:               reportScheduleAdvisor.UserID,
				UserNRP:              reportScheduleAdvisor.UserNRP,
				ActivityName:         activityName,
				RegistrationID:       reportScheduleAdvisor.RegistrationID,
				AcademicAdvisorID:    reportScheduleAdvisor.AcademicAdvisorID,
				AcademicAdvisorEmail: reportScheduleAdvisor.AcademicAdvisorEmail,
				ReportType:           reportScheduleAdvisor.ReportType,
				Week:                 reportScheduleAdvisor.Week,
				StartDate:            reportScheduleAdvisor.StartDate.Format(time.RFC3339),
				EndDate:              reportScheduleAdvisor.EndDate.Format(time.RFC3339),
			}

			// Only set Report if there are reports
			if len(reportScheduleAdvisor.Report) > 0 {
				response.Report = &dto.ReportResponse{
					ID:                    reportScheduleAdvisor.Report[0].ID.String(),
					ReportScheduleID:      reportScheduleAdvisor.ID.String(),
					Title:                 reportScheduleAdvisor.Report[0].Title,
					Content:               reportScheduleAdvisor.Report[0].Content,
					ReportType:            reportScheduleAdvisor.Report[0].ReportType,
					Feedback:              reportScheduleAdvisor.Report[0].Feedback,
					AcademicAdvisorStatus: reportScheduleAdvisor.Report[0].AcademicAdvisorStatus,
					FileStorageID:         reportScheduleAdvisor.Report[0].FileStorageID,
				}
			}

			reportSchedule = append(reportSchedule, response)
		}

		// log.Printf("Adding %d schedules to map with key %s", len(reportSchedule), userNRP)
		reportScheduleResponses[userNRP] = reportSchedule
	}

	// log.Println("Final map keys in reportScheduleResponses:", getMapKeys(reportScheduleResponses))

	// Generate pagination metadata using the helper
	paginationResponse := helper.MetaDataPagination(totalCount, pagReq)

	return dto.ReportScheduleByAdvisorResponse{
		Reports: reportScheduleResponses,
	}, paginationResponse, nil
}

// Helper function to get map keys
func getMapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (s *reportScheduleService) FindByUserID(ctx context.Context, token string) ([]dto.ReportScheduleResponse, error) {
	user := s.userManagementService.GetUserData("GET", token)

	userNRP, ok := user["nrp"].(string)
	if !ok {
		return nil, errors.New("user NRP not found")
	}

	reportSchedules, err := s.reportScheduleRepo.FindByUserID(ctx, userNRP, nil)
	if err != nil {
		return nil, err
	}

	var reportScheduleResponses []dto.ReportScheduleResponse
	for _, reportSchedule := range reportSchedules {
		reportScheduleResponse := dto.ReportScheduleResponse{
			ID:                   reportSchedule.ID.String(),
			UserID:               reportSchedule.UserID,
			UserNRP:              reportSchedule.UserNRP,
			RegistrationID:       reportSchedule.RegistrationID,
			AcademicAdvisorID:    reportSchedule.AcademicAdvisorID,
			AcademicAdvisorEmail: reportSchedule.AcademicAdvisorEmail,
			ReportType:           reportSchedule.ReportType,
			Week:                 reportSchedule.Week,
			StartDate:            reportSchedule.StartDate.Format(time.RFC3339),
			EndDate:              reportSchedule.EndDate.Format(time.RFC3339),
		}

		if len(reportSchedule.Report) > 0 {
			reportScheduleResponse.Report = &dto.ReportResponse{
				ID:                    reportSchedule.Report[0].ID.String(),
				ReportScheduleID:      reportSchedule.ID.String(),
				Title:                 reportSchedule.Report[0].Title,
				Content:               reportSchedule.Report[0].Content,
				ReportType:            reportSchedule.Report[0].ReportType,
				Feedback:              reportSchedule.Report[0].Feedback,
				AcademicAdvisorStatus: reportSchedule.Report[0].AcademicAdvisorStatus,
				FileStorageID:         reportSchedule.Report[0].FileStorageID,
			}
		}
		reportScheduleResponses = append(reportScheduleResponses, reportScheduleResponse)
	}

	return reportScheduleResponses, nil
}

// Index retrieves all report schedules
func (s *reportScheduleService) Index(ctx context.Context, token string) (dto.ReportScheduleByAdvisorResponse, error) {
	reportSchedules, err := s.reportScheduleRepo.Index(ctx, nil)
	if err != nil {
		return dto.ReportScheduleByAdvisorResponse{}, err
	}

	reportScheduleResponses := make(map[string][]dto.ReportScheduleResponse)

	for userID, reportScheduleAdvisors := range reportSchedules {
		// get user by user id
		user := s.userManagementService.GetUserByID("GET", token, userID)
		userName, ok := user["name"].(string)
		if !ok {
			return dto.ReportScheduleByAdvisorResponse{}, errors.New("user name not found")
		}

		var reportSchedule []dto.ReportScheduleResponse
		for _, reportScheduleAdvisor := range reportScheduleAdvisors {
			// Create the base response
			response := dto.ReportScheduleResponse{
				ID:                   reportScheduleAdvisor.ID.String(),
				UserID:               reportScheduleAdvisor.UserID,
				UserNRP:              reportScheduleAdvisor.UserNRP,
				RegistrationID:       reportScheduleAdvisor.RegistrationID,
				AcademicAdvisorID:    reportScheduleAdvisor.AcademicAdvisorID,
				AcademicAdvisorEmail: reportScheduleAdvisor.AcademicAdvisorEmail,
				ReportType:           reportScheduleAdvisor.ReportType,
				Week:                 reportScheduleAdvisor.Week,
				StartDate:            reportScheduleAdvisor.StartDate.Format(time.RFC3339),
				EndDate:              reportScheduleAdvisor.EndDate.Format(time.RFC3339),
			}

			// Only set Report if there are reports
			if len(reportScheduleAdvisor.Report) > 0 {
				response.Report = &dto.ReportResponse{
					ID:                    reportScheduleAdvisor.Report[0].ID.String(),
					ReportScheduleID:      reportScheduleAdvisor.ID.String(),
					Title:                 reportScheduleAdvisor.Report[0].Title,
					Content:               reportScheduleAdvisor.Report[0].Content,
					ReportType:            reportScheduleAdvisor.Report[0].ReportType,
					Feedback:              reportScheduleAdvisor.Report[0].Feedback,
					AcademicAdvisorStatus: reportScheduleAdvisor.Report[0].AcademicAdvisorStatus,
					FileStorageID:         reportScheduleAdvisor.Report[0].FileStorageID,
				}
			}

			reportSchedule = append(reportSchedule, response)
		}
		reportScheduleResponses[userName] = reportSchedule // Convert to string
	}

	return dto.ReportScheduleByAdvisorResponse{
		Reports: reportScheduleResponses,
	}, nil
}

func (s *reportScheduleService) ReportScheduleAccess(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (bool, error) {
	user := s.userManagementService.GetUserData("GET", token)

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
	} else if userRole != "ADMIN" && userRole != "LO-MBKM" {
		return false, errors.New("user role not allowed")
	}

	return true, nil
}

// Create creates a new report schedule
func (s *reportScheduleService) Create(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (dto.ReportScheduleResponse, error) {
	// first we need to get user email and user id
	access, err := s.ReportScheduleAccess(ctx, reportSchedule, token)
	if err != nil {
		return dto.ReportScheduleResponse{}, err
	}

	if !access {
		return dto.ReportScheduleResponse{}, errors.New("user role not allowed")
	}

	var reportScheduleEntity entity.ReportSchedule
	reportScheduleEntity.ID = uuid.New()
	reportScheduleEntity.UserID = reportSchedule.UserID
	reportScheduleEntity.UserNRP = reportSchedule.UserNRP
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

func (s *reportScheduleService) Update(ctx context.Context, id string, subject dto.ReportScheduleRequest, token string) error {

	access, err := s.ReportScheduleAccess(ctx, subject, token)
	if err != nil {
		return err
	}

	if !access {
		return errors.New("user role not allowed")
	}

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
func (s *reportScheduleService) FindByID(ctx context.Context, id string, token string) (dto.ReportScheduleResponse, error) {

	reportSchedule, err := s.reportScheduleRepo.FindByID(ctx, id, nil)
	if err != nil {
		return dto.ReportScheduleResponse{}, err
	}

	reportScheduleRequest := dto.ReportScheduleRequest{
		AcademicAdvisorEmail: reportSchedule.AcademicAdvisorEmail,
	}

	access, err := s.ReportScheduleAccess(ctx, reportScheduleRequest, token)
	if err != nil {
		return dto.ReportScheduleResponse{}, err
	}

	if !access {
		return dto.ReportScheduleResponse{}, errors.New("user role not allowed")
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

	if len(reportSchedule.Report) > 0 {
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
	}

	return reportScheduleResponse, nil
}

// Destroy deletes a report schedule
func (s *reportScheduleService) Destroy(ctx context.Context, id string, token string) error {

	access, err := s.ReportScheduleAccess(ctx, dto.ReportScheduleRequest{}, token)
	if err != nil {
		return err
	}

	if !access {
		return errors.New("user role not allowed")
	}

	err = s.reportScheduleRepo.Destroy(ctx, id, nil)
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
