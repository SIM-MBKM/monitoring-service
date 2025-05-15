package service_test

import (
	"context"
	"errors"
	"mime/multipart"
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

type SyllabusServiceTestSuite struct {
	suite.Suite
	mockSyllabusRepo          *repository_mock.MockSyllabusRepository
	mockUserManagementService *service_mock.MockUserManagementService
	mockRegistrationService   *service_mock.MockRegistrationService
	mockFileService           *service_mock.MockFileService
	service                   service.SyllabusService
}

func (suite *SyllabusServiceTestSuite) SetupTest() {
	suite.mockSyllabusRepo = new(repository_mock.MockSyllabusRepository)
	suite.mockUserManagementService = new(service_mock.MockUserManagementService)
	suite.mockRegistrationService = new(service_mock.MockRegistrationService)
	suite.mockFileService = new(service_mock.MockFileService)

	createService := func(
		syllabusRepo *repository_mock.MockSyllabusRepository,
		userManagementService *service_mock.MockUserManagementService,
		registrationService *service_mock.MockRegistrationService,
		fileService *service_mock.MockFileService,
	) service.SyllabusService {
		svc := &mockSyllabusService{
			syllabusRepo:          syllabusRepo,
			userManagementService: userManagementService,
			registrationService:   registrationService,
			fileService:           fileService,
		}
		return svc
	}

	suite.service = createService(
		suite.mockSyllabusRepo,
		suite.mockUserManagementService,
		suite.mockRegistrationService,
		suite.mockFileService,
	)
}

type mockSyllabusService struct {
	syllabusRepo          *repository_mock.MockSyllabusRepository
	userManagementService *service_mock.MockUserManagementService
	registrationService   *service_mock.MockRegistrationService
	fileService           *service_mock.MockFileService
}

// Implementation of mock methods
func (m *mockSyllabusService) Index(ctx context.Context) ([]dto.SyllabusResponse, error) {
	syllabuses, err := m.syllabusRepo.Index(ctx, nil)
	if err != nil {
		return nil, err
	}

	var syllabusResponses []dto.SyllabusResponse
	for _, syllabus := range syllabuses {
		syllabusResponses = append(syllabusResponses, dto.SyllabusResponse{
			ID:                   syllabus.ID.String(),
			UserID:               syllabus.UserID,
			UserNRP:              syllabus.UserNRP,
			AcademicAdvisorID:    syllabus.AcademicAdvisorID,
			AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
			RegistrationID:       syllabus.RegistrationID,
			Title:                syllabus.Title,
			FileStorageID:        syllabus.FileStorageID,
		})
	}

	return syllabusResponses, nil
}

func (m *mockSyllabusService) Create(ctx context.Context, syllabus dto.SyllabusRequest, file *multipart.FileHeader, token string) (dto.SyllabusResponse, error) {
	syllabusCheck, err := m.syllabusRepo.FindByRegistrationID(ctx, syllabus.RegistrationID, nil)
	if err != nil {
		if err.Error() != "record not found" {
			return dto.SyllabusResponse{}, err
		}
	}

	if syllabusCheck.ID != uuid.Nil {
		return dto.SyllabusResponse{}, errors.New("syllabus already exists")
	}

	if file == nil {
		return dto.SyllabusResponse{}, errors.New("file is required")
	}

	// Upload file to storage
	result, err := m.fileService.Upload(file, "sim_mbkm", "", "")
	if err != nil {
		return dto.SyllabusResponse{}, err
	}

	// Verify user has access to this registration
	user := m.userManagementService.GetUserData("GET", token)
	if user == nil {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	registration := m.registrationService.GetRegistrationByID("GET", syllabus.RegistrationID, token)
	if registration == nil {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	userID, ok := registration["user_id"].(string)
	if !ok {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	if userID != user["id"] {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	// Create syllabus entity
	var syllabusEntity entity.Syllabus
	syllabusEntity.ID = uuid.New()
	syllabusEntity.UserID = userID
	syllabusEntity.UserNRP = registration["user_nrp"].(string)
	syllabusEntity.AcademicAdvisorID = registration["academic_advisor"].(string)
	syllabusEntity.AcademicAdvisorEmail = registration["academic_advisor_email"].(string)
	syllabusEntity.RegistrationID = syllabus.RegistrationID
	syllabusEntity.Title = syllabus.Title
	syllabusEntity.FileStorageID = result["file_id"].(string)

	// Set timestamps
	now := time.Now()
	syllabusEntity.CreatedAt = &now
	syllabusEntity.UpdatedAt = &now

	// Save to database
	syllabusResponse, err := m.syllabusRepo.Create(ctx, syllabusEntity, nil)
	if err != nil {
		return dto.SyllabusResponse{}, err
	}

	return dto.SyllabusResponse{
		ID:                   syllabusResponse.ID.String(),
		UserID:               syllabusResponse.UserID,
		UserNRP:              syllabusResponse.UserNRP,
		AcademicAdvisorID:    syllabusResponse.AcademicAdvisorID,
		AcademicAdvisorEmail: syllabusResponse.AcademicAdvisorEmail,
		RegistrationID:       syllabusResponse.RegistrationID,
		Title:                syllabusResponse.Title,
		FileStorageID:        syllabusResponse.FileStorageID,
	}, nil
}

func (m *mockSyllabusService) Update(ctx context.Context, id string, subject dto.SyllabusRequest) error {
	res, err := m.syllabusRepo.FindByID(ctx, id, nil)
	if err != nil {
		return err
	}

	// Create syllabusEntity with updated fields
	syllabusEntity := entity.Syllabus{
		ID:                   res.ID,
		UserID:               res.UserID,
		UserNRP:              res.UserNRP,
		AcademicAdvisorID:    res.AcademicAdvisorID,
		AcademicAdvisorEmail: res.AcademicAdvisorEmail,
		FileStorageID:        res.FileStorageID,
	}

	// Update fields from request
	if subject.RegistrationID != "" {
		syllabusEntity.RegistrationID = subject.RegistrationID
	} else {
		syllabusEntity.RegistrationID = res.RegistrationID
	}

	if subject.Title != "" {
		syllabusEntity.Title = subject.Title
	} else {
		syllabusEntity.Title = res.Title
	}

	// Update timestamps
	now := time.Now()
	syllabusEntity.UpdatedAt = &now

	// Perform the update
	err = m.syllabusRepo.Update(ctx, id, syllabusEntity, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *mockSyllabusService) FindByID(ctx context.Context, id string, token string) (dto.SyllabusResponse, error) {
	syllabus, err := m.syllabusRepo.FindByID(ctx, id, nil)
	if err != nil {
		return dto.SyllabusResponse{}, err
	}

	user := m.userManagementService.GetUserData("GET", token)
	if user == nil {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	userRole, ok := user["role"].(string)
	if !ok {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	if userRole == "DOSEN PEMBIMBING" {
		userEmail, ok := user["email"].(string)
		if !ok || syllabus.AcademicAdvisorEmail != userEmail {
			return dto.SyllabusResponse{}, errors.New("unauthorized")
		}
	} else if userRole == "MAHASISWA" {
		userNRP, ok := user["nrp"].(string)
		if !ok || syllabus.UserNRP != userNRP {
			return dto.SyllabusResponse{}, errors.New("unauthorized")
		}
	} else if userRole != "ADMIN" && userRole != "LO-MBKM" {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	return dto.SyllabusResponse{
		ID:                   syllabus.ID.String(),
		UserID:               syllabus.UserID,
		UserNRP:              syllabus.UserNRP,
		AcademicAdvisorID:    syllabus.AcademicAdvisorID,
		AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
		RegistrationID:       syllabus.RegistrationID,
		Title:                syllabus.Title,
		FileStorageID:        syllabus.FileStorageID,
	}, nil
}

func (m *mockSyllabusService) Destroy(ctx context.Context, id string) error {
	err := m.syllabusRepo.Destroy(ctx, id, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *mockSyllabusService) FindByRegistrationID(ctx context.Context, registrationID string) (dto.SyllabusResponse, error) {
	syllabus, err := m.syllabusRepo.FindByRegistrationID(ctx, registrationID, nil)
	if err != nil {
		return dto.SyllabusResponse{}, err
	}

	return dto.SyllabusResponse{
		ID:                   syllabus.ID.String(),
		UserID:               syllabus.UserID,
		UserNRP:              syllabus.UserNRP,
		AcademicAdvisorID:    syllabus.AcademicAdvisorID,
		AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
		RegistrationID:       syllabus.RegistrationID,
		Title:                syllabus.Title,
		FileStorageID:        syllabus.FileStorageID,
	}, nil
}

func (m *mockSyllabusService) FindAllByRegistrationID(ctx context.Context, registrationID string) ([]dto.SyllabusResponse, error) {
	syllabuses, err := m.syllabusRepo.FindAllByRegistrationID(ctx, registrationID, nil)
	if err != nil {
		return nil, err
	}

	var syllabusResponses []dto.SyllabusResponse
	for _, syllabus := range syllabuses {
		syllabusResponses = append(syllabusResponses, dto.SyllabusResponse{
			ID:                   syllabus.ID.String(),
			UserID:               syllabus.UserID,
			UserNRP:              syllabus.UserNRP,
			AcademicAdvisorID:    syllabus.AcademicAdvisorID,
			AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
			RegistrationID:       syllabus.RegistrationID,
			Title:                syllabus.Title,
			FileStorageID:        syllabus.FileStorageID,
		})
	}

	return syllabusResponses, nil
}

func (m *mockSyllabusService) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, token string, pagReq dto.PaginationRequest, filter dto.SyllabusAdvisorFilterRequest) (dto.SyllabusAdvisorResponse, dto.PaginationResponse, error) {
	user := m.userManagementService.GetUserData("GET", token)
	advisorEmail, ok := user["email"].(string)
	if !ok {
		return dto.SyllabusAdvisorResponse{}, dto.PaginationResponse{}, errors.New("advisor email not found")
	}

	syllabuses, totalCount, err := m.syllabusRepo.FindByAdvisorEmailAndGroupByUserNRP(ctx, advisorEmail, nil, &pagReq, filter.UserNRP)
	if err != nil {
		return dto.SyllabusAdvisorResponse{}, dto.PaginationResponse{}, err
	}

	syllabusResponses := make(map[string][]dto.SyllabusResponse)
	for userNRP, syllabus := range syllabuses {
		// Get registration details to add activity name
		registration := m.registrationService.GetRegistrationByID("GET", syllabus.RegistrationID, token)

		var activityName string
		if registration != nil {
			activityNameVal, ok := registration["activity_name"].(string)
			if ok {
				activityName = activityNameVal
			}
		}

		syllabusResponses[userNRP] = append(syllabusResponses[userNRP], dto.SyllabusResponse{
			ID:                   syllabus.ID.String(),
			UserID:               syllabus.UserID,
			UserNRP:              syllabus.UserNRP,
			ActivityName:         activityName,
			AcademicAdvisorID:    syllabus.AcademicAdvisorID,
			AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
			RegistrationID:       syllabus.RegistrationID,
			Title:                syllabus.Title,
			FileStorageID:        syllabus.FileStorageID,
		})
	}

	// Generate pagination metadata
	paginationResponse := dto.PaginationResponse{
		Total:       totalCount,
		TotalPages:  (totalCount + int64(pagReq.Limit) - 1) / int64(pagReq.Limit),
		CurrentPage: (pagReq.Offset / pagReq.Limit) + 1,
		PerPage:     pagReq.Limit,
	}

	return dto.SyllabusAdvisorResponse{
		Syllabuses: syllabusResponses,
	}, paginationResponse, nil
}

func (m *mockSyllabusService) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.SyllabusByStudentResponse, error) {
	user := m.userManagementService.GetUserData("GET", token)

	userNRP, ok := user["nrp"].(string)
	if !ok {
		return dto.SyllabusByStudentResponse{}, errors.New("user NRP not found")
	}

	syllabuses, err := m.syllabusRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)
	if err != nil {
		return dto.SyllabusByStudentResponse{}, err
	}

	syllabusResponses := make(map[string][]dto.SyllabusResponse)

	for registrationID, syllabusList := range syllabuses {
		registration := m.registrationService.GetRegistrationByID("GET", registrationID, token)
		registrationActivityName, ok := registration["activity_name"].(string)
		if !ok {
			return dto.SyllabusByStudentResponse{}, errors.New("registration activity name not found")
		}

		approvalStatus, ok := registration["approval_status"].(bool)
		if !ok {
			return dto.SyllabusByStudentResponse{}, errors.New("registration approval status not found")
		}

		if !approvalStatus {
			continue
		}

		var syllabusResponseList []dto.SyllabusResponse
		for _, syllabus := range syllabusList {
			response := dto.SyllabusResponse{
				ID:                   syllabus.ID.String(),
				UserID:               syllabus.UserID,
				UserNRP:              syllabus.UserNRP,
				AcademicAdvisorID:    syllabus.AcademicAdvisorID,
				AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
				RegistrationID:       syllabus.RegistrationID,
				Title:                syllabus.Title,
				FileStorageID:        syllabus.FileStorageID,
			}

			syllabusResponseList = append(syllabusResponseList, response)
		}

		syllabusResponses[registrationActivityName] = syllabusResponseList
	}

	return dto.SyllabusByStudentResponse{
		Syllabuses: syllabusResponses,
	}, nil
}

// Test Index - Success
func (suite *SyllabusServiceTestSuite) TestIndex_Success() {
	// Prepare test data
	ctx := context.Background()
	syllabuses := []entity.Syllabus{
		{
			ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
			UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
			UserNRP:              "5025211111",
			RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
			AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
			AcademicAdvisorEmail: "test@gmail.com",
			Title:                "Test Syllabus",
			FileStorageID:        "file-123",
			BaseModel: entity.BaseModel{
				CreatedAt: func() *time.Time { t := time.Now(); return &t }(),
				UpdatedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
		},
		{
			ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b5"),
			UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac2",
			UserNRP:              "5025211112",
			RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b6",
			AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac3",
			AcademicAdvisorEmail: "test2@gmail.com",
			Title:                "Test Syllabus 2",
			FileStorageID:        "file-456",
			BaseModel: entity.BaseModel{
				CreatedAt: func() *time.Time { t := time.Now(); return &t }(),
				UpdatedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
		},
	}

	// Set up expectations
	suite.mockSyllabusRepo.On("Index", ctx, mock.Anything).Return(syllabuses, nil)

	// Call the method
	result, err := suite.service.Index(ctx)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b3", result[0].ID)
	assert.Equal(suite.T(), "Test Syllabus", result[0].Title)
	assert.Equal(suite.T(), "file-123", result[0].FileStorageID)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b5", result[1].ID)
	assert.Equal(suite.T(), "Test Syllabus 2", result[1].Title)
	assert.Equal(suite.T(), "file-456", result[1].FileStorageID)
}

// Test Index - Database Error
func (suite *SyllabusServiceTestSuite) TestIndex_DatabaseError() {
	// Prepare test data
	ctx := context.Background()

	// Set up expectations
	suite.mockSyllabusRepo.On("Index", ctx, mock.Anything).Return([]entity.Syllabus{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.Index(ctx)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Nil(suite.T(), result)
}

// Test Index - Empty Result
func (suite *SyllabusServiceTestSuite) TestIndex_EmptyResult() {
	// Prepare test data
	ctx := context.Background()

	// Set up expectations
	suite.mockSyllabusRepo.On("Index", ctx, mock.Anything).Return([]entity.Syllabus{}, nil)

	// Call the method
	result, err := suite.service.Index(ctx)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

// Test Create - Success
func (suite *SyllabusServiceTestSuite) TestCreate_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	syllabusRequest := dto.SyllabusRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Test Syllabus",
	}

	file := &multipart.FileHeader{
		Filename: "test.pdf",
		Size:     1024,
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	registration := map[string]interface{}{
		"id":                     "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_nrp":               "5025211111",
		"academic_advisor":       "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"academic_advisor_email": "advisor@example.com",
		"activity_name":          "Test Activity",
		"approval_status":        true,
	}

	fileUploadResult := map[string]interface{}{
		"file_id": "file-123",
	}

	syllabusEntity := entity.Syllabus{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@example.com",
		Title:                "Test Syllabus",
		FileStorageID:        "file-123",
		BaseModel: entity.BaseModel{
			CreatedAt: func() *time.Time { t := time.Now(); return &t }(),
			UpdatedAt: func() *time.Time { t := time.Now(); return &t }(),
		},
	}

	// Set up expectations
	suite.mockSyllabusRepo.On("FindByRegistrationID", ctx, syllabusRequest.RegistrationID, mock.Anything).Return(entity.Syllabus{}, errors.New("record not found"))
	suite.mockFileService.On("Upload", file, "sim_mbkm", "", "").Return(fileUploadResult, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", syllabusRequest.RegistrationID, token).Return(registration)
	suite.mockSyllabusRepo.On("Create", ctx, mock.AnythingOfType("entity.Syllabus"), mock.Anything).Return(syllabusEntity, nil)

	// Call the method
	result, err := suite.service.Create(ctx, syllabusRequest, file, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b3", result.ID)
	assert.Equal(suite.T(), "Test Syllabus", result.Title)
	assert.Equal(suite.T(), "file-123", result.FileStorageID)
}

// Test Create - File Upload Error
func (suite *SyllabusServiceTestSuite) TestCreate_FileUploadError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	syllabusRequest := dto.SyllabusRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Test Syllabus",
	}

	file := &multipart.FileHeader{
		Filename: "test.pdf",
		Size:     1024,
	}

	// Set up expectations
	suite.mockSyllabusRepo.On("FindByRegistrationID", ctx, syllabusRequest.RegistrationID, mock.Anything).Return(entity.Syllabus{}, errors.New("record not found"))
	suite.mockFileService.On("Upload", file, "sim_mbkm", "", "").Return(map[string]interface{}{}, errors.New("file upload error"))

	// Call the method
	result, err := suite.service.Create(ctx, syllabusRequest, file, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "file upload error", err.Error())
	assert.Equal(suite.T(), dto.SyllabusResponse{}, result)
}

// Test Create - Syllabus Already Exists
func (suite *SyllabusServiceTestSuite) TestCreate_SyllabusAlreadyExists() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	syllabusRequest := dto.SyllabusRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Test Syllabus",
	}

	file := &multipart.FileHeader{
		Filename: "test.pdf",
		Size:     1024,
	}

	existingSyllabus := entity.Syllabus{
		ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@example.com",
		Title:                "Existing Syllabus",
		FileStorageID:        "file-123",
	}

	// Set up expectations
	suite.mockSyllabusRepo.On("FindByRegistrationID", ctx, syllabusRequest.RegistrationID, mock.Anything).Return(existingSyllabus, nil)

	// Call the method
	result, err := suite.service.Create(ctx, syllabusRequest, file, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "syllabus already exists", err.Error())
	assert.Equal(suite.T(), dto.SyllabusResponse{}, result)
}

// Test Create - No File Provided
func (suite *SyllabusServiceTestSuite) TestCreate_NoFileProvided() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	syllabusRequest := dto.SyllabusRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Test Syllabus",
	}

	// Set up expectations
	suite.mockSyllabusRepo.On("FindByRegistrationID", ctx, syllabusRequest.RegistrationID, mock.Anything).Return(entity.Syllabus{}, errors.New("record not found"))

	// Call the method
	result, err := suite.service.Create(ctx, syllabusRequest, nil, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "file is required", err.Error())
	assert.Equal(suite.T(), dto.SyllabusResponse{}, result)
}

// Test Create - Unauthorized (User ID Mismatch)
func (suite *SyllabusServiceTestSuite) TestCreate_Unauthorized() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	syllabusRequest := dto.SyllabusRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Test Syllabus",
	}

	file := &multipart.FileHeader{
		Filename: "test.pdf",
		Size:     1024,
	}

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "Dimas Fadilah",
		"role":  "MAHASISWA",
		"email": "dimasfadilah20@gmail.com",
	}

	registration := map[string]interface{}{
		"id":                     "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		"user_id":                "different-user-id", // Different user ID
		"user_nrp":               "5025211111",
		"academic_advisor":       "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"academic_advisor_email": "advisor@example.com",
		"activity_name":          "Test Activity",
		"approval_status":        true,
	}

	fileUploadResult := map[string]interface{}{
		"file_id": "file-123",
	}

	// Set up expectations
	suite.mockSyllabusRepo.On("FindByRegistrationID", ctx, syllabusRequest.RegistrationID, mock.Anything).Return(entity.Syllabus{}, errors.New("record not found"))
	suite.mockFileService.On("Upload", file, "sim_mbkm", "", "").Return(fileUploadResult, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", syllabusRequest.RegistrationID, token).Return(registration)

	// Call the method
	result, err := suite.service.Create(ctx, syllabusRequest, file, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "unauthorized", err.Error())
	assert.Equal(suite.T(), dto.SyllabusResponse{}, result)
}

// Test FindByAdvisorEmailAndGroupByUserNRP - Success
func (suite *SyllabusServiceTestSuite) TestFindByAdvisorEmailAndGroupByUserNRP_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	advisorEmail := "advisor@example.com"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"nrp":   "1001",
		"name":  "Advisor User",
		"role":  "DOSEN PEMBIMBING",
		"email": advisorEmail,
	}

	pagReq := dto.PaginationRequest{
		Limit:  10,
		Offset: 0,
	}

	filter := dto.SyllabusAdvisorFilterRequest{
		UserNRP: userNRP,
	}

	syllabusesByNRP := map[string]entity.Syllabus{
		userNRP: {
			ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
			UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
			UserNRP:              userNRP,
			RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
			AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
			AcademicAdvisorEmail: advisorEmail,
			Title:                "Test Syllabus",
			FileStorageID:        "file-123",
			BaseModel: entity.BaseModel{
				CreatedAt: func() *time.Time { t := time.Now(); return &t }(),
				UpdatedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
		},
	}

	registration := map[string]interface{}{
		"id":                     "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_nrp":               userNRP,
		"academic_advisor":       "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"academic_advisor_email": advisorEmail,
		"activity_name":          "Test Activity",
		"approval_status":        true,
	}

	totalCount := int64(1)

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockSyllabusRepo.On("FindByAdvisorEmailAndGroupByUserNRP", ctx, advisorEmail, mock.Anything, &pagReq, userNRP).Return(syllabusesByNRP, totalCount, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", "9c2fc428-3cca-4c76-a690-e6ba24d135b4", token).Return(registration)

	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmailAndGroupByUserNRP(ctx, token, pagReq, filter)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), totalCount, pagination.Total)
	assert.Equal(suite.T(), int64(1), pagination.TotalPages)
	assert.Equal(suite.T(), 1, pagination.CurrentPage)
	assert.Equal(suite.T(), 10, pagination.PerPage)

	// Verify response structure
	syllabuses := result.Syllabuses
	assert.Contains(suite.T(), syllabuses, userNRP)
	assert.Len(suite.T(), syllabuses[userNRP], 1)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b3", syllabuses[userNRP][0].ID)
	assert.Equal(suite.T(), "Test Syllabus", syllabuses[userNRP][0].Title)
	assert.Equal(suite.T(), "Test Activity", syllabuses[userNRP][0].ActivityName)
}

// Test FindByAdvisorEmailAndGroupByUserNRP - Advisor Email Not Found
func (suite *SyllabusServiceTestSuite) TestFindByAdvisorEmailAndGroupByUserNRP_AdvisorEmailNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	usersData := map[string]interface{}{
		"id":   "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"nrp":  "1001",
		"name": "Advisor User",
		"role": "DOSEN PEMBIMBING",
		// Email is missing
	}

	pagReq := dto.PaginationRequest{
		Limit:  10,
		Offset: 0,
	}

	filter := dto.SyllabusAdvisorFilterRequest{
		UserNRP: "5025211111",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmailAndGroupByUserNRP(ctx, token, pagReq, filter)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "advisor email not found", err.Error())
	assert.Equal(suite.T(), dto.SyllabusAdvisorResponse{}, result)
	assert.Equal(suite.T(), dto.PaginationResponse{}, pagination)
}

// Test FindByAdvisorEmailAndGroupByUserNRP - Repository Error
func (suite *SyllabusServiceTestSuite) TestFindByAdvisorEmailAndGroupByUserNRP_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	advisorEmail := "advisor@example.com"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"nrp":   "1001",
		"name":  "Advisor User",
		"role":  "DOSEN PEMBIMBING",
		"email": advisorEmail,
	}

	pagReq := dto.PaginationRequest{
		Limit:  10,
		Offset: 0,
	}

	filter := dto.SyllabusAdvisorFilterRequest{
		UserNRP: "5025211111",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockSyllabusRepo.On("FindByAdvisorEmailAndGroupByUserNRP", ctx, advisorEmail, mock.Anything, &pagReq, filter.UserNRP).Return(map[string]entity.Syllabus{}, int64(0), errors.New("database error"))

	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmailAndGroupByUserNRP(ctx, token, pagReq, filter)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Equal(suite.T(), dto.SyllabusAdvisorResponse{}, result)
	assert.Equal(suite.T(), dto.PaginationResponse{}, pagination)
}

// Test FindByUserNRPAndGroupByRegistrationID - Success
func (suite *SyllabusServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Student User",
		"role":  "MAHASISWA",
		"email": "student@example.com",
	}

	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"
	syllabuses := map[string][]entity.Syllabus{
		registrationID: {
			{
				ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
				UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
				UserNRP:              userNRP,
				RegistrationID:       registrationID,
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: "advisor@example.com",
				Title:                "Test Syllabus",
				FileStorageID:        "file-123",
				BaseModel: entity.BaseModel{
					CreatedAt: func() *time.Time { t := time.Now(); return &t }(),
					UpdatedAt: func() *time.Time { t := time.Now(); return &t }(),
				},
			},
		},
	}

	registration := map[string]interface{}{
		"id":                     registrationID,
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_nrp":               userNRP,
		"academic_advisor":       "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"academic_advisor_email": "advisor@example.com",
		"activity_name":          "Test Activity",
		"approval_status":        true,
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockSyllabusRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(syllabuses, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", registrationID, token).Return(registration)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Contains(suite.T(), result.Syllabuses, "Test Activity")
	assert.Len(suite.T(), result.Syllabuses["Test Activity"], 1)
	assert.Equal(suite.T(), "9c2fc428-3cca-4c76-a690-e6ba24d135b3", result.Syllabuses["Test Activity"][0].ID)
	assert.Equal(suite.T(), "Test Syllabus", result.Syllabuses["Test Activity"][0].Title)
}

// Test FindByUserNRPAndGroupByRegistrationID - User NRP Not Found
func (suite *SyllabusServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_UserNRPNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	usersData := map[string]interface{}{
		"id": "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		// NRP is missing
		"name":  "Student User",
		"role":  "MAHASISWA",
		"email": "student@example.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user NRP not found", err.Error())
	assert.Equal(suite.T(), dto.SyllabusByStudentResponse{}, result)
}

// Test FindByUserNRPAndGroupByRegistrationID - Repository Error
func (suite *SyllabusServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Student User",
		"role":  "MAHASISWA",
		"email": "student@example.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockSyllabusRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(map[string][]entity.Syllabus{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Equal(suite.T(), dto.SyllabusByStudentResponse{}, result)
}

// Test FindByUserNRPAndGroupByRegistrationID - Activity Name Not Found
func (suite *SyllabusServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_ActivityNameNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	usersData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "Student User",
		"role":  "MAHASISWA",
		"email": "student@example.com",
	}

	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"
	syllabuses := map[string][]entity.Syllabus{
		registrationID: {
			{
				ID:                   uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3"),
				UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
				UserNRP:              userNRP,
				RegistrationID:       registrationID,
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: "advisor@example.com",
				Title:                "Test Syllabus",
				FileStorageID:        "file-123",
			},
		},
	}

	registration := map[string]interface{}{
		"id":                     registrationID,
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_nrp":               userNRP,
		"academic_advisor":       "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"academic_advisor_email": "advisor@example.com",
		// Activity name is missing
		"approval_status": true,
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(usersData)
	suite.mockSyllabusRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(syllabuses, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", registrationID, token).Return(registration)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "registration activity name not found", err.Error())
	assert.Equal(suite.T(), dto.SyllabusByStudentResponse{}, result)
}

func TestSyllabusServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SyllabusServiceTestSuite))
}
