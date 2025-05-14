package repository_test

import (
	"context"
	"errors"
	"monitoring-service/dto"
	"monitoring-service/entity"
	repository_mock "monitoring-service/mocks/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockTranscript() entity.Transcript {
	now := time.Now()

	return entity.Transcript{
		ID:                   uuid.New(),
		UserNRP:              "5123123123",
		RegistrationID:       uuid.NewString(),
		FileStorageID:        uuid.NewString(),
		AcademicAdvisorEmail: "advisor@example.com",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}
}

func TestTranscriptRepository_Index(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	mockTranscripts := []entity.Transcript{createMockTranscript(), createMockTranscript()}

	mockRepo.On("Index", ctx, mock.Anything).Return(mockTranscripts, nil)

	result, err := mockRepo.Index(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockTranscripts, result)
}

func TestTranscriptRepository_Index_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()

	mockRepo.On("Index", ctx, mock.Anything).Return([]entity.Transcript{}, errors.New("error"))

	result, err := mockRepo.Index(ctx, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestTranscriptRepository_Create(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	mockTranscript := createMockTranscript()

	mockRepo.On("Create", ctx, mock.AnythingOfType("entity.Transcript"), mock.Anything).Return(mockTranscript, nil)

	result, err := mockRepo.Create(ctx, mockTranscript, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockTranscript.ID, result.ID)
}

func TestTranscriptRepository_Create_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	mockTranscript := createMockTranscript()

	mockRepo.On("Create", ctx, mock.AnythingOfType("entity.Transcript"), mock.Anything).Return(entity.Transcript{}, errors.New("error"))

	result, err := mockRepo.Create(ctx, mockTranscript, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.Transcript{}, result)
}

func TestTranscriptRepository_Update(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	mockTranscript := createMockTranscript()
	id := mockTranscript.ID.String()

	mockRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Transcript"), mock.Anything).Return(nil)

	err := mockRepo.Update(ctx, id, mockTranscript, nil)

	assert.NoError(t, err)
}

func TestTranscriptRepository_Update_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	mockTranscript := createMockTranscript()
	id := mockTranscript.ID.String()

	mockRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Transcript"), mock.Anything).Return(errors.New("error"))

	err := mockRepo.Update(ctx, id, mockTranscript, nil)

	assert.Error(t, err)
}

func TestTranscriptRepository_FindByID(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	mockTranscript := createMockTranscript()
	id := mockTranscript.ID.String()

	mockRepo.On("FindByID", ctx, id, mock.Anything).Return(mockTranscript, nil)

	result, err := mockRepo.FindByID(ctx, id, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockTranscript.ID, result.ID)
}

func TestTranscriptRepository_FindByID_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	id := uuid.New().String()

	mockRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.Transcript{}, errors.New("error"))

	result, err := mockRepo.FindByID(ctx, id, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.Transcript{}, result)
}

func TestTranscriptRepository_Destroy(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	id := uuid.New().String()

	mockRepo.On("Destroy", ctx, id, mock.Anything).Return(nil)

	err := mockRepo.Destroy(ctx, id, nil)

	assert.NoError(t, err)
}

func TestTranscriptRepository_Destroy_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	id := uuid.New().String()

	mockRepo.On("Destroy", ctx, id, mock.Anything).Return(errors.New("error"))

	err := mockRepo.Destroy(ctx, id, nil)

	assert.Error(t, err)
}

func TestTranscriptRepository_FindByRegistrationID(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()
	mockTranscript := createMockTranscript()

	mockRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).Return(mockTranscript, nil)

	result, err := mockRepo.FindByRegistrationID(ctx, registrationID, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockTranscript.ID, result.ID)
}

func TestTranscriptRepository_FindByRegistrationID_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()

	mockRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).Return(entity.Transcript{}, errors.New("error"))

	result, err := mockRepo.FindByRegistrationID(ctx, registrationID, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.Transcript{}, result)
}

func TestTranscriptRepository_FindAllByRegistrationID(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()
	mockTranscripts := []entity.Transcript{createMockTranscript(), createMockTranscript()}

	mockRepo.On("FindAllByRegistrationID", ctx, registrationID, mock.Anything).Return(mockTranscripts, nil)

	result, err := mockRepo.FindAllByRegistrationID(ctx, registrationID, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockTranscripts, result)
}

func TestTranscriptRepository_FindAllByRegistrationID_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()

	mockRepo.On("FindAllByRegistrationID", ctx, registrationID, mock.Anything).Return([]entity.Transcript{}, errors.New("error"))

	result, err := mockRepo.FindAllByRegistrationID(ctx, registrationID, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestTranscriptRepository_FindByAdvisorEmailAndGroupByUserNRP(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	advisorEmail := "advisor@example.com"
	pagReq := &dto.PaginationRequest{Limit: 10, Offset: 0}
	userNRPFilter := ""
	mockTranscriptMap := map[string]entity.Transcript{
		"5123123123": createMockTranscript(),
	}
	totalCount := int64(1)

	mockRepo.On("FindByAdvisorEmailAndGroupByUserNRP", ctx, advisorEmail, mock.Anything, pagReq, userNRPFilter).Return(mockTranscriptMap, totalCount, nil)

	result, count, err := mockRepo.FindByAdvisorEmailAndGroupByUserNRP(ctx, advisorEmail, nil, pagReq, userNRPFilter)

	assert.NoError(t, err)
	assert.Equal(t, mockTranscriptMap, result)
	assert.Equal(t, totalCount, count)
}

func TestTranscriptRepository_FindByAdvisorEmailAndGroupByUserNRP_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	advisorEmail := "advisor@example.com"
	pagReq := &dto.PaginationRequest{Limit: 10, Offset: 0}
	userNRPFilter := ""

	mockRepo.On("FindByAdvisorEmailAndGroupByUserNRP", ctx, advisorEmail, mock.Anything, pagReq, userNRPFilter).Return(map[string]entity.Transcript{}, int64(0), errors.New("error"))

	result, count, err := mockRepo.FindByAdvisorEmailAndGroupByUserNRP(ctx, advisorEmail, nil, pagReq, userNRPFilter)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), count)
}

func TestTranscriptRepository_FindByUserNRPAndGroupByRegistrationID(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	userNRP := "5123123123"
	mockTranscriptMap := map[string][]entity.Transcript{
		"reg-id-1": {createMockTranscript(), createMockTranscript()},
	}

	mockRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(mockTranscriptMap, nil)

	result, err := mockRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockTranscriptMap, result)
}

func TestTranscriptRepository_FindByUserNRPAndGroupByRegistrationID_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockTranscriptRepository)

	ctx := context.Background()
	userNRP := "5123123123"

	mockRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(map[string][]entity.Transcript{}, errors.New("error"))

	result, err := mockRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}
