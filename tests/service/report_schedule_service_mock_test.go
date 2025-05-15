package service_test

import (
	"context"
	"errors"
	"log"
	"monitoring-service/dto"
	"monitoring-service/entity"
	repository_mock "monitoring-service/mocks/repository"
	service_mock "monitoring-service/mocks/service"
	"monitoring-service/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ReportScheduleServiceTestSuite struct {
	suite.Suite
	mockReportScheduleRepo    *repository_mock.MockReportScheduleRepository
	mockUserManagementService *service_mock.MockUserManagementService
	mockRegistrationService   *service_mock.MockRegistrationService
	service                   service.ReportScheduleService
}

func (suite *ReportScheduleServiceTestSuite) SetupTest() {
	suite.mockReportScheduleRepo = new(repository_mock.MockReportScheduleRepository)
	suite.mockUserManagementService = new(service_mock.MockUserManagementService)
	suite.mockRegistrationService = new(service_mock.MockRegistrationService)

	createService := func(repo *repository_mock.MockReportScheduleRepository,
		userManagementService *service_mock.MockUserManagementService,
		registrationService *service_mock.MockRegistrationService,
	) service.ReportScheduleService {

		svc := &mockReportScheduleService{
			reportScheduleRepo:    repo,
			userManagementService: userManagementService,
			registrationService:   registrationService,
		}

		return svc
	}

	suite.service = createService(
		suite.mockReportScheduleRepo,
		suite.mockUserManagementService,
		suite.mockRegistrationService,
	)
}

type mockReportScheduleService struct {
	reportScheduleRepo    *repository_mock.MockReportScheduleRepository
	userManagementService *service_mock.MockUserManagementService
	registrationService   *service_mock.MockRegistrationService
}

// Test FindByUserNRPAndGroupByRegistrationID - Success
func (suite *ReportScheduleServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"
	reportSchedules := map[string][]entity.ReportSchedule{
		registrationID: []entity.ReportSchedule{
			{
				ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
				UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
				UserNRP:              userNRP,
				RegistrationID:       registrationID,
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: "test@gmail.com",
				ReportType:           "FINAL_REPORT",
				Week:                 1,
				StartDate:            func() *time.Time { t := time.Now(); return &t }(),
				EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
				BaseModel: entity.BaseModel{
					CreatedAt: func() *time.Time { t := time.Now(); return &t }(),
					UpdatedAt: func() *time.Time { t := time.Now(); return &t }(),
				},
				Report: []entity.Report{
					{
						ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
						ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
						Title:                 "Test Report",
						Content:               "Test Content",
						ReportType:            "FINAL_REPORT",
						Feedback:              "Great work",
						AcademicAdvisorStatus: "APPROVED",
						FileStorageID:         "file-123",
					},
				},
			},
		},
	}

	registration := map[string]interface{}{
		"id":                     registrationID,
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_name":              "Dimas Fadilah",
		"academic_advisor":       "test",
		"academic_advisor_email": "test@gmail.com",
		"activity_name":          "Test Activity",
		"approval_status":        true,
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(reportSchedules, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", registrationID, token).Return(registration)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Contains(suite.T(), result.Reports, "Test Activity")

	reports := result.Reports["Test Activity"]
	assert.Equal(suite.T(), 1, len(reports))
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b3", reports[0].ID)
	assert.Equal(suite.T(), "Test Activity", reports[0].ActivityName)
	assert.NotNil(suite.T(), reports[0].Report)
	assert.Equal(suite.T(), "Test Report", reports[0].Report.Title)
}

// Test FindByUserNRPAndGroupByRegistrationID - User NRP Not Found
func (suite *ReportScheduleServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_UserNRPNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	usersData := map[string]interface{}{
		"id": "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		// NRP is missing
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user NRP not found", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleByStudentResponse{}, result)
}

// Test FindByUserNRPAndGroupByRegistrationID - Repository Error
func (suite *ReportScheduleServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(map[string][]entity.ReportSchedule{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleByStudentResponse{}, result)
}

// Test FindByUserNRPAndGroupByRegistrationID - Activity Name Not Found
func (suite *ReportScheduleServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_ActivityNameNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"
	reportSchedules := map[string][]entity.ReportSchedule{
		registrationID: []entity.ReportSchedule{
			{
				ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
				UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
				UserNRP:              userNRP,
				RegistrationID:       registrationID,
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: "test@gmail.com",
				ReportType:           "FINAL_REPORT",
				Week:                 1,
				StartDate:            func() *time.Time { t := time.Now(); return &t }(),
				EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
			},
		},
	}

	registration := map[string]interface{}{
		"id":                     registrationID,
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_name":              "Dimas Fadilah",
		"academic_advisor":       "test",
		"academic_advisor_email": "test@gmail.com",
		// Activity name is missing
		"approval_status": true,
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(reportSchedules, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", registrationID, token).Return(registration)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "registration activity name not found", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleByStudentResponse{}, result)
}

// Test FindByUserNRPAndGroupByRegistrationID - Approval Status Not Found
func (suite *ReportScheduleServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_ApprovalStatusNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"
	reportSchedules := map[string][]entity.ReportSchedule{
		registrationID: []entity.ReportSchedule{
			{
				ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
				UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
				UserNRP:              userNRP,
				RegistrationID:       registrationID,
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: "test@gmail.com",
				ReportType:           "FINAL_REPORT",
				Week:                 1,
				StartDate:            func() *time.Time { t := time.Now(); return &t }(),
				EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
			},
		},
	}

	registration := map[string]interface{}{
		"id":                     registrationID,
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_name":              "Dimas Fadilah",
		"academic_advisor":       "test",
		"academic_advisor_email": "test@gmail.com",
		"activity_name":          "Test Activity",
		// Approval status is missing
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(reportSchedules, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", registrationID, token).Return(registration)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "registration approval status not found", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleByStudentResponse{}, result)
}

// Test FindByUserNRPAndGroupByRegistrationID - No Approved Schedules
func (suite *ReportScheduleServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_NoApprovedSchedules() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"
	reportSchedules := map[string][]entity.ReportSchedule{
		registrationID: []entity.ReportSchedule{
			{
				ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
				UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
				UserNRP:              userNRP,
				RegistrationID:       registrationID,
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: "test@gmail.com",
				ReportType:           "FINAL_REPORT",
				Week:                 1,
				StartDate:            func() *time.Time { t := time.Now(); return &t }(),
				EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
			},
		},
	}

	registration := map[string]interface{}{
		"id":                     registrationID,
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_name":              "Dimas Fadilah",
		"academic_advisor":       "test",
		"academic_advisor_email": "test@gmail.com",
		"activity_name":          "Test Activity",
		"approval_status":        false, // Not approved
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(reportSchedules, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", registrationID, token).Return(registration)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result.Reports)
}

func (s *mockReportScheduleService) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.ReportScheduleByStudentResponse, error) {
	usersData := s.userManagementService.GetUserData("GET", token)

	userNRP, ok := usersData["nrp"].(string)
	if !ok {
		return dto.ReportScheduleByStudentResponse{}, errors.New("user NRP not found")
	}

	reportSchedules, err := s.reportScheduleRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)
	if err != nil {
		return dto.ReportScheduleByStudentResponse{}, err
	}

	reportScheduleResponses := make(map[string][]dto.ReportScheduleResponse)

	for registrationID, reportSchedules := range reportSchedules {
		registration := s.registrationService.GetRegistrationByID("GET", registrationID, token)
		registrationActivityName, ok := registration["activity_name"].(string)
		if !ok {
			return dto.ReportScheduleByStudentResponse{}, errors.New("registration activity name not found")
		}

		approvalStatus, ok := registration["approval_status"].(bool)
		if !ok {
			return dto.ReportScheduleByStudentResponse{}, errors.New("registration approval status not found")
		}

		if !approvalStatus {
			continue
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

// Test Create - Success
func (suite *ReportScheduleServiceTestSuite) TestCreate_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
		ReportType:           "FINAL_REPORT",
		Week:                 1,
		StartDate:            time.Now().Format(time.RFC3339),
		EndDate:              time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN", // Role with access
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("Create", ctx, mock.AnythingOfType("entity.ReportSchedule"), mock.Anything).Return(entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               reportScheduleRequest.UserID,
		UserNRP:              reportScheduleRequest.UserNRP,
		RegistrationID:       reportScheduleRequest.RegistrationID,
		AcademicAdvisorID:    reportScheduleRequest.AcademicAdvisorID,
		AcademicAdvisorEmail: reportScheduleRequest.AcademicAdvisorEmail,
		ReportType:           reportScheduleRequest.ReportType,
		Week:                 reportScheduleRequest.Week,
		StartDate:            func() *time.Time { t, _ := time.Parse(time.RFC3339, reportScheduleRequest.StartDate); return &t }(),
		EndDate:              func() *time.Time { t, _ := time.Parse(time.RFC3339, reportScheduleRequest.EndDate); return &t }(),
	}, nil)

	// Call the method
	result, err := suite.service.Create(ctx, reportScheduleRequest, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b3", result.ID)
	assert.Equal(suite.T(), reportScheduleRequest.UserID, result.UserID)
	assert.Equal(suite.T(), reportScheduleRequest.ReportType, result.ReportType)
}

// Test Create - Invalid User Role
func (suite *ReportScheduleServiceTestSuite) TestCreate_InvalidUserRole() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorEmail: "test@gmail.com",
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",           // Invalid role for creating schedules
		"email": "different@gmail.com", // Different email from academic advisor
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.Create(ctx, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user role not allowed", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleResponse{}, result)
}

// Test Create - Advisor Email Mismatch
func (suite *ReportScheduleServiceTestSuite) TestCreate_AdvisorEmailMismatch() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportSchedule := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@example.com",
		ReportType:           "WEEKLY_REPORT",
		Week:                 1,
		StartDate:            "2023-05-01T00:00:00Z",
		EndDate:              "2023-05-07T23:59:59Z",
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"nrp":   "10000000",
		"name":  "John Doe",
		"role":  "DOSEN PEMBIMBING",
		"email": "different@example.com", // Different from the request
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	_, err := suite.service.Create(ctx, reportSchedule, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user email not found", err.Error())
}

// Test Create - Repository Error
func (suite *ReportScheduleServiceTestSuite) TestCreate_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
		ReportType:           "FINAL_REPORT",
		Week:                 1,
		StartDate:            time.Now().Format(time.RFC3339),
		EndDate:              time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("Create", ctx, mock.AnythingOfType("entity.ReportSchedule"), mock.Anything).Return(entity.ReportSchedule{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.Create(ctx, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleResponse{}, result)
}

// Test Create - Invalid Date Format
func (suite *ReportScheduleServiceTestSuite) TestCreate_InvalidDateFormat() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
		ReportType:           "FINAL_REPORT",
		Week:                 1,
		StartDate:            "invalid-date-format", // Invalid date format
		EndDate:              time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.Create(ctx, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "parsing time")
	assert.Equal(suite.T(), dto.ReportScheduleResponse{}, result)
}

// Test FindByID - Success
func (suite *ReportScheduleServiceTestSuite) TestFindByID_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN",
		"email": "dimasfadilah20@gmail.com",
	}

	reportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
		ReportType:           "FINAL_REPORT",
		Week:                 1,
		StartDate:            func() *time.Time { t := time.Now(); return &t }(),
		EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
		Report: []entity.Report{
			{
				ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
				ReportScheduleID:      id,
				Title:                 "Test Report",
				Content:               "Test Content",
				ReportType:            "FINAL_REPORT",
				Feedback:              "Great work",
				AcademicAdvisorStatus: "APPROVED",
				FileStorageID:         "file-123",
			},
		},
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByID", ctx, id, mock.Anything).Return(reportSchedule, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), id, result.ID)
	assert.Equal(suite.T(), reportSchedule.UserID, result.UserID)
	assert.Equal(suite.T(), reportSchedule.ReportType, result.ReportType)
	assert.NotNil(suite.T(), result.Report)
	assert.Equal(suite.T(), "Test Report", result.Report.Title)
}

// Test FindByID - Not Found
func (suite *ReportScheduleServiceTestSuite) TestFindByID_NotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.ReportSchedule{}, errors.New("record not found"))

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleResponse{}, result)
}

// Test FindByID - Access Denied
func (suite *ReportScheduleServiceTestSuite) TestFindByID_AccessDenied() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA", // Invalid role
		"email": "dimasfadilah20@gmail.com",
	}

	reportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
		ReportType:           "FINAL_REPORT",
		Week:                 1,
		StartDate:            func() *time.Time { t := time.Now(); return &t }(),
		EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByID", ctx, id, mock.Anything).Return(reportSchedule, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user role not allowed", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleResponse{}, result)
}

// Test Update - Success
func (suite *ReportScheduleServiceTestSuite) TestUpdate_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
		ReportType:           "UPDATED_TYPE",
		Week:                 2,
		StartDate:            time.Now().Format(time.RFC3339),
		EndDate:              time.Now().AddDate(0, 0, 14).Format(time.RFC3339),
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN",
		"email": "dimasfadilah20@gmail.com",
	}

	originalReportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
		ReportType:           "FINAL_REPORT",
		Week:                 1,
		StartDate:            func() *time.Time { t := time.Now(); return &t }(),
		EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByID", ctx, id, mock.Anything).Return(originalReportSchedule, nil)
	suite.mockReportScheduleRepo.On("Update", ctx, id, mock.AnythingOfType("entity.ReportSchedule"), mock.Anything).Return(nil)

	// Call the method
	err := suite.service.Update(ctx, id, reportScheduleRequest, token)

	// Assertions
	assert.NoError(suite.T(), err)
}

// Test Update - Access Denied
func (suite *ReportScheduleServiceTestSuite) TestUpdate_AccessDenied() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorEmail: "test@gmail.com",
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA", // Invalid role
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	err := suite.service.Update(ctx, id, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user role not allowed", err.Error())
}

// Test Update - Record Not Found
func (suite *ReportScheduleServiceTestSuite) TestUpdate_RecordNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.ReportSchedule{}, errors.New("record not found"))

	// Call the method
	err := suite.service.Update(ctx, id, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
}

// Test Update - Repository Error
func (suite *ReportScheduleServiceTestSuite) TestUpdate_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
		ReportType:           "UPDATED_TYPE",
		Week:                 2,
		StartDate:            time.Now().Format(time.RFC3339),
		EndDate:              time.Now().AddDate(0, 0, 14).Format(time.RFC3339),
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN",
		"email": "dimasfadilah20@gmail.com",
	}

	originalReportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "test@gmail.com",
		ReportType:           "FINAL_REPORT",
		Week:                 1,
		StartDate:            func() *time.Time { t := time.Now(); return &t }(),
		EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByID", ctx, id, mock.Anything).Return(originalReportSchedule, nil)
	suite.mockReportScheduleRepo.On("Update", ctx, id, mock.AnythingOfType("entity.ReportSchedule"), mock.Anything).Return(errors.New("database error"))

	// Call the method
	err := suite.service.Update(ctx, id, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

// Test Destroy - Success
func (suite *ReportScheduleServiceTestSuite) TestDestroy_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("Destroy", ctx, id, mock.Anything).Return(nil)

	// Call the method
	err := suite.service.Destroy(ctx, id, token)

	// Assertions
	assert.NoError(suite.T(), err)
}

// Test Destroy - Access Denied
func (suite *ReportScheduleServiceTestSuite) TestDestroy_AccessDenied() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA", // Invalid role
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	err := suite.service.Destroy(ctx, id, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user role not allowed", err.Error())
}

// Test Destroy - Repository Error
func (suite *ReportScheduleServiceTestSuite) TestDestroy_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("Destroy", ctx, id, mock.Anything).Return(errors.New("database error"))

	// Call the method
	err := suite.service.Destroy(ctx, id, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

// Test FindByRegistrationID - Success
func (suite *ReportScheduleServiceTestSuite) TestFindByRegistrationID_Success() {
	// Prepare test data
	ctx := context.Background()
	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"

	reportSchedules := []entity.ReportSchedule{
		{
			ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
			UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
			UserNRP:              "5025211111",
			RegistrationID:       registrationID,
			AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
			AcademicAdvisorEmail: "test@gmail.com",
			ReportType:           "FINAL_REPORT",
			Week:                 1,
			StartDate:            func() *time.Time { t := time.Now(); return &t }(),
			EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
			Report: []entity.Report{
				{
					ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
					ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
					Title:                 "Test Report",
					Content:               "Test Content",
					ReportType:            "FINAL_REPORT",
					Feedback:              "Great work",
					AcademicAdvisorStatus: "APPROVED",
					FileStorageID:         "file-123",
				},
			},
		},
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).Return(reportSchedules, nil)

	// Call the method
	result, err := suite.service.FindByRegistrationID(ctx, registrationID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b3", result[0].ID)
	assert.NotNil(suite.T(), result[0].Report)
	assert.Equal(suite.T(), "Test Report", result[0].Report.Title)
}

// Test FindByRegistrationID - Not Found
func (suite *ReportScheduleServiceTestSuite) TestFindByRegistrationID_NotFound() {
	// Prepare test data
	ctx := context.Background()
	registrationID := "non-existent-id"

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).Return([]entity.ReportSchedule{}, errors.New("record not found"))

	// Call the method
	result, err := suite.service.FindByRegistrationID(ctx, registrationID)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
	assert.Nil(suite.T(), result)
}

// Test FindByUserID - Success
func (suite *ReportScheduleServiceTestSuite) TestFindByUserID_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	reportSchedules := []entity.ReportSchedule{
		{
			ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
			UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
			UserNRP:              userNRP,
			RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
			AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
			AcademicAdvisorEmail: "test@gmail.com",
			ReportType:           "FINAL_REPORT",
			Week:                 1,
			StartDate:            func() *time.Time { t := time.Now(); return &t }(),
			EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
			Report: []entity.Report{
				{
					ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
					ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
					Title:                 "Test Report",
					Content:               "Test Content",
					ReportType:            "FINAL_REPORT",
					Feedback:              "Great work",
					AcademicAdvisorStatus: "APPROVED",
					FileStorageID:         "file-123",
				},
			},
		},
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByUserID", ctx, userNRP, mock.Anything).Return(reportSchedules, nil)

	// Call the method
	result, err := suite.service.FindByUserID(ctx, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b3", result[0].ID)
	assert.NotNil(suite.T(), result[0].Report)
	assert.Equal(suite.T(), "Test Report", result[0].Report.Title)
}

// Test FindByUserID - User NRP Not Found
func (suite *ReportScheduleServiceTestSuite) TestFindByUserID_UserNRPNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	usersData := map[string]interface{}{
		"id": "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		// NRP is missing
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.FindByUserID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user NRP not found", err.Error())
	assert.Nil(suite.T(), result)
}

// Test FindByUserID - Repository Error
func (suite *ReportScheduleServiceTestSuite) TestFindByUserID_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Dimas Fadilah",
		"role":  "ADMIN",
		"email": "dimas@example.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByUserID", ctx, userNRP, mock.Anything).Return([]entity.ReportSchedule{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.FindByUserID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Empty(suite.T(), result)
}

// Test FindByAdvisorEmail - Success
func (suite *ReportScheduleServiceTestSuite) TestFindByAdvisorEmail_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	advisorEmail := "advisor@example.com"
	totalCount := int64(1)
	pagReq := dto.PaginationRequest{Limit: 10, Offset: 0}
	reportScheduleReq := dto.ReportScheduleAdvisorRequest{UserNRP: "5025211111"}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "DOSEN PEMBIMBING",
		"email": advisorEmail,
	}

	uuid1, _ := uuid.Parse("9c2fc428-3cca-4c76-a690-e6ba24d135b3")
	uuid2, _ := uuid.Parse("9c2fc428-3cca-4c76-a690-e6ba24d135b5")
	startDate := time.Now().AddDate(1, 0, 0) // One year from now
	endDate := startDate.AddDate(0, 0, 7)    // One week later

	reportScheduleEntity := entity.ReportSchedule{
		ID:                   uuid1,
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: advisorEmail,
		ReportType:           "FINAL_REPORT",
		Week:                 1,
		StartDate:            &startDate,
		EndDate:              &endDate,
		Report: []entity.Report{
			{
				ID:                    uuid2,
				ReportScheduleID:      uuid1.String(),
				Title:                 "Test Report",
				Content:               "Test Content",
				ReportType:            "FINAL_REPORT",
				FileStorageID:         "file-123",
				Feedback:              "Great work",
				AcademicAdvisorStatus: "APPROVED",
			},
		},
	}

	reportSchedulesByUserNRP := map[string][]entity.ReportSchedule{
		"5025211111": {reportScheduleEntity},
	}

	// Expected DTO responses
	expectedReportResponse := &dto.ReportResponse{
		ID:                    uuid2.String(),
		ReportScheduleID:      uuid1.String(),
		Title:                 "Test Report",
		Content:               "Test Content",
		ReportType:            "FINAL_REPORT",
		FileStorageID:         "file-123",
		Feedback:              "Great work",
		AcademicAdvisorStatus: "APPROVED",
	}

	expectedReportScheduleResponse := dto.ReportScheduleResponse{
		ID:                   uuid1.String(),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		ActivityName:         "Test Activity",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: advisorEmail,
		ReportType:           "FINAL_REPORT",
		Week:                 1,
		StartDate:            startDate.Format(time.RFC3339),
		EndDate:              endDate.Format(time.RFC3339),
		Report:               expectedReportResponse,
	}

	expectedResponse := map[string][]dto.ReportScheduleResponse{
		"5025211111": {expectedReportScheduleResponse},
	}

	registration := map[string]interface{}{
		"id":                     "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_nrp":               "5025211111",
		"user_name":              "Dimas Fadilah",
		"academic_advisor":       "test",
		"academic_advisor_email": "advisor@example.com",
		"activity_name":          "Test Activity",
		"approval_status":        "APPROVED",
	}
	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByAdvisorEmailAndGroupByUserID", ctx, advisorEmail, mock.Anything, &pagReq, reportScheduleReq.UserNRP).Return(reportSchedulesByUserNRP, totalCount, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", "9c2fc428-3cca-4c76-a690-e6ba24d135b4", token).Return(registration)
	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmail(ctx, token, pagReq, reportScheduleReq)

	log.Println("RESULT BROOO", result.Reports)
	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dto.ReportScheduleByAdvisorResponse{Reports: expectedResponse}, result)
	assert.Equal(suite.T(), "Test Activity", result.Reports["5025211111"][0].ActivityName)

	// Check pagination
	paginationResp := dto.PaginationResponse{
		Total:       totalCount,
		TotalPages:  1,
		CurrentPage: 1,
		PerPage:     pagReq.Limit,
	}
	assert.Equal(suite.T(), paginationResp.Total, pagination.Total)
}

// Test FindByAdvisorEmail - Advisor Email Not Found
func (suite *ReportScheduleServiceTestSuite) TestFindByAdvisorEmail_AdvisorEmailNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	pagReq := dto.PaginationRequest{Limit: 10, Offset: 0}
	reportScheduleReq := dto.ReportScheduleAdvisorRequest{UserNRP: ""}

	usersData := map[string]interface{}{
		"id":   "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":  "5025211111",
		"name": "Dimas Fadilah",
		"role": "DOSEN PEMBIMBING",
		// Email is missing
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmail(ctx, token, pagReq, reportScheduleReq)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "advisor email not found", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleByAdvisorResponse{}, result)
	assert.Equal(suite.T(), dto.PaginationResponse{}, pagination)
}

// Test FindByAdvisorEmail - Repository Error
func (suite *ReportScheduleServiceTestSuite) TestFindByAdvisorEmail_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	advisorEmail := "advisor@example.com"
	pagReq := dto.PaginationRequest{Limit: 10, Offset: 0}
	reportScheduleReq := dto.ReportScheduleAdvisorRequest{UserNRP: ""}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "DOSEN PEMBIMBING",
		"email": advisorEmail,
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByAdvisorEmailAndGroupByUserID", ctx, advisorEmail, mock.Anything, &pagReq, reportScheduleReq.UserNRP).Return(map[string][]entity.ReportSchedule{}, int64(0), errors.New("database error"))

	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmail(ctx, token, pagReq, reportScheduleReq)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleByAdvisorResponse{}, result)
	assert.Equal(suite.T(), dto.PaginationResponse{}, pagination)
}

// Test FindByAdvisorEmail - Activity Name Not Found
func (suite *ReportScheduleServiceTestSuite) TestFindByAdvisorEmail_ActivityNameNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	advisorEmail := "advisor@example.com"
	pagReq := dto.PaginationRequest{Limit: 10, Offset: 0}
	reportScheduleReq := dto.ReportScheduleAdvisorRequest{UserNRP: ""}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "DOSEN PEMBIMBING",
		"email": advisorEmail,
	}

	userNRP := "5025211111"
	reportSchedules := map[string][]entity.ReportSchedule{
		userNRP: {
			{
				ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
				UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
				UserNRP:              userNRP,
				RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: advisorEmail,
				ReportType:           "FINAL_REPORT",
				Week:                 1,
				StartDate:            func() *time.Time { t := time.Now(); return &t }(),
				EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
			},
		},
	}

	totalCount := int64(1)

	registration := map[string]interface{}{
		"id":                     "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_name":              "Dimas Fadilah",
		"academic_advisor":       "test",
		"academic_advisor_email": advisorEmail,
		// Activity name is missing
		"approval_status": true,
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockReportScheduleRepo.On("FindByAdvisorEmailAndGroupByUserID", ctx, advisorEmail, mock.Anything, &pagReq, reportScheduleReq.UserNRP).Return(reportSchedules, totalCount, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", "9c2fc428-3cca-4c76-a690-e6ba24d135b4", token).Return(registration)

	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmail(ctx, token, pagReq, reportScheduleReq)

	log.Println("ERRORRR", err)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "activity name not found", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleByAdvisorResponse{}, result)
	assert.Equal(suite.T(), dto.PaginationResponse{}, pagination)
}

// Test Index - Success
func (suite *ReportScheduleServiceTestSuite) TestIndex_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	userID := "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0"
	reportSchedules := map[string][]entity.ReportSchedule{
		userID: {
			{
				ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
				UserID:               userID,
				UserNRP:              "5025211111",
				RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: "test@gmail.com",
				ReportType:           "FINAL_REPORT",
				Week:                 1,
				StartDate:            func() *time.Time { t := time.Now(); return &t }(),
				EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
				Report: []entity.Report{
					{
						ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
						ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
						Title:                 "Test Report",
						Content:               "Test Content",
						ReportType:            "FINAL_REPORT",
						Feedback:              "Great work",
						AcademicAdvisorStatus: "APPROVED",
						FileStorageID:         "file-123",
					},
				},
			},
		},
	}

	userData := map[string]interface{}{
		"id":    userID,
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("Index", ctx, mock.Anything).Return(reportSchedules, nil)
	suite.mockUserManagementService.On("GetUserByID", "GET", token, userID).Return(userData)

	// Call the method
	result, err := suite.service.Index(ctx, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Contains(suite.T(), result.Reports, "Dimas Fadilah")
	assert.Len(suite.T(), result.Reports["Dimas Fadilah"], 1)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b3", result.Reports["Dimas Fadilah"][0].ID)
}

// Test Index - Repository Error
func (suite *ReportScheduleServiceTestSuite) TestIndex_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	// Set up expectations
	suite.mockReportScheduleRepo.On("Index", ctx, mock.Anything).Return(map[string][]entity.ReportSchedule{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.Index(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleByAdvisorResponse{}, result)
}

// Test Index - User Name Not Found
func (suite *ReportScheduleServiceTestSuite) TestIndex_UserNameNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	userID := "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0"
	reportSchedules := map[string][]entity.ReportSchedule{
		userID: {
			{
				ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
				UserID:               userID,
				UserNRP:              "5025211111",
				RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: "test@gmail.com",
				ReportType:           "FINAL_REPORT",
				Week:                 1,
				StartDate:            func() *time.Time { t := time.Now(); return &t }(),
				EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
			},
		},
	}

	userData := map[string]interface{}{
		"id":  userID,
		"nrp": "5025211111",
		// Name is missing
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("Index", ctx, mock.Anything).Return(reportSchedules, nil)
	suite.mockUserManagementService.On("GetUserByID", "GET", token, userID).Return(userData)

	// Call the method
	result, err := suite.service.Index(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user name not found", err.Error())
	assert.Equal(suite.T(), dto.ReportScheduleByAdvisorResponse{}, result)
}

// Test ReportScheduleAccess - Admin Role
func (suite *ReportScheduleServiceTestSuite) TestReportScheduleAccess_AdminRole() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorEmail: "test@gmail.com",
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "ADMIN", // Admin role should always have access
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.ReportScheduleAccess(ctx, reportScheduleRequest, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result)
}

// Test ReportScheduleAccess - LO-MBKM Role
func (suite *ReportScheduleServiceTestSuite) TestReportScheduleAccess_LOMBKMRole() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorEmail: "test@gmail.com",
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "LO-MBKM", // LO-MBKM role should always have access
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.ReportScheduleAccess(ctx, reportScheduleRequest, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result)
}

// Test ReportScheduleAccess - Academic Advisor with Matching Email
func (suite *ReportScheduleServiceTestSuite) TestReportScheduleAccess_AdvisorMatchingEmail() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	advisorEmail := "advisor@example.com"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorEmail: advisorEmail,
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "DOSEN PEMBIMBING",
		"email": advisorEmail, // Matching email
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.ReportScheduleAccess(ctx, reportScheduleRequest, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result)
}

// Test ReportScheduleAccess - Academic Advisor with Non-Matching Email
func (suite *ReportScheduleServiceTestSuite) TestReportScheduleAccess_AdvisorNonMatchingEmail() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	usersData := map[string]interface{}{
		"id":   "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":  "5025211111",
		"name": "Dimas Fadilah",
		"role": "DOSEN PEMBIMBING",
		// Email is missing
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.ReportScheduleAccess(ctx, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user email not found", err.Error())
	assert.False(suite.T(), result)
}

// Test ReportScheduleAccess - Invalid Role
func (suite *ReportScheduleServiceTestSuite) TestReportScheduleAccess_InvalidRole() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA", // Invalid role
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.ReportScheduleAccess(ctx, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user role not allowed", err.Error())
	assert.False(suite.T(), result)
}

// Test ReportScheduleAccess - Role Not Found
func (suite *ReportScheduleServiceTestSuite) TestReportScheduleAccess_RoleNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	usersData := map[string]interface{}{
		"id":   "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":  "5025211111",
		"name": "Dimas Fadilah",
		// Role is missing
		"email": "dimasfadilah20@gmail.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.ReportScheduleAccess(ctx, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user role not found", err.Error())
	assert.False(suite.T(), result)
}

// Test ReportScheduleAccess - Advisor Email Not Found
func (suite *ReportScheduleServiceTestSuite) TestReportScheduleAccess_AdvisorEmailNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportScheduleRequest := dto.ReportScheduleRequest{
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	usersData := map[string]interface{}{
		"id":   "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":  "5025211111",
		"name": "Dimas Fadilah",
		"role": "DOSEN PEMBIMBING",
		// Email is missing
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.ReportScheduleAccess(ctx, reportScheduleRequest, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user email not found", err.Error())
	assert.False(suite.T(), result)
}

func (s *mockReportScheduleService) Create(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (dto.ReportScheduleResponse, error) {
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
		return dto.ReportScheduleResponse{}, err
	}
	reportScheduleEntity.StartDate = &start_date
	end_date, err := time.Parse(time.RFC3339, reportSchedule.EndDate)
	if err != nil {
		return dto.ReportScheduleResponse{}, err
	}
	reportScheduleEntity.EndDate = &end_date

	// Set timestamps
	now := time.Now()
	reportScheduleEntity.CreatedAt = &now
	reportScheduleEntity.UpdatedAt = &now

	createdReportSchedule, err := s.reportScheduleRepo.Create(ctx, reportScheduleEntity, nil)
	if err != nil {
		return dto.ReportScheduleResponse{}, err
	}

	return dto.ReportScheduleResponse{
		ID:                   createdReportSchedule.ID.String(),
		UserID:               createdReportSchedule.UserID,
		RegistrationID:       createdReportSchedule.RegistrationID,
		AcademicAdvisorID:    createdReportSchedule.AcademicAdvisorID,
		AcademicAdvisorEmail: createdReportSchedule.AcademicAdvisorEmail,
		ReportType:           createdReportSchedule.ReportType,
		Week:                 createdReportSchedule.Week,
		StartDate:            createdReportSchedule.StartDate.Format(time.RFC3339),
		EndDate:              createdReportSchedule.EndDate.Format(time.RFC3339),
	}, nil
}

func (s *mockReportScheduleService) Update(ctx context.Context, id string, subject dto.ReportScheduleRequest, token string) error {
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

	// Create a new entity with original ID
	reportScheduleEntity := entity.ReportSchedule{
		ID:                   res.ID,
		UserID:               subject.UserID,
		UserNRP:              subject.UserNRP,
		RegistrationID:       subject.RegistrationID,
		AcademicAdvisorID:    subject.AcademicAdvisorID,
		AcademicAdvisorEmail: subject.AcademicAdvisorEmail,
		ReportType:           subject.ReportType,
		Week:                 subject.Week,
	}

	// Handle date parsing
	startDate, err := time.Parse(time.RFC3339, subject.StartDate)
	if err != nil {
		return err
	}
	reportScheduleEntity.StartDate = &startDate

	endDate, err := time.Parse(time.RFC3339, subject.EndDate)
	if err != nil {
		return err
	}
	reportScheduleEntity.EndDate = &endDate

	// Perform the update
	err = s.reportScheduleRepo.Update(ctx, id, reportScheduleEntity, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *mockReportScheduleService) FindByID(ctx context.Context, id string, token string) (dto.ReportScheduleResponse, error) {
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

func (s *mockReportScheduleService) Destroy(ctx context.Context, id string, token string) error {
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

func (s *mockReportScheduleService) FindByRegistrationID(ctx context.Context, registrationID string) ([]dto.ReportScheduleResponse, error) {
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

func (s *mockReportScheduleService) ReportScheduleAccess(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (bool, error) {
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
			return false, errors.New("user email not found")
		}
	} else if userRole != "ADMIN" && userRole != "LO-MBKM" {
		return false, errors.New("user role not allowed")
	}

	return true, nil
}

func (s *mockReportScheduleService) FindByUserID(ctx context.Context, token string) ([]dto.ReportScheduleResponse, error) {
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

func (s *mockReportScheduleService) FindByAdvisorEmail(ctx context.Context, token string, pagReq dto.PaginationRequest, reportScheduleRequest dto.ReportScheduleAdvisorRequest) (dto.ReportScheduleByAdvisorResponse, dto.PaginationResponse, error) {
	user := s.userManagementService.GetUserData("GET", token)
	advisorEmail, ok := user["email"].(string)
	if !ok {
		return dto.ReportScheduleByAdvisorResponse{}, dto.PaginationResponse{}, errors.New("advisor email not found")
	}

	reportSchedules, totalCount, err := s.reportScheduleRepo.FindByAdvisorEmailAndGroupByUserID(ctx, advisorEmail, nil, &pagReq, reportScheduleRequest.UserNRP)
	if err != nil {
		return dto.ReportScheduleByAdvisorResponse{}, dto.PaginationResponse{}, err
	}

	reportScheduleResponses := make(map[string][]dto.ReportScheduleResponse)

	for userNRP, reportScheduleAdvisors := range reportSchedules {
		var reportSchedule []dto.ReportScheduleResponse
		for _, reportScheduleAdvisor := range reportScheduleAdvisors {
			// Get activity name
			activity := s.registrationService.GetRegistrationByID("GET", reportScheduleAdvisor.RegistrationID, token)
			activityName, ok := activity["activity_name"].(string)

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
		reportScheduleResponses[userNRP] = reportSchedule
	}

	// Generate pagination metadata
	paginationResponse := dto.PaginationResponse{
		Total:       totalCount,
		TotalPages:  (totalCount + int64(pagReq.Limit) - 1) / int64(pagReq.Limit),
		CurrentPage: (pagReq.Offset / pagReq.Limit) + 1,
		PerPage:     pagReq.Limit,
	}

	return dto.ReportScheduleByAdvisorResponse{
		Reports: reportScheduleResponses,
	}, paginationResponse, nil
}

func (s *mockReportScheduleService) Index(ctx context.Context, token string) (dto.ReportScheduleByAdvisorResponse, error) {
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
		reportScheduleResponses[userName] = reportSchedule
	}

	return dto.ReportScheduleByAdvisorResponse{
		Reports: reportScheduleResponses,
	}, nil
}

func TestReportScheduleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReportScheduleServiceTestSuite))
}
