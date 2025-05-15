package service_test

import (
	"monitoring-service/dto"
	"monitoring-service/entity"
	repository_mock "monitoring-service/mocks/repository"
	service_mock "monitoring-service/mocks/service"
	"monitoring-service/service"
	"testing"
	"time"

	"context"
	"errors"

	"mime/multipart"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ReportServiceTestSuite struct {
	suite.Suite
	mockReportRepo            *repository_mock.MockReportRepository
	mockReportScheduleRepo    *repository_mock.MockReportScheduleRepository
	mockUserManagementService *service_mock.MockUserManagementService
	mockFileService           *service_mock.MockFileService
	service                   service.ReportService
}

func (suite *ReportServiceTestSuite) SetupTest() {
	suite.mockReportRepo = new(repository_mock.MockReportRepository)
	suite.mockReportScheduleRepo = new(repository_mock.MockReportScheduleRepository)
	suite.mockUserManagementService = new(service_mock.MockUserManagementService)
	suite.mockFileService = new(service_mock.MockFileService)

	createService := func(
		reportRepo *repository_mock.MockReportRepository,
		reportScheduleRepo *repository_mock.MockReportScheduleRepository,
		userManagementService *service_mock.MockUserManagementService,
		fileService *service_mock.MockFileService,
	) service.ReportService {
		svc := &mockReportService{
			reportRepo:            reportRepo,
			reportScheduleRepo:    reportScheduleRepo,
			userManagementService: userManagementService,
			fileService:           fileService,
		}
		return svc
	}

	suite.service = createService(
		suite.mockReportRepo,
		suite.mockReportScheduleRepo,
		suite.mockUserManagementService,
		suite.mockFileService,
	)
}

type mockReportService struct {
	reportRepo            *repository_mock.MockReportRepository
	reportScheduleRepo    *repository_mock.MockReportScheduleRepository
	userManagementService *service_mock.MockUserManagementService
	fileService           *service_mock.MockFileService
}

// Implement required methods for ReportService interface
func (m *mockReportService) Index(ctx context.Context) ([]dto.ReportResponse, error) {
	reports, err := m.reportRepo.Index(ctx, nil)
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

func (m *mockReportService) Create(ctx context.Context, report dto.ReportRequest, file *multipart.FileHeader, token string) (dto.ReportResponse, error) {
	reportSchedule, err := m.reportScheduleRepo.FindByID(ctx, report.ReportScheduleID, nil)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	user := m.userManagementService.GetUserData("GET", token)

	if user["id"] != reportSchedule.UserID {
		return dto.ReportResponse{}, errors.New("unauthorized")
	}

	var fileID string
	if file != nil {
		result, err := m.fileService.Upload(file, "sim_mbkm", "", "")
		if err != nil {
			return dto.ReportResponse{}, err
		}
		fileID = result["file_id"].(string)
	}

	reportEntity := entity.Report{
		ID:                    uuid.New(),
		ReportScheduleID:      report.ReportScheduleID,
		FileStorageID:         fileID,
		Title:                 report.Title,
		Content:               report.Content,
		ReportType:            report.ReportType,
		AcademicAdvisorStatus: "PENDING",
	}

	createdReport, err := m.reportRepo.Create(ctx, reportEntity, nil)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	return dto.ReportResponse{
		ID:                    createdReport.ID.String(),
		ReportScheduleID:      createdReport.ReportScheduleID,
		FileStorageID:         createdReport.FileStorageID,
		Title:                 createdReport.Title,
		Content:               createdReport.Content,
		ReportType:            createdReport.ReportType,
		Feedback:              createdReport.Feedback,
		AcademicAdvisorStatus: createdReport.AcademicAdvisorStatus,
	}, nil
}

func (m *mockReportService) Update(ctx context.Context, id string, report dto.ReportRequest) error {
	res, err := m.reportRepo.FindByID(ctx, id, nil)
	if err != nil {
		return err
	}

	reportEntity := entity.Report{
		ID:                    res.ID,
		ReportScheduleID:      report.ReportScheduleID,
		Title:                 report.Title,
		Content:               report.Content,
		ReportType:            report.ReportType,
		FileStorageID:         res.FileStorageID,
		Feedback:              res.Feedback,
		AcademicAdvisorStatus: res.AcademicAdvisorStatus,
	}

	// Keep original values if request fields are empty
	if report.Title == "" {
		reportEntity.Title = res.Title
	}
	if report.Content == "" {
		reportEntity.Content = res.Content
	}
	if report.ReportType == "" {
		reportEntity.ReportType = res.ReportType
	}
	if report.ReportScheduleID == "" {
		reportEntity.ReportScheduleID = res.ReportScheduleID
	}

	return m.reportRepo.Update(ctx, id, reportEntity, nil)
}

func (m *mockReportService) FindByID(ctx context.Context, id string, token string) (dto.ReportResponse, error) {
	report, err := m.reportRepo.FindByID(ctx, id, nil)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	reportSchedule, err := m.reportScheduleRepo.FindByID(ctx, report.ReportScheduleID, nil)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	user := m.userManagementService.GetUserData("GET", token)
	userID, _ := user["id"].(string)
	userRole, _ := user["role"].(string)

	if userRole != "ADMIN" && userID != reportSchedule.UserID {
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

func (m *mockReportService) Destroy(ctx context.Context, id string) error {
	return m.reportRepo.Destroy(ctx, id, nil)
}

func (m *mockReportService) FindByReportScheduleID(ctx context.Context, reportScheduleID string) ([]dto.ReportResponse, error) {
	reports, err := m.reportRepo.FindByReportScheduleID(ctx, reportScheduleID, nil)
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

// Implement the Approval method
func (m *mockReportService) Approval(ctx context.Context, token string, report dto.ReportApprovalRequest) error {
	if len(report.IDs) == 0 {
		return errors.New("at least one report ID is required")
	}

	for _, reportID := range report.IDs {
		reportEntity, err := m.reportRepo.FindByID(ctx, reportID, nil)
		if err != nil {
			return err
		}

		reportSchedule, err := m.reportScheduleRepo.FindByID(ctx, reportEntity.ReportScheduleID, nil)
		if err != nil {
			return err
		}

		advisor := m.userManagementService.GetUserData("GET", token)
		advisorEmail, ok := advisor["email"].(string)
		if !ok {
			return errors.New("advisor email not found")
		}

		if advisorEmail != reportSchedule.AcademicAdvisorEmail {
			return errors.New("unauthorized")
		}

		reportEntity.AcademicAdvisorStatus = report.Status
		reportEntity.Feedback = report.Feedback

		err = m.reportRepo.Approval(ctx, reportID, reportEntity, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// Test Index - Success case
func (suite *ReportServiceTestSuite) TestIndex_Success() {
	// Prepare test data
	ctx := context.Background()
	reports := []entity.Report{
		{
			ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
			ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
			Title:                 "Test Report 1",
			Content:               "Test Content 1",
			ReportType:            "WEEKLY_REPORT",
			Feedback:              "Good work",
			AcademicAdvisorStatus: "APPROVED",
			FileStorageID:         "file-123",
		},
		{
			ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b6"),
			ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
			Title:                 "Test Report 2",
			Content:               "Test Content 2",
			ReportType:            "FINAL_REPORT",
			Feedback:              "Need improvements",
			AcademicAdvisorStatus: "REJECTED",
			FileStorageID:         "file-456",
		},
	}

	// Set up expectations
	suite.mockReportRepo.On("Index", ctx, mock.Anything).Return(reports, nil)

	// Call the method
	result, err := suite.service.Index(ctx)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 2, len(result))
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b5", result[0].ID)
	assert.Equal(suite.T(), "Test Report 1", result[0].Title)
	assert.Equal(suite.T(), "WEEKLY_REPORT", result[0].ReportType)
	assert.Equal(suite.T(), "APPROVED", result[0].AcademicAdvisorStatus)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b6", result[1].ID)
	assert.Equal(suite.T(), "Test Report 2", result[1].Title)
	assert.Equal(suite.T(), "FINAL_REPORT", result[1].ReportType)
	assert.Equal(suite.T(), "REJECTED", result[1].AcademicAdvisorStatus)
}

// Test Index - Database Error
func (suite *ReportServiceTestSuite) TestIndex_DatabaseError() {
	// Prepare test data
	ctx := context.Background()

	// Set up expectations
	suite.mockReportRepo.On("Index", ctx, mock.Anything).Return([]entity.Report{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.Index(ctx)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Empty(suite.T(), result)
}

// Test Index - Empty Result
func (suite *ReportServiceTestSuite) TestIndex_EmptyResult() {
	// Prepare test data
	ctx := context.Background()

	// Set up expectations
	suite.mockReportRepo.On("Index", ctx, mock.Anything).Return([]entity.Report{}, nil)

	// Call the method
	result, err := suite.service.Index(ctx)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

// Test FindByID - Success
func (suite *ReportServiceTestSuite) TestFindByID_Success() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"
	token := "test-token"

	report := entity.Report{
		ID:                    uuid.MustParse(id),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:                 "Test Report",
		Content:               "Test Content",
		ReportType:            "WEEKLY_REPORT",
		Feedback:              "Good work",
		AcademicAdvisorStatus: "APPROVED",
		FileStorageID:         "file-123",
	}

	reportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	userAuth := map[string]interface{}{
		"id":    "user-123",
		"email": "user@example.com",
		"role":  "ADMIN",
	}

	// Set up expectations
	suite.mockReportRepo.On("FindByID", ctx, id, mock.Anything).Return(report, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, report.ReportScheduleID, mock.Anything).Return(reportSchedule, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), id, result.ID)
	assert.Equal(suite.T(), "Test Report", result.Title)
	assert.Equal(suite.T(), "Test Content", result.Content)
	assert.Equal(suite.T(), "WEEKLY_REPORT", result.ReportType)
	assert.Equal(suite.T(), "Good work", result.Feedback)
	assert.Equal(suite.T(), "APPROVED", result.AcademicAdvisorStatus)
	assert.Equal(suite.T(), "file-123", result.FileStorageID)
}

// Test FindByID - Report Not Found
func (suite *ReportServiceTestSuite) TestFindByID_ReportNotFound() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"
	token := "test-token"

	// Set up expectations
	suite.mockReportRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.Report{}, errors.New("record not found"))

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
	assert.Equal(suite.T(), dto.ReportResponse{}, result)
}

// Test FindByID - Report Schedule Not Found
func (suite *ReportServiceTestSuite) TestFindByID_ReportScheduleNotFound() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"
	token := "test-token"

	report := entity.Report{
		ID:                    uuid.MustParse(id),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:                 "Test Report",
		Content:               "Test Content",
		ReportType:            "WEEKLY_REPORT",
		Feedback:              "Good work",
		AcademicAdvisorStatus: "APPROVED",
		FileStorageID:         "file-123",
	}

	// Set up expectations
	suite.mockReportRepo.On("FindByID", ctx, id, mock.Anything).Return(report, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, report.ReportScheduleID, mock.Anything).Return(entity.ReportSchedule{}, errors.New("record not found"))

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
	assert.Equal(suite.T(), dto.ReportResponse{}, result)
}

// Test FindByID - Unauthorized User
func (suite *ReportServiceTestSuite) TestFindByID_Unauthorized() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"
	token := "test-token"

	report := entity.Report{
		ID:                    uuid.MustParse(id),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:                 "Test Report",
		Content:               "Test Content",
		ReportType:            "WEEKLY_REPORT",
		Feedback:              "Good work",
		AcademicAdvisorStatus: "APPROVED",
		FileStorageID:         "file-123",
	}

	reportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	userAuth := map[string]interface{}{
		"id":    "different-user", // Different user ID
		"email": "user@example.com",
		"role":  "MAHASISWA", // Not ADMIN or DOSEN
	}

	// Set up expectations
	suite.mockReportRepo.On("FindByID", ctx, id, mock.Anything).Return(report, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, report.ReportScheduleID, mock.Anything).Return(reportSchedule, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "unauthorized", err.Error())
	assert.Equal(suite.T(), dto.ReportResponse{}, result)
}

// Test Create - Success with File
func (suite *ReportServiceTestSuite) TestCreate_SuccessWithFile() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	reportID := uuid.New()
	file := &multipart.FileHeader{} // Mock file

	reportRequest := dto.ReportRequest{
		ReportScheduleID: "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:            "New Test Report",
		Content:          "New Test Content",
		ReportType:       "WEEKLY_REPORT",
	}

	reportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	userAuth := map[string]interface{}{
		"id":    "user-123", // Same as report schedule's UserID
		"email": "user@example.com",
		"role":  "MAHASISWA",
	}

	fileResponse := map[string]interface{}{
		"file_id": "file-123",
		"url":     "https://example.com/file-123",
	}

	createdReport := entity.Report{
		ID:                    reportID,
		ReportScheduleID:      reportRequest.ReportScheduleID,
		Title:                 reportRequest.Title,
		Content:               reportRequest.Content,
		ReportType:            reportRequest.ReportType,
		FileStorageID:         "file-123",
		Feedback:              "",
		AcademicAdvisorStatus: "PENDING",
	}

	// Set up expectations
	suite.mockFileService.On("Upload", file, mock.Anything, mock.Anything, mock.Anything).Return(fileResponse, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, reportRequest.ReportScheduleID, mock.Anything).Return(reportSchedule, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)
	suite.mockReportRepo.On("Create", ctx, mock.AnythingOfType("entity.Report"), mock.Anything).Return(createdReport, nil)

	// Call the method
	result, err := suite.service.Create(ctx, reportRequest, file, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), reportID.String(), result.ID)
	assert.Equal(suite.T(), reportRequest.ReportScheduleID, result.ReportScheduleID)
	assert.Equal(suite.T(), reportRequest.Title, result.Title)
	assert.Equal(suite.T(), reportRequest.Content, result.Content)
	assert.Equal(suite.T(), reportRequest.ReportType, result.ReportType)
	assert.Equal(suite.T(), "file-123", result.FileStorageID)
	assert.Equal(suite.T(), "PENDING", result.AcademicAdvisorStatus)
}

// Test Create - Success without File
func (suite *ReportServiceTestSuite) TestCreate_SuccessWithoutFile() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	reportID := uuid.New()

	reportRequest := dto.ReportRequest{
		ReportScheduleID: "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:            "New Test Report",
		Content:          "New Test Content",
		ReportType:       "WEEKLY_REPORT",
	}

	reportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	userAuth := map[string]interface{}{
		"id":    "user-123", // Same as report schedule's UserID
		"email": "user@example.com",
		"role":  "MAHASISWA",
	}

	createdReport := entity.Report{
		ID:                    reportID,
		ReportScheduleID:      reportRequest.ReportScheduleID,
		Title:                 reportRequest.Title,
		Content:               reportRequest.Content,
		ReportType:            reportRequest.ReportType,
		FileStorageID:         "",
		Feedback:              "",
		AcademicAdvisorStatus: "PENDING",
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByID", ctx, reportRequest.ReportScheduleID, mock.Anything).Return(reportSchedule, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)
	suite.mockReportRepo.On("Create", ctx, mock.AnythingOfType("entity.Report"), mock.Anything).Return(createdReport, nil)

	// Call the method
	result, err := suite.service.Create(ctx, reportRequest, nil, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), reportID.String(), result.ID)
	assert.Equal(suite.T(), reportRequest.ReportScheduleID, result.ReportScheduleID)
	assert.Equal(suite.T(), reportRequest.Title, result.Title)
	assert.Equal(suite.T(), reportRequest.Content, result.Content)
	assert.Equal(suite.T(), reportRequest.ReportType, result.ReportType)
	assert.Empty(suite.T(), result.FileStorageID)
	assert.Equal(suite.T(), "PENDING", result.AcademicAdvisorStatus)
}

// Test Create - File Upload Error
func (suite *ReportServiceTestSuite) TestCreate_FileUploadError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	file := &multipart.FileHeader{} // Mock file

	userData := map[string]interface{}{
		"id":   "9c2fc428-3cca-4c76-a690-e6ba24d135b1",
		"nrp":  "5025211111",
		"name": "Dimas Fadilah",
		"role": "MAHASISWA",
	}

	reportRequest := dto.ReportRequest{
		ReportScheduleID: "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:            "New Test Report",
		Content:          "New Test Content",
		ReportType:       "WEEKLY_REPORT",
	}

	reportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "9c2fc428-3cca-4c76-a690-e6ba24d135b1",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByID", ctx, reportRequest.ReportScheduleID, mock.Anything).Return(reportSchedule, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userData)

	suite.mockFileService = new(service_mock.MockFileService)

	// Set up the file upload to return an error
	suite.mockFileService.On("Upload", file, mock.Anything, mock.Anything, mock.Anything).Return(
		map[string]interface{}{}, errors.New("file upload error"))

	// Recreate the service with our updated mocks
	suite.service = &mockReportService{
		reportRepo:            suite.mockReportRepo,
		reportScheduleRepo:    suite.mockReportScheduleRepo,
		userManagementService: suite.mockUserManagementService,
		fileService:           suite.mockFileService,
	}

	// Call the method
	result, err := suite.service.Create(ctx, reportRequest, file, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "file upload error", err.Error())
	assert.Equal(suite.T(), dto.ReportResponse{}, result)
}

// Test Create - Report Schedule Not Found
func (suite *ReportServiceTestSuite) TestCreate_ReportScheduleNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportRequest := dto.ReportRequest{
		ReportScheduleID: "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:            "New Test Report",
		Content:          "New Test Content",
		ReportType:       "WEEKLY_REPORT",
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByID", ctx, reportRequest.ReportScheduleID, mock.Anything).Return(entity.ReportSchedule{}, errors.New("record not found"))

	// Call the method
	result, err := suite.service.Create(ctx, reportRequest, nil, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
	assert.Equal(suite.T(), dto.ReportResponse{}, result)
}

// Test Create - Unauthorized
func (suite *ReportServiceTestSuite) TestCreate_Unauthorized() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportRequest := dto.ReportRequest{
		ReportScheduleID: "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:            "New Test Report",
		Content:          "New Test Content",
		ReportType:       "WEEKLY_REPORT",
	}

	reportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	userAuth := map[string]interface{}{
		"id":    "different-user", // Different from report schedule's UserID
		"email": "user@example.com",
		"role":  "MAHASISWA",
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByID", ctx, reportRequest.ReportScheduleID, mock.Anything).Return(reportSchedule, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)

	// Call the method
	result, err := suite.service.Create(ctx, reportRequest, nil, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "unauthorized", err.Error())
	assert.Equal(suite.T(), dto.ReportResponse{}, result)
}

// Test Create - Repository Error
func (suite *ReportServiceTestSuite) TestCreate_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportRequest := dto.ReportRequest{
		ReportScheduleID: "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:            "New Test Report",
		Content:          "New Test Content",
		ReportType:       "WEEKLY_REPORT",
	}

	reportSchedule := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	userAuth := map[string]interface{}{
		"id":    "user-123", // Same as report schedule's UserID
		"email": "user@example.com",
		"role":  "MAHASISWA",
	}

	// Set up expectations
	suite.mockReportScheduleRepo.On("FindByID", ctx, reportRequest.ReportScheduleID, mock.Anything).Return(reportSchedule, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)
	suite.mockReportRepo.On("Create", ctx, mock.AnythingOfType("entity.Report"), mock.Anything).Return(entity.Report{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.Create(ctx, reportRequest, nil, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Equal(suite.T(), dto.ReportResponse{}, result)
}

// Test Update - Success
func (suite *ReportServiceTestSuite) TestUpdate_Success() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"

	reportRequest := dto.ReportRequest{
		ReportScheduleID: "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:            "Updated Test Report",
		Content:          "Updated Test Content",
		ReportType:       "FINAL_REPORT",
	}

	existingReport := entity.Report{
		ID:                    uuid.MustParse(id),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:                 "Original Test Report",
		Content:               "Original Test Content",
		ReportType:            "WEEKLY_REPORT",
		FileStorageID:         "file-123",
		Feedback:              "Good work",
		AcademicAdvisorStatus: "APPROVED",
	}

	// Set up expectations
	suite.mockReportRepo.On("FindByID", ctx, id, mock.Anything).Return(existingReport, nil)
	suite.mockReportRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Report"), mock.Anything).Return(nil)

	// Call the method
	err := suite.service.Update(ctx, id, reportRequest)

	// Assertions
	assert.NoError(suite.T(), err)
}

// Test Update - Report Not Found
func (suite *ReportServiceTestSuite) TestUpdate_ReportNotFound() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"

	reportRequest := dto.ReportRequest{
		ReportScheduleID: "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:            "Updated Test Report",
		Content:          "Updated Test Content",
		ReportType:       "FINAL_REPORT",
	}

	// Set up expectations
	suite.mockReportRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.Report{}, errors.New("record not found"))

	// Call the method
	err := suite.service.Update(ctx, id, reportRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
}

// Test Update - Repository Error
func (suite *ReportServiceTestSuite) TestUpdate_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"

	reportRequest := dto.ReportRequest{
		ReportScheduleID: "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:            "Updated Test Report",
		Content:          "Updated Test Content",
		ReportType:       "FINAL_REPORT",
	}

	existingReport := entity.Report{
		ID:                    uuid.MustParse(id),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:                 "Original Test Report",
		Content:               "Original Test Content",
		ReportType:            "WEEKLY_REPORT",
		FileStorageID:         "file-123",
		Feedback:              "Good work",
		AcademicAdvisorStatus: "APPROVED",
	}

	// Set up expectations
	suite.mockReportRepo.On("FindByID", ctx, id, mock.Anything).Return(existingReport, nil)
	suite.mockReportRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Report"), mock.Anything).Return(errors.New("database error"))

	// Call the method
	err := suite.service.Update(ctx, id, reportRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

// Test Update - Partial Update
func (suite *ReportServiceTestSuite) TestUpdate_PartialUpdate() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"

	// Only updating the title
	reportRequest := dto.ReportRequest{
		Title:      "Updated Test Report",
		ReportType: "", // Empty field should be ignored
		Content:    "", // Empty field should be ignored
	}

	existingReport := entity.Report{
		ID:                    uuid.MustParse(id),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:                 "Original Test Report",
		Content:               "Original Test Content",
		ReportType:            "WEEKLY_REPORT",
		FileStorageID:         "file-123",
		Feedback:              "Good work",
		AcademicAdvisorStatus: "APPROVED",
	}

	// Set up expectations
	suite.mockReportRepo.On("FindByID", ctx, id, mock.Anything).Return(existingReport, nil)

	// Mock the update to verify the expected entity
	suite.mockReportRepo.On("Update", ctx, id, mock.MatchedBy(func(r entity.Report) bool {
		// The title should be updated but other fields should be preserved
		return r.Title == "Updated Test Report" &&
			r.Content == existingReport.Content &&
			r.ReportType == existingReport.ReportType
	}), mock.Anything).Return(nil)

	// Call the method
	err := suite.service.Update(ctx, id, reportRequest)

	// Assertions
	assert.NoError(suite.T(), err)
}

// Test Destroy - Success
func (suite *ReportServiceTestSuite) TestDestroy_Success() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"

	// Set up expectations
	suite.mockReportRepo.On("Destroy", ctx, id, mock.Anything).Return(nil)

	// Call the method
	err := suite.service.Destroy(ctx, id)

	// Assertions
	assert.NoError(suite.T(), err)
}

// Test Destroy - Repository Error
func (suite *ReportServiceTestSuite) TestDestroy_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b5"

	// Set up expectations
	suite.mockReportRepo.On("Destroy", ctx, id, mock.Anything).Return(errors.New("database error"))

	// Call the method
	err := suite.service.Destroy(ctx, id)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

// Test FindByReportScheduleID - Success
func (suite *ReportServiceTestSuite) TestFindByReportScheduleID_Success() {
	// Prepare test data
	ctx := context.Background()
	reportScheduleID := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	reports := []entity.Report{
		{
			ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
			ReportScheduleID:      reportScheduleID,
			Title:                 "Test Report 1",
			Content:               "Test Content 1",
			ReportType:            "WEEKLY_REPORT",
			Feedback:              "Good work",
			AcademicAdvisorStatus: "APPROVED",
			FileStorageID:         "file-123",
		},
		{
			ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b6"),
			ReportScheduleID:      reportScheduleID,
			Title:                 "Test Report 2",
			Content:               "Test Content 2",
			ReportType:            "FINAL_REPORT",
			Feedback:              "Need improvements",
			AcademicAdvisorStatus: "REJECTED",
			FileStorageID:         "file-456",
		},
	}

	// Set up expectations
	suite.mockReportRepo.On("FindByReportScheduleID", ctx, reportScheduleID, mock.Anything).Return(reports, nil)

	// Call the method
	result, err := suite.service.FindByReportScheduleID(ctx, reportScheduleID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 2, len(result))
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b5", result[0].ID)
	assert.Equal(suite.T(), "Test Report 1", result[0].Title)
	assert.Equal(suite.T(), "WEEKLY_REPORT", result[0].ReportType)
	assert.Equal(suite.T(), "APPROVED", result[0].AcademicAdvisorStatus)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b6", result[1].ID)
	assert.Equal(suite.T(), "Test Report 2", result[1].Title)
	assert.Equal(suite.T(), "FINAL_REPORT", result[1].ReportType)
	assert.Equal(suite.T(), "REJECTED", result[1].AcademicAdvisorStatus)
}

// Test FindByReportScheduleID - Repository Error
func (suite *ReportServiceTestSuite) TestFindByReportScheduleID_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	reportScheduleID := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	// Set up expectations
	suite.mockReportRepo.On("FindByReportScheduleID", ctx, reportScheduleID, mock.Anything).Return([]entity.Report{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.FindByReportScheduleID(ctx, reportScheduleID)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Empty(suite.T(), result)
}

// Test FindByReportScheduleID - Empty Result
func (suite *ReportServiceTestSuite) TestFindByReportScheduleID_EmptyResult() {
	// Prepare test data
	ctx := context.Background()
	reportScheduleID := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	// Set up expectations
	suite.mockReportRepo.On("FindByReportScheduleID", ctx, reportScheduleID, mock.Anything).Return([]entity.Report{}, nil)

	// Call the method
	result, err := suite.service.FindByReportScheduleID(ctx, reportScheduleID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

// Test Approval - Success
func (suite *ReportServiceTestSuite) TestApproval_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportApprovalRequest := dto.ReportApprovalRequest{
		Status:   "APPROVED",
		Feedback: "Great work!",
		IDs:      []string{"9c2fc428-3cca-4c76-a690-e6ba24d135b5", "9c2fc428-3cca-4c76-a690-e6ba24d135b6"},
	}

	report1 := entity.Report{
		ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b1",
		Title:                 "Test Report 1",
		Content:               "Test Content 1",
		ReportType:            "WEEKLY_REPORT",
		Feedback:              "",
		AcademicAdvisorStatus: "PENDING",
		FileStorageID:         "file-123",
	}

	report2 := entity.Report{
		ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b6"),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b2",
		Title:                 "Test Report 2",
		Content:               "Test Content 2",
		ReportType:            "FINAL_REPORT",
		Feedback:              "",
		AcademicAdvisorStatus: "PENDING",
		FileStorageID:         "file-456",
	}

	schedule1 := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b1"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com", // Same as user's email
	}

	schedule2 := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b2"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com", // Same as user's email
	}

	userAuth := map[string]interface{}{
		"id":    "advisor-123",
		"email": "advisor@example.com", // Same as academic advisor email
		"role":  "DOSEN PEMBIMBING",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)

	// First report
	suite.mockReportRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b5", mock.Anything).Return(report1, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b1", mock.Anything).Return(schedule1, nil)
	suite.mockReportRepo.On("Approval", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b5", mock.MatchedBy(func(r entity.Report) bool {
		return r.AcademicAdvisorStatus == "APPROVED" && r.Feedback == "Great work!"
	}), mock.Anything).Return(nil)

	// Second report
	suite.mockReportRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b6", mock.Anything).Return(report2, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b2", mock.Anything).Return(schedule2, nil)
	suite.mockReportRepo.On("Approval", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b6", mock.MatchedBy(func(r entity.Report) bool {
		return r.AcademicAdvisorStatus == "APPROVED" && r.Feedback == "Great work!"
	}), mock.Anything).Return(nil)

	// Call the method
	err := suite.service.Approval(ctx, token, reportApprovalRequest)

	// Assertions
	assert.NoError(suite.T(), err)
}

// Test Approval - No Report IDs Provided
func (suite *ReportServiceTestSuite) TestApproval_NoReportIDs() {
	ctx := context.Background()
	token := "test-token"

	reportApprovalRequest := dto.ReportApprovalRequest{
		Status:   "APPROVED",
		Feedback: "Great work!",
		IDs:      []string{},
	}

	err := suite.service.Approval(ctx, token, reportApprovalRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "at least one report ID is required", err.Error())
}

// Test Approval - Advisor Email Not Found
func (suite *ReportServiceTestSuite) TestApproval_AdvisorEmailNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	uuid1 := uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5")
	uuid2 := uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b6")
	uuid3 := uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3")

	reportApprovalRequest := dto.ReportApprovalRequest{
		Status:   "APPROVED",
		Feedback: "Great work!",
		IDs:      []string{"9c2fc428-3cca-4c76-a690-e6ba24d135b5"},
	}

	userData := map[string]interface{}{
		"id":   uuid2.String(),
		"nrp":  "5025211111",
		"name": "Dimas Fadilah",
		"role": "MAHASISWA",
		// "email": "dimas@gmail.com",
	}

	reportEntity := entity.Report{
		ID:                    uuid1,
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:                 "Test Report 1",
		Content:               "Test Content 1",
		ReportType:            "WEEKLY_REPORT",
		Feedback:              "Great work!",
		AcademicAdvisorStatus: "PENDING",
		FileStorageID:         "file-123",
	}

	reportScheduleEntity := entity.ReportSchedule{
		ID:                   uuid3,
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com",
		ReportType:           "WEEKLY_REPORT",
		Week:                 1,
		StartDate:            func() *time.Time { t := time.Now(); return &t }(),
		EndDate:              func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
		Report:               []entity.Report{reportEntity},
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userData)
	suite.mockReportRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b5", mock.Anything).Return(reportEntity, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, reportEntity.ReportScheduleID, mock.Anything).Return(reportScheduleEntity, nil)

	// Call the method
	err := suite.service.Approval(ctx, token, reportApprovalRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "advisor email not found", err.Error())
}

// Test Approval - Report Not Found
func (suite *ReportServiceTestSuite) TestApproval_ReportNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportApprovalRequest := dto.ReportApprovalRequest{
		Status:   "APPROVED",
		Feedback: "Great work!",
		IDs:      []string{"9c2fc428-3cca-4c76-a690-e6ba24d135b5"},
	}

	userAuth := map[string]interface{}{
		"id":    "advisor-123",
		"email": "advisor@example.com",
		"role":  "DOSEN PEMBIMBING",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)
	suite.mockReportRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b5", mock.Anything).Return(entity.Report{}, errors.New("record not found"))

	// Call the method
	err := suite.service.Approval(ctx, token, reportApprovalRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
}

// Test Approval - Report Schedule Not Found
func (suite *ReportServiceTestSuite) TestApproval_ReportScheduleNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportApprovalRequest := dto.ReportApprovalRequest{
		Status:   "APPROVED",
		Feedback: "Great work!",
		IDs:      []string{"9c2fc428-3cca-4c76-a690-e6ba24d135b5"},
	}

	report := entity.Report{
		ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
		ReportScheduleID:      "schedule-1",
		Title:                 "Test Report",
		Content:               "Test Content",
		ReportType:            "WEEKLY_REPORT",
		Feedback:              "",
		AcademicAdvisorStatus: "PENDING",
		FileStorageID:         "file-123",
	}

	userAuth := map[string]interface{}{
		"id":    "advisor-123",
		"email": "advisor@example.com",
		"role":  "DOSEN PEMBIMBING",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)
	suite.mockReportRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b5", mock.Anything).Return(report, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, "schedule-1", mock.Anything).Return(entity.ReportSchedule{}, errors.New("record not found"))

	// Call the method
	err := suite.service.Approval(ctx, token, reportApprovalRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
}

// Test Approval - Unauthorized (Email Mismatch)
func (suite *ReportServiceTestSuite) TestApproval_Unauthorized() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportApprovalRequest := dto.ReportApprovalRequest{
		Status:   "APPROVED",
		Feedback: "Great work!",
		IDs:      []string{"9c2fc428-3cca-4c76-a690-e6ba24d135b5"},
	}

	report := entity.Report{
		ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b1",
		Title:                 "Test Report",
		Content:               "Test Content",
		ReportType:            "WEEKLY_REPORT",
		Feedback:              "",
		AcademicAdvisorStatus: "PENDING",
		FileStorageID:         "file-123",
	}

	schedule := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b1"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com",
	}

	userAuth := map[string]interface{}{
		"id":    "advisor-123",
		"email": "different@example.com", // Different from academic advisor email
		"role":  "DOSEN PEMBIMBING",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)
	suite.mockReportRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b5", mock.Anything).Return(report, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b1", mock.Anything).Return(schedule, nil)

	// Call the method
	err := suite.service.Approval(ctx, token, reportApprovalRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "unauthorized", err.Error())
}

// Test Approval - Repository Error
func (suite *ReportServiceTestSuite) TestApproval_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	reportApprovalRequest := dto.ReportApprovalRequest{
		Status:   "APPROVED",
		Feedback: "Great work!",
		IDs:      []string{"9c2fc428-3cca-4c76-a690-e6ba24d135b5"},
	}

	report := entity.Report{
		ID:                    uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
		ReportScheduleID:      "9c2fc428-3cca-4c76-a690-e6ba24d135b3",
		Title:                 "Test Report",
		Content:               "Test Content",
		ReportType:            "WEEKLY_REPORT",
		Feedback:              "",
		AcademicAdvisorStatus: "PENDING",
		FileStorageID:         "file-123",
	}

	schedule := entity.ReportSchedule{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "user-123",
		UserNRP:              "5025211111",
		RegistrationID:       "reg-123",
		AcademicAdvisorID:    "advisor-123",
		AcademicAdvisorEmail: "advisor@example.com", // Same as user's email
	}

	userAuth := map[string]interface{}{
		"id":    "advisor-123",
		"email": "advisor@example.com", // Same as academic advisor email
		"role":  "DOSEN PEMBIMBING",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userAuth)
	suite.mockReportRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b5", mock.Anything).Return(report, nil)
	suite.mockReportScheduleRepo.On("FindByID", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b3", mock.Anything).Return(schedule, nil)
	suite.mockReportRepo.On("Approval", ctx, "9c2fc428-3cca-4c76-a690-e6ba24d135b5", mock.MatchedBy(func(r entity.Report) bool {
		return r.AcademicAdvisorStatus == "APPROVED" && r.Feedback == "Great work!"
	}), mock.Anything).Return(errors.New("database error"))

	// Call the method
	err := suite.service.Approval(ctx, token, reportApprovalRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

func TestReportServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReportServiceTestSuite))
}
