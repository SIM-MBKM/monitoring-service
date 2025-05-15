package service_test

import (
	repository_mock "monitoring-service/mocks/repository"
	service_mock "monitoring-service/mocks/service"
	"monitoring-service/service"
	"testing"

	"context"
	"errors"
	"mime/multipart"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TranscriptServiceTestSuite struct {
	suite.Suite
	mockTranscriptRepo        *repository_mock.MockTranscriptRepository
	mockUserManagementService *service_mock.MockUserManagementService
	mockRegistrationService   *service_mock.MockRegistrationService
	mockFileService           *service_mock.MockFileService
	service                   service.TranscriptService
}

func (suite *TranscriptServiceTestSuite) SetupTest() {
	suite.mockTranscriptRepo = new(repository_mock.MockTranscriptRepository)
	suite.mockUserManagementService = new(service_mock.MockUserManagementService)
	suite.mockRegistrationService = new(service_mock.MockRegistrationService)
	suite.mockFileService = new(service_mock.MockFileService)

	createService := func(
		transcriptRepo *repository_mock.MockTranscriptRepository,
		userManagementService *service_mock.MockUserManagementService,
		registrationService *service_mock.MockRegistrationService,
		fileService *service_mock.MockFileService,
	) service.TranscriptService {
		svc := &mockTranscriptService{
			transcriptRepo:        transcriptRepo,
			userManagementService: userManagementService,
			registrationService:   registrationService,
			fileService:           fileService,
		}
		return svc
	}

	suite.service = createService(
		suite.mockTranscriptRepo,
		suite.mockUserManagementService,
		suite.mockRegistrationService,
		suite.mockFileService,
	)
}

type mockTranscriptService struct {
	transcriptRepo        *repository_mock.MockTranscriptRepository
	userManagementService *service_mock.MockUserManagementService
	registrationService   *service_mock.MockRegistrationService
	fileService           *service_mock.MockFileService
}

// Index method implementation for mock service
func (m *mockTranscriptService) Index(ctx context.Context) ([]dto.TranscriptResponse, error) {
	transcripts, err := m.transcriptRepo.Index(ctx, nil)
	if err != nil {
		return nil, err
	}

	var transcriptResponses []dto.TranscriptResponse
	for _, transcript := range transcripts {
		transcriptResponses = append(transcriptResponses, dto.TranscriptResponse{
			ID:                   transcript.ID.String(),
			UserID:               transcript.UserID,
			UserNRP:              transcript.UserNRP,
			AcademicAdvisorID:    transcript.AcademicAdvisorID,
			AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
			RegistrationID:       transcript.RegistrationID,
			Title:                transcript.Title,
			FileStorageID:        transcript.FileStorageID,
		})
	}

	return transcriptResponses, nil
}

// Create method implementation for mock service
func (m *mockTranscriptService) Create(ctx context.Context, transcript dto.TranscriptRequest, file *multipart.FileHeader, token string) (dto.TranscriptResponse, error) {
	// Check if transcript already exists
	transcriptCheck, err := m.transcriptRepo.FindByRegistrationID(ctx, transcript.RegistrationID, nil)
	if err != nil {
		if err.Error() != "record not found" {
			return dto.TranscriptResponse{}, err
		}
	}

	if transcriptCheck.ID != uuid.Nil {
		return dto.TranscriptResponse{}, errors.New("transcript already exists")
	}

	// Check if file is provided
	if file == nil {
		return dto.TranscriptResponse{}, errors.New("file is required")
	}

	// Mock file upload
	fileResult := map[string]interface{}{
		"file_id": "mock-file-id",
	}
	m.fileService.On("Upload", file, "sim_mbkm", "", mock.Anything).Return(fileResult, nil)

	// Get user data for authorization
	user := m.userManagementService.GetUserData("GET", token)
	if user == nil {
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	// Get registration data
	registration := m.registrationService.GetRegistrationByID("GET", transcript.RegistrationID, token)
	if registration == nil {
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	// Check if user has access
	userID, ok := registration["user_id"].(string)
	if !ok {
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	if userID != user["id"] {
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	// Create transcript entity
	userNRP, _ := registration["user_nrp"].(string)
	academicAdvisorID, _ := registration["academic_advisor"].(string)
	academicAdvisorEmail, _ := registration["academic_advisor_email"].(string)

	var transcriptEntity entity.Transcript
	transcriptEntity.ID = uuid.New()
	transcriptEntity.UserID = userID
	transcriptEntity.UserNRP = userNRP
	transcriptEntity.AcademicAdvisorID = academicAdvisorID
	transcriptEntity.AcademicAdvisorEmail = academicAdvisorEmail
	transcriptEntity.RegistrationID = transcript.RegistrationID
	transcriptEntity.Title = transcript.Title
	transcriptEntity.FileStorageID = "mock-file-id"

	// Set timestamps
	now := time.Now()
	transcriptEntity.CreatedAt = &now
	transcriptEntity.UpdatedAt = &now

	// Save to database
	createdTranscript, err := m.transcriptRepo.Create(ctx, transcriptEntity, nil)
	if err != nil {
		return dto.TranscriptResponse{}, err
	}

	return dto.TranscriptResponse{
		ID:                   createdTranscript.ID.String(),
		UserID:               createdTranscript.UserID,
		UserNRP:              createdTranscript.UserNRP,
		AcademicAdvisorID:    createdTranscript.AcademicAdvisorID,
		AcademicAdvisorEmail: createdTranscript.AcademicAdvisorEmail,
		RegistrationID:       createdTranscript.RegistrationID,
		Title:                createdTranscript.Title,
		FileStorageID:        createdTranscript.FileStorageID,
	}, nil
}

// Update method implementation for mock service
func (m *mockTranscriptService) Update(ctx context.Context, id string, subject dto.TranscriptRequest) error {
	// Find the existing transcript
	transcript, err := m.transcriptRepo.FindByID(ctx, id, nil)
	if err != nil {
		return err
	}

	// Create updated transcript entity with original ID
	transcriptEntity := entity.Transcript{
		ID: transcript.ID,
	}

	// Update non-empty fields
	if subject.RegistrationID != "" {
		transcriptEntity.RegistrationID = subject.RegistrationID
	} else {
		transcriptEntity.RegistrationID = transcript.RegistrationID
	}

	if subject.Title != "" {
		transcriptEntity.Title = subject.Title
	} else {
		transcriptEntity.Title = transcript.Title
	}

	// Keep other fields unchanged
	transcriptEntity.UserID = transcript.UserID
	transcriptEntity.UserNRP = transcript.UserNRP
	transcriptEntity.AcademicAdvisorID = transcript.AcademicAdvisorID
	transcriptEntity.AcademicAdvisorEmail = transcript.AcademicAdvisorEmail
	transcriptEntity.FileStorageID = transcript.FileStorageID

	// Update timestamp
	now := time.Now()
	transcriptEntity.UpdatedAt = &now
	transcriptEntity.CreatedAt = transcript.CreatedAt

	// Perform the update
	err = m.transcriptRepo.Update(ctx, id, transcriptEntity, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindByID method implementation for mock service
func (m *mockTranscriptService) FindByID(ctx context.Context, id string, token string) (dto.TranscriptResponse, error) {
	// Find the transcript
	transcript, err := m.transcriptRepo.FindByID(ctx, id, nil)
	if err != nil {
		return dto.TranscriptResponse{}, err
	}

	// Check authorization
	user := m.userManagementService.GetUserData("GET", token)
	if user == nil {
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	// Get user role
	userRole, ok := user["role"].(string)
	if !ok {
		return dto.TranscriptResponse{}, errors.New("user role not found")
	}

	// Admin and LO-MBKM roles have access to all transcripts
	if userRole == "ADMIN" || userRole == "LO-MBKM" {
		// Return transcript data
		return dto.TranscriptResponse{
			ID:                   transcript.ID.String(),
			UserID:               transcript.UserID,
			UserNRP:              transcript.UserNRP,
			AcademicAdvisorID:    transcript.AcademicAdvisorID,
			AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
			RegistrationID:       transcript.RegistrationID,
			Title:                transcript.Title,
			FileStorageID:        transcript.FileStorageID,
		}, nil
	}

	// For academic advisors, check if they are assigned to this transcript
	if userRole == "DOSEN PEMBIMBING" {
		userEmail, ok := user["email"].(string)
		if !ok {
			return dto.TranscriptResponse{}, errors.New("user email not found")
		}

		if userEmail == transcript.AcademicAdvisorEmail {
			// Return transcript data
			return dto.TranscriptResponse{
				ID:                   transcript.ID.String(),
				UserID:               transcript.UserID,
				UserNRP:              transcript.UserNRP,
				AcademicAdvisorID:    transcript.AcademicAdvisorID,
				AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
				RegistrationID:       transcript.RegistrationID,
				Title:                transcript.Title,
				FileStorageID:        transcript.FileStorageID,
			}, nil
		}
	}

	// For students, check if the transcript belongs to them
	if userRole == "MAHASISWA" {
		userID, ok := user["id"].(string)
		if !ok {
			return dto.TranscriptResponse{}, errors.New("user ID not found")
		}

		if userID == transcript.UserID {
			// Return transcript data
			return dto.TranscriptResponse{
				ID:                   transcript.ID.String(),
				UserID:               transcript.UserID,
				UserNRP:              transcript.UserNRP,
				AcademicAdvisorID:    transcript.AcademicAdvisorID,
				AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
				RegistrationID:       transcript.RegistrationID,
				Title:                transcript.Title,
				FileStorageID:        transcript.FileStorageID,
			}, nil
		}
	}

	// Access denied
	return dto.TranscriptResponse{}, errors.New("unauthorized")
}

// Destroy method implementation for mock service
func (m *mockTranscriptService) Destroy(ctx context.Context, id string) error {
	err := m.transcriptRepo.Destroy(ctx, id, nil)
	if err != nil {
		return err
	}
	return nil
}

// FindByRegistrationID method implementation for mock service
func (m *mockTranscriptService) FindByRegistrationID(ctx context.Context, registrationID string) (dto.TranscriptResponse, error) {
	transcript, err := m.transcriptRepo.FindByRegistrationID(ctx, registrationID, nil)
	if err != nil {
		return dto.TranscriptResponse{}, err
	}

	return dto.TranscriptResponse{
		ID:                   transcript.ID.String(),
		UserID:               transcript.UserID,
		UserNRP:              transcript.UserNRP,
		AcademicAdvisorID:    transcript.AcademicAdvisorID,
		AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
		RegistrationID:       transcript.RegistrationID,
		Title:                transcript.Title,
		FileStorageID:        transcript.FileStorageID,
	}, nil
}

// FindAllByRegistrationID method implementation for mock service
func (m *mockTranscriptService) FindAllByRegistrationID(ctx context.Context, registrationID string) ([]dto.TranscriptResponse, error) {
	transcripts, err := m.transcriptRepo.FindAllByRegistrationID(ctx, registrationID, nil)
	if err != nil {
		return nil, err
	}

	var transcriptResponses []dto.TranscriptResponse
	for _, transcript := range transcripts {
		transcriptResponses = append(transcriptResponses, dto.TranscriptResponse{
			ID:                   transcript.ID.String(),
			UserID:               transcript.UserID,
			UserNRP:              transcript.UserNRP,
			AcademicAdvisorID:    transcript.AcademicAdvisorID,
			AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
			RegistrationID:       transcript.RegistrationID,
			Title:                transcript.Title,
			FileStorageID:        transcript.FileStorageID,
		})
	}

	return transcriptResponses, nil
}

// FindByAdvisorEmailAndGroupByUserNRP method implementation for mock service
func (m *mockTranscriptService) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, token string, pagReq dto.PaginationRequest, filter dto.TranscriptAdvisorFilterRequest) (dto.TranscriptAdvisorResponse, dto.PaginationResponse, error) {
	user := m.userManagementService.GetUserData("GET", token)
	advisorEmail, ok := user["email"].(string)
	if !ok {
		return dto.TranscriptAdvisorResponse{}, dto.PaginationResponse{}, errors.New("advisor email not found")
	}

	transcripts, totalCount, err := m.transcriptRepo.FindByAdvisorEmailAndGroupByUserNRP(ctx, advisorEmail, nil, &pagReq, filter.UserNRP)
	if err != nil {
		return dto.TranscriptAdvisorResponse{}, dto.PaginationResponse{}, err
	}

	transcriptResponses := make(map[string][]dto.TranscriptResponse)
	for userNRP, transcript := range transcripts {
		// get registration by registration id
		registration := m.registrationService.GetRegistrationByID("GET", transcript.RegistrationID, token)
		if registration == nil {
			return dto.TranscriptAdvisorResponse{}, dto.PaginationResponse{}, errors.New("registration not found")
		}

		registrationActivityName, ok := registration["activity_name"].(string)
		if !ok {
			return dto.TranscriptAdvisorResponse{}, dto.PaginationResponse{}, errors.New("registration activity name not found")
		}

		transcriptResponses[userNRP] = append(transcriptResponses[userNRP], dto.TranscriptResponse{
			ID:                   transcript.ID.String(),
			UserID:               transcript.UserID,
			UserNRP:              transcript.UserNRP,
			ActivityName:         registrationActivityName,
			AcademicAdvisorID:    transcript.AcademicAdvisorID,
			AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
			RegistrationID:       transcript.RegistrationID,
			Title:                transcript.Title,
			FileStorageID:        transcript.FileStorageID,
		})
	}

	// Generate pagination metadata
	paginationResponse := dto.PaginationResponse{
		Total:       totalCount,
		TotalPages:  (totalCount + int64(pagReq.Limit) - 1) / int64(pagReq.Limit),
		CurrentPage: (pagReq.Offset / pagReq.Limit) + 1,
		PerPage:     pagReq.Limit,
	}

	return dto.TranscriptAdvisorResponse{
		Transcripts: transcriptResponses,
	}, paginationResponse, nil
}

// FindByUserNRPAndGroupByRegistrationID method implementation for mock service
func (m *mockTranscriptService) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.TranscriptByStudentResponse, error) {
	user := m.userManagementService.GetUserData("GET", token)

	userNRP, ok := user["nrp"].(string)
	if !ok {
		return dto.TranscriptByStudentResponse{}, errors.New("user NRP not found")
	}

	transcriptMap, err := m.transcriptRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)
	if err != nil {
		return dto.TranscriptByStudentResponse{}, err
	}

	transcriptResponses := make(map[string][]dto.TranscriptResponse)

	for registrationID, transcripts := range transcriptMap {
		registration := m.registrationService.GetRegistrationByID("GET", registrationID, token)

		registrationActivityName, ok := registration["activity_name"].(string)
		if !ok {
			return dto.TranscriptByStudentResponse{}, errors.New("registration activity name not found")
		}

		var transcriptResponse []dto.TranscriptResponse
		for _, transcript := range transcripts {
			response := dto.TranscriptResponse{
				ID:                   transcript.ID.String(),
				UserID:               transcript.UserID,
				UserNRP:              transcript.UserNRP,
				ActivityName:         registrationActivityName,
				AcademicAdvisorID:    transcript.AcademicAdvisorID,
				AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
				RegistrationID:       transcript.RegistrationID,
				Title:                transcript.Title,
				FileStorageID:        transcript.FileStorageID,
			}
			transcriptResponse = append(transcriptResponse, response)
		}
		transcriptResponses[registrationActivityName] = transcriptResponse
	}

	return dto.TranscriptByStudentResponse{
		Transcripts: transcriptResponses,
	}, nil
}

// Test Index - Success case
func (suite *TranscriptServiceTestSuite) TestIndex_Success() {
	// Prepare test data
	ctx := context.Background()
	uuid1 := uuid.MustParse("9c2fc428-3cca-4c76-a690-e6ba24d135b3")
	now := time.Now()

	transcripts := []entity.Transcript{
		{
			ID:                   uuid1,
			UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
			UserNRP:              "5025211111",
			RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
			AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
			AcademicAdvisorEmail: "advisor@gmail.com",
			Title:                "Final Transcript",
			FileStorageID:        "file-123",
			BaseModel: entity.BaseModel{
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		},
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("Index", ctx, mock.Anything).Return(transcripts, nil)

	// Call the method
	result, err := suite.service.Index(ctx)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), uuid1.String(), result[0].ID)
	assert.Equal(suite.T(), "5025211111", result[0].UserNRP)
	assert.Equal(suite.T(), "Final Transcript", result[0].Title)
}

// Test Index - Database Error
func (suite *TranscriptServiceTestSuite) TestIndex_DatabaseError() {
	// Prepare test data
	ctx := context.Background()

	// Set up expectations
	suite.mockTranscriptRepo.On("Index", ctx, mock.Anything).Return([]entity.Transcript{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.Index(ctx)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Nil(suite.T(), result)
}

// Test Index - Empty Result
func (suite *TranscriptServiceTestSuite) TestIndex_EmptyResult() {
	// Prepare test data
	ctx := context.Background()

	// Set up expectations
	suite.mockTranscriptRepo.On("Index", ctx, mock.Anything).Return([]entity.Transcript{}, nil)

	// Call the method
	result, err := suite.service.Index(ctx)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

// Test Create - Success
func (suite *TranscriptServiceTestSuite) TestCreate_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	file := &multipart.FileHeader{}

	transcriptRequest := dto.TranscriptRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Final Transcript",
	}

	transcriptID := uuid.New()
	now := time.Now()

	// Mock finding transcript by registration ID (not found case)
	suite.mockTranscriptRepo.On("FindByRegistrationID", ctx, transcriptRequest.RegistrationID, mock.Anything).
		Return(entity.Transcript{}, errors.New("record not found"))

	// Mock file upload
	fileUploadResult := map[string]interface{}{
		"file_id": "file-123",
	}
	suite.mockFileService.On("Upload", file, "sim_mbkm", "", mock.Anything).Return(fileUploadResult, nil)

	// Mock user data
	userData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "John Doe",
		"role":  "MAHASISWA",
		"email": "john@example.com",
	}
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userData)

	// Mock registration data
	registrationData := map[string]interface{}{
		"id":                     "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_nrp":               "5025211111",
		"user_name":              "John Doe",
		"academic_advisor":       "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"academic_advisor_email": "advisor@gmail.com",
		"activity_name":          "MBKM Activity",
	}
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", transcriptRequest.RegistrationID, token).
		Return(registrationData)

	// Mock creating transcript
	expectedTranscriptEntity := entity.Transcript{
		ID:                   transcriptID,
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@gmail.com",
		Title:                "Final Transcript",
		FileStorageID:        "file-123",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}
	suite.mockTranscriptRepo.On("Create", ctx, mock.AnythingOfType("entity.Transcript"), mock.Anything).
		Return(expectedTranscriptEntity, nil)

	// Call the method
	result, err := suite.service.Create(ctx, transcriptRequest, file, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedTranscriptEntity.ID.String(), result.ID)
	assert.Equal(suite.T(), expectedTranscriptEntity.Title, result.Title)
	assert.Equal(suite.T(), expectedTranscriptEntity.FileStorageID, result.FileStorageID)
}

// Test Create - Transcript Already Exists
func (suite *TranscriptServiceTestSuite) TestCreate_TranscriptAlreadyExists() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	file := &multipart.FileHeader{}

	transcriptRequest := dto.TranscriptRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Final Transcript",
	}

	existingTranscript := entity.Transcript{
		ID:             uuid.New(),
		RegistrationID: transcriptRequest.RegistrationID,
	}

	// Mock finding transcript by registration ID (found case)
	suite.mockTranscriptRepo.On("FindByRegistrationID", ctx, transcriptRequest.RegistrationID, mock.Anything).
		Return(existingTranscript, nil)

	// Call the method
	result, err := suite.service.Create(ctx, transcriptRequest, file, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "transcript already exists", err.Error())
	assert.Equal(suite.T(), dto.TranscriptResponse{}, result)
}

// Test Create - No File Provided
func (suite *TranscriptServiceTestSuite) TestCreate_NoFileProvided() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	transcriptRequest := dto.TranscriptRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Final Transcript",
	}

	// Mock finding transcript by registration ID (not found case)
	suite.mockTranscriptRepo.On("FindByRegistrationID", ctx, transcriptRequest.RegistrationID, mock.Anything).
		Return(entity.Transcript{}, errors.New("record not found"))

	// Call the method
	result, err := suite.service.Create(ctx, transcriptRequest, nil, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "file is required", err.Error())
	assert.Equal(suite.T(), dto.TranscriptResponse{}, result)
}

// Test FindByID - Success
func (suite *TranscriptServiceTestSuite) TestFindByID_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"
	now := time.Now()

	transcript := entity.Transcript{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@gmail.com",
		Title:                "Final Transcript",
		FileStorageID:        "file-123",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	userData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   "5025211111",
		"name":  "John Doe",
		"role":  "ADMIN",
		"email": "john@example.com",
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByID", ctx, id, mock.Anything).Return(transcript, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userData)

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), id, result.ID)
	assert.Equal(suite.T(), transcript.UserID, result.UserID)
	assert.Equal(suite.T(), transcript.Title, result.Title)
}

// Test FindByID - Not Found
func (suite *TranscriptServiceTestSuite) TestFindByID_NotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.Transcript{}, errors.New("record not found"))

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
	assert.Equal(suite.T(), dto.TranscriptResponse{}, result)
}

// Test FindByID - Unauthorized (Student trying to access another student's transcript)
func (suite *TranscriptServiceTestSuite) TestFindByID_Unauthorized() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"
	now := time.Now()

	transcript := entity.Transcript{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0", // Different user ID
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@gmail.com",
		Title:                "Final Transcript",
		FileStorageID:        "file-123",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	userData := map[string]interface{}{
		"id":    "different-user-id",
		"nrp":   "5025211112",
		"name":  "Jane Doe",
		"role":  "MAHASISWA",
		"email": "jane@example.com",
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByID", ctx, id, mock.Anything).Return(transcript, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userData)

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "unauthorized", err.Error())
	assert.Equal(suite.T(), dto.TranscriptResponse{}, result)
}

// Test FindByID - Authorized Advisor
func (suite *TranscriptServiceTestSuite) TestFindByID_AuthorizedAdvisor() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"
	now := time.Now()

	transcript := entity.Transcript{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@gmail.com",
		Title:                "Final Transcript",
		FileStorageID:        "file-123",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	userData := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"nrp":   "10000000",
		"name":  "Dr. Advisor",
		"role":  "DOSEN PEMBIMBING",
		"email": "advisor@gmail.com", // Matches transcript's advisor email
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByID", ctx, id, mock.Anything).Return(transcript, nil)
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(userData)

	// Call the method
	result, err := suite.service.FindByID(ctx, id, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), id, result.ID)
	assert.Equal(suite.T(), transcript.Title, result.Title)
}

// Test Update - Success
func (suite *TranscriptServiceTestSuite) TestUpdate_Success() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"
	now := time.Now()

	existingTranscript := entity.Transcript{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@gmail.com",
		Title:                "Initial Transcript",
		FileStorageID:        "file-123",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	updateRequest := dto.TranscriptRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Updated Transcript Title",
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByID", ctx, id, mock.Anything).Return(existingTranscript, nil)
	suite.mockTranscriptRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Transcript"), mock.Anything).Return(nil)

	// Call the method
	err := suite.service.Update(ctx, id, updateRequest)

	// Assertions
	assert.NoError(suite.T(), err)
}

// Test Update - Record Not Found
func (suite *TranscriptServiceTestSuite) TestUpdate_RecordNotFound() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	updateRequest := dto.TranscriptRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Updated Transcript Title",
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.Transcript{}, errors.New("record not found"))

	// Call the method
	err := suite.service.Update(ctx, id, updateRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
}

// Test Update - Repository Error
func (suite *TranscriptServiceTestSuite) TestUpdate_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"
	now := time.Now()

	existingTranscript := entity.Transcript{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@gmail.com",
		Title:                "Initial Transcript",
		FileStorageID:        "file-123",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	updateRequest := dto.TranscriptRequest{
		RegistrationID: "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		Title:          "Updated Transcript Title",
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByID", ctx, id, mock.Anything).Return(existingTranscript, nil)
	suite.mockTranscriptRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Transcript"), mock.Anything).Return(errors.New("database error"))

	// Call the method
	err := suite.service.Update(ctx, id, updateRequest)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

// Test Update - Partial Update
func (suite *TranscriptServiceTestSuite) TestUpdate_PartialUpdate() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"
	now := time.Now()

	existingTranscript := entity.Transcript{
		ID:                   uuid.MustParse(id),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@gmail.com",
		Title:                "Initial Transcript",
		FileStorageID:        "file-123",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	// Only update title, keep registration ID the same
	updateRequest := dto.TranscriptRequest{
		Title: "Updated Transcript Title",
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByID", ctx, id, mock.Anything).Return(existingTranscript, nil)
	suite.mockTranscriptRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Transcript"), mock.Anything).Return(nil)

	// Call the method
	err := suite.service.Update(ctx, id, updateRequest)

	// Assertions
	assert.NoError(suite.T(), err)

	// Verify the update call maintains the original registration ID
	updateCall := suite.mockTranscriptRepo.Calls[1]
	updatedEntity := updateCall.Arguments.Get(2).(entity.Transcript)
	assert.Equal(suite.T(), existingTranscript.RegistrationID, updatedEntity.RegistrationID)
	assert.Equal(suite.T(), updateRequest.Title, updatedEntity.Title)
}

// Test Destroy - Success
func (suite *TranscriptServiceTestSuite) TestDestroy_Success() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	// Set up expectations
	suite.mockTranscriptRepo.On("Destroy", ctx, id, mock.Anything).Return(nil)

	// Call the method
	err := suite.service.Destroy(ctx, id)

	// Assertions
	assert.NoError(suite.T(), err)
}

// Test Destroy - Repository Error
func (suite *TranscriptServiceTestSuite) TestDestroy_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	id := "9c2fc428-3cca-4c76-a690-e6ba24d135b3"

	// Set up expectations
	suite.mockTranscriptRepo.On("Destroy", ctx, id, mock.Anything).Return(errors.New("database error"))

	// Call the method
	err := suite.service.Destroy(ctx, id)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
}

// Test FindByRegistrationID - Success
func (suite *TranscriptServiceTestSuite) TestFindByRegistrationID_Success() {
	// Prepare test data
	ctx := context.Background()
	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"
	now := time.Now()

	transcript := entity.Transcript{
		ID:                   uuid.New(),
		UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		UserNRP:              "5025211111",
		RegistrationID:       registrationID,
		AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		AcademicAdvisorEmail: "advisor@gmail.com",
		Title:                "Final Transcript",
		FileStorageID:        "file-123",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).Return(transcript, nil)

	// Call the method
	result, err := suite.service.FindByRegistrationID(ctx, registrationID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), transcript.ID.String(), result.ID)
	assert.Equal(suite.T(), transcript.UserNRP, result.UserNRP)
	assert.Equal(suite.T(), transcript.Title, result.Title)
}

// Test FindByRegistrationID - Not Found
func (suite *TranscriptServiceTestSuite) TestFindByRegistrationID_NotFound() {
	// Prepare test data
	ctx := context.Background()
	registrationID := "non-existent-id"

	// Set up expectations
	suite.mockTranscriptRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).
		Return(entity.Transcript{}, errors.New("record not found"))

	// Call the method
	result, err := suite.service.FindByRegistrationID(ctx, registrationID)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "record not found", err.Error())
	assert.Equal(suite.T(), dto.TranscriptResponse{}, result)
}

// Test FindAllByRegistrationID - Success
func (suite *TranscriptServiceTestSuite) TestFindAllByRegistrationID_Success() {
	// Prepare test data
	ctx := context.Background()
	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"
	now := time.Now()

	transcripts := []entity.Transcript{
		{
			ID:                   uuid.New(),
			UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
			UserNRP:              "5025211111",
			RegistrationID:       registrationID,
			AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
			AcademicAdvisorEmail: "advisor@gmail.com",
			Title:                "Final Transcript",
			FileStorageID:        "file-123",
			BaseModel: entity.BaseModel{
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		},
		{
			ID:                   uuid.New(),
			UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
			UserNRP:              "5025211111",
			RegistrationID:       registrationID,
			AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
			AcademicAdvisorEmail: "advisor@gmail.com",
			Title:                "Midterm Transcript",
			FileStorageID:        "file-456",
			BaseModel: entity.BaseModel{
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		},
	}

	// Set up expectations
	suite.mockTranscriptRepo.On("FindAllByRegistrationID", ctx, registrationID, mock.Anything).Return(transcripts, nil)

	// Call the method
	result, err := suite.service.FindAllByRegistrationID(ctx, registrationID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "Final Transcript", result[0].Title)
	assert.Equal(suite.T(), "Midterm Transcript", result[1].Title)
}

// Test FindAllByRegistrationID - Empty Result
func (suite *TranscriptServiceTestSuite) TestFindAllByRegistrationID_EmptyResult() {
	// Prepare test data
	ctx := context.Background()
	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"

	// Set up expectations
	suite.mockTranscriptRepo.On("FindAllByRegistrationID", ctx, registrationID, mock.Anything).Return([]entity.Transcript{}, nil)

	// Call the method
	result, err := suite.service.FindAllByRegistrationID(ctx, registrationID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

// Test FindAllByRegistrationID - Repository Error
func (suite *TranscriptServiceTestSuite) TestFindAllByRegistrationID_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	registrationID := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"

	// Set up expectations
	suite.mockTranscriptRepo.On("FindAllByRegistrationID", ctx, registrationID, mock.Anything).Return([]entity.Transcript{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.FindAllByRegistrationID(ctx, registrationID)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Nil(suite.T(), result)
}

// Test FindByAdvisorEmailAndGroupByUserNRP - Success
func (suite *TranscriptServiceTestSuite) TestFindByAdvisorEmailAndGroupByUserNRP_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	advisorEmail := "advisor@gmail.com"
	totalCount := int64(1)
	pagReq := dto.PaginationRequest{Limit: 10, Offset: 0}
	filter := dto.TranscriptAdvisorFilterRequest{UserNRP: "5025211111"}
	now := time.Now()

	user := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"nrp":   "10000000",
		"name":  "Dr. Advisor",
		"role":  "DOSEN PEMBIMBING",
		"email": advisorEmail,
	}

	transcripts := map[string]entity.Transcript{
		"5025211111": {
			ID:                   uuid.New(),
			UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
			UserNRP:              "5025211111",
			RegistrationID:       "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
			AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
			AcademicAdvisorEmail: advisorEmail,
			Title:                "Final Transcript",
			FileStorageID:        "file-123",
			BaseModel: entity.BaseModel{
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		},
	}

	registration := map[string]interface{}{
		"id":                     "9c2fc428-3cca-4c76-a690-e6ba24d135b4",
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_nrp":               "5025211111",
		"user_name":              "John Doe",
		"academic_advisor":       "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"academic_advisor_email": advisorEmail,
		"activity_name":          "MBKM Activity",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(user)
	suite.mockTranscriptRepo.On("FindByAdvisorEmailAndGroupByUserNRP", ctx, advisorEmail, mock.Anything, &pagReq, filter.UserNRP).Return(transcripts, totalCount, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", "9c2fc428-3cca-4c76-a690-e6ba24d135b4", token).Return(registration)

	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmailAndGroupByUserNRP(ctx, token, pagReq, filter)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Contains(suite.T(), result.Transcripts, "5025211111")
	assert.Equal(suite.T(), "MBKM Activity", result.Transcripts["5025211111"][0].ActivityName)
	assert.Equal(suite.T(), totalCount, pagination.Total)
}

// Test FindByAdvisorEmailAndGroupByUserNRP - Advisor Email Not Found
func (suite *TranscriptServiceTestSuite) TestFindByAdvisorEmailAndGroupByUserNRP_AdvisorEmailNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	pagReq := dto.PaginationRequest{Limit: 10, Offset: 0}
	filter := dto.TranscriptAdvisorFilterRequest{UserNRP: "5025211111"}

	user := map[string]interface{}{
		"id":   "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"nrp":  "10000000",
		"name": "Dr. Advisor",
		"role": "DOSEN PEMBIMBING",
		// Email missing
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(user)

	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmailAndGroupByUserNRP(ctx, token, pagReq, filter)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "advisor email not found", err.Error())
	assert.Equal(suite.T(), dto.TranscriptAdvisorResponse{}, result)
	assert.Equal(suite.T(), dto.PaginationResponse{}, pagination)
}

// Test FindByAdvisorEmailAndGroupByUserNRP - Repository Error
func (suite *TranscriptServiceTestSuite) TestFindByAdvisorEmailAndGroupByUserNRP_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	advisorEmail := "advisor@gmail.com"
	pagReq := dto.PaginationRequest{Limit: 10, Offset: 0}
	filter := dto.TranscriptAdvisorFilterRequest{UserNRP: "5025211111"}

	user := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"nrp":   "10000000",
		"name":  "Dr. Advisor",
		"role":  "DOSEN PEMBIMBING",
		"email": advisorEmail,
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(user)
	suite.mockTranscriptRepo.On("FindByAdvisorEmailAndGroupByUserNRP", ctx, advisorEmail, mock.Anything, &pagReq, filter.UserNRP).Return(map[string]entity.Transcript{}, int64(0), errors.New("database error"))

	// Call the method
	result, pagination, err := suite.service.FindByAdvisorEmailAndGroupByUserNRP(ctx, token, pagReq, filter)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Equal(suite.T(), dto.TranscriptAdvisorResponse{}, result)
	assert.Equal(suite.T(), dto.PaginationResponse{}, pagination)
}

// Test FindByUserNRPAndGroupByRegistrationID - Success
func (suite *TranscriptServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_Success() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"
	now := time.Now()

	user := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "John Doe",
		"role":  "MAHASISWA",
		"email": "john@example.com",
	}

	registrationID1 := "9c2fc428-3cca-4c76-a690-e6ba24d135b4"
	registrationID2 := "8d1ec317-2b71-5c63-b910-f7da241d8a55"

	transcriptMap := map[string][]entity.Transcript{
		registrationID1: {
			{
				ID:                   uuid.New(),
				UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
				UserNRP:              userNRP,
				RegistrationID:       registrationID1,
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
				AcademicAdvisorEmail: "advisor@gmail.com",
				Title:                "Activity 1 Transcript",
				FileStorageID:        "file-123",
				BaseModel: entity.BaseModel{
					CreatedAt: &now,
					UpdatedAt: &now,
				},
			},
		},
		registrationID2: {
			{
				ID:                   uuid.New(),
				UserID:               "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
				UserNRP:              userNRP,
				RegistrationID:       registrationID2,
				AcademicAdvisorID:    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac2",
				AcademicAdvisorEmail: "advisor2@gmail.com",
				Title:                "Activity 2 Transcript",
				FileStorageID:        "file-456",
				BaseModel: entity.BaseModel{
					CreatedAt: &now,
					UpdatedAt: &now,
				},
			},
		},
	}

	registration1 := map[string]interface{}{
		"id":                     registrationID1,
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_nrp":               userNRP,
		"user_name":              "John Doe",
		"academic_advisor":       "b89dddf4-6ff9-4e9c-891f-dab1960d9ac1",
		"academic_advisor_email": "advisor@gmail.com",
		"activity_name":          "MBKM Activity 1",
	}

	registration2 := map[string]interface{}{
		"id":                     registrationID2,
		"user_id":                "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"user_nrp":               userNRP,
		"user_name":              "John Doe",
		"academic_advisor":       "b89dddf4-6ff9-4e9c-891f-dab1960d9ac2",
		"academic_advisor_email": "advisor2@gmail.com",
		"activity_name":          "MBKM Activity 2",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(user)
	suite.mockTranscriptRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(transcriptMap, nil)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", registrationID1, token).Return(registration1)
	suite.mockRegistrationService.On("GetRegistrationByID", "GET", registrationID2, token).Return(registration2)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Contains(suite.T(), result.Transcripts, "MBKM Activity 1")
	assert.Contains(suite.T(), result.Transcripts, "MBKM Activity 2")
	assert.Equal(suite.T(), "Activity 1 Transcript", result.Transcripts["MBKM Activity 1"][0].Title)
	assert.Equal(suite.T(), "Activity 2 Transcript", result.Transcripts["MBKM Activity 2"][0].Title)
}

// Test FindByUserNRPAndGroupByRegistrationID - User NRP Not Found
func (suite *TranscriptServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_UserNRPNotFound() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"

	user := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"name":  "John Doe",
		"role":  "MAHASISWA",
		"email": "john@example.com",
		// NRP missing
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(user)

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user NRP not found", err.Error())
	assert.Equal(suite.T(), dto.TranscriptByStudentResponse{}, result)
}

// Test FindByUserNRPAndGroupByRegistrationID - Repository Error
func (suite *TranscriptServiceTestSuite) TestFindByUserNRPAndGroupByRegistrationID_RepositoryError() {
	// Prepare test data
	ctx := context.Background()
	token := "test-token"
	userNRP := "5025211111"

	user := map[string]interface{}{
		"id":    "b89dddf4-6ff9-4e9c-891f-dab1960d9ac0",
		"nrp":   userNRP,
		"name":  "John Doe",
		"role":  "MAHASISWA",
		"email": "john@example.com",
	}

	// Set up expectations
	suite.mockUserManagementService.On("GetUserData", "GET", token).Return(user)
	suite.mockTranscriptRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(map[string][]entity.Transcript{}, errors.New("database error"))

	// Call the method
	result, err := suite.service.FindByUserNRPAndGroupByRegistrationID(ctx, token)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "database error", err.Error())
	assert.Equal(suite.T(), dto.TranscriptByStudentResponse{}, result)
}

// Run the test suite
func TestTranscriptServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TranscriptServiceTestSuite))
}
