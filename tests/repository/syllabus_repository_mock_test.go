package repository_test

import (
	"context"
	"errors"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"monitoring-service/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockSyllabus() entity.Syllabus {
	now := time.Now()

	return entity.Syllabus{
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

func TestSyllabusRepository_Index(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	mockSyllabuses := []entity.Syllabus{createMockSyllabus(), createMockSyllabus()}

	mockRepo.On("Index", ctx, mock.Anything).Return(mockSyllabuses, nil)

	result, err := mockRepo.Index(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSyllabuses, result)
}

func TestSyllabusRepository_Index_Error(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()

	mockRepo.On("Index", ctx, mock.Anything).Return([]entity.Syllabus{}, errors.New("error"))

	result, err := mockRepo.Index(ctx, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestSyllabusRepository_Create(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	mockSyllabus := createMockSyllabus()

	mockRepo.On("Create", ctx, mock.AnythingOfType("entity.Syllabus"), mock.Anything).Return(mockSyllabus, nil)

	result, err := mockRepo.Create(ctx, mockSyllabus, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSyllabus.ID, result.ID)
}

func TestSyllabusRepository_Create_Error(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	mockSyllabus := createMockSyllabus()

	mockRepo.On("Create", ctx, mock.AnythingOfType("entity.Syllabus"), mock.Anything).Return(entity.Syllabus{}, errors.New("error"))

	result, err := mockRepo.Create(ctx, mockSyllabus, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.Syllabus{}, result)
}

func TestSyllabusRepository_Update(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	mockSyllabus := createMockSyllabus()
	id := mockSyllabus.ID.String()

	mockRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Syllabus"), mock.Anything).Return(nil)

	err := mockRepo.Update(ctx, id, mockSyllabus, nil)

	assert.NoError(t, err)
}

func TestSyllabusRepository_Update_Error(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	mockSyllabus := createMockSyllabus()
	id := mockSyllabus.ID.String()

	mockRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Syllabus"), mock.Anything).Return(errors.New("error"))

	err := mockRepo.Update(ctx, id, mockSyllabus, nil)

	assert.Error(t, err)
}

func TestSyllabusRepository_FindByID(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	mockSyllabus := createMockSyllabus()
	id := mockSyllabus.ID.String()

	mockRepo.On("FindByID", ctx, id, mock.Anything).Return(mockSyllabus, nil)

	result, err := mockRepo.FindByID(ctx, id, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSyllabus.ID, result.ID)
}

func TestSyllabusRepository_FindByID_Error(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	id := uuid.New().String()

	mockRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.Syllabus{}, errors.New("error"))

	result, err := mockRepo.FindByID(ctx, id, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.Syllabus{}, result)
}

func TestSyllabusRepository_Destroy(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	id := uuid.New().String()

	mockRepo.On("Destroy", ctx, id, mock.Anything).Return(nil)

	err := mockRepo.Destroy(ctx, id, nil)

	assert.NoError(t, err)
}

func TestSyllabusRepository_Destroy_Error(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	id := uuid.New().String()

	mockRepo.On("Destroy", ctx, id, mock.Anything).Return(errors.New("error"))

	err := mockRepo.Destroy(ctx, id, nil)

	assert.Error(t, err)
}

func TestSyllabusRepository_FindByRegistrationID(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()
	mockSyllabus := createMockSyllabus()

	mockRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).Return(mockSyllabus, nil)

	result, err := mockRepo.FindByRegistrationID(ctx, registrationID, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSyllabus.ID, result.ID)
}

func TestSyllabusRepository_FindByRegistrationID_Error(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()

	mockRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).Return(entity.Syllabus{}, errors.New("error"))

	result, err := mockRepo.FindByRegistrationID(ctx, registrationID, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.Syllabus{}, result)
}

func TestSyllabusRepository_FindAllByRegistrationID(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()
	mockSyllabuses := []entity.Syllabus{createMockSyllabus(), createMockSyllabus()}

	mockRepo.On("FindAllByRegistrationID", ctx, registrationID, mock.Anything).Return(mockSyllabuses, nil)

	result, err := mockRepo.FindAllByRegistrationID(ctx, registrationID, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSyllabuses, result)
}

func TestSyllabusRepository_FindAllByRegistrationID_Error(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()

	mockRepo.On("FindAllByRegistrationID", ctx, registrationID, mock.Anything).Return([]entity.Syllabus{}, errors.New("error"))

	result, err := mockRepo.FindAllByRegistrationID(ctx, registrationID, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestSyllabusRepository_FindByAdvisorEmailAndGroupByUserNRP(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	advisorEmail := "advisor@example.com"
	pagReq := &dto.PaginationRequest{Limit: 10, Offset: 0}
	userNRPFilter := ""
	mockSyllabusMap := map[string]entity.Syllabus{
		"5123123123": createMockSyllabus(),
	}
	totalCount := int64(1)

	mockRepo.On("FindByAdvisorEmailAndGroupByUserNRP", ctx, advisorEmail, mock.Anything, pagReq, userNRPFilter).Return(mockSyllabusMap, totalCount, nil)

	result, count, err := mockRepo.FindByAdvisorEmailAndGroupByUserNRP(ctx, advisorEmail, nil, pagReq, userNRPFilter)

	assert.NoError(t, err)
	assert.Equal(t, mockSyllabusMap, result)
	assert.Equal(t, totalCount, count)
}

func TestSyllabusRepository_FindByAdvisorEmailAndGroupByUserNRP_Error(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	advisorEmail := "advisor@example.com"
	pagReq := &dto.PaginationRequest{Limit: 10, Offset: 0}
	userNRPFilter := ""

	mockRepo.On("FindByAdvisorEmailAndGroupByUserNRP", ctx, advisorEmail, mock.Anything, pagReq, userNRPFilter).Return(map[string]entity.Syllabus{}, int64(0), errors.New("error"))

	result, count, err := mockRepo.FindByAdvisorEmailAndGroupByUserNRP(ctx, advisorEmail, nil, pagReq, userNRPFilter)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), count)
}

func TestSyllabusRepository_FindByUserNRPAndGroupByRegistrationID(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	userNRP := "5123123123"
	mockSyllabusMap := map[string][]entity.Syllabus{
		"reg-id-1": {createMockSyllabus(), createMockSyllabus()},
	}

	mockRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(mockSyllabusMap, nil)

	result, err := mockRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSyllabusMap, result)
}

func TestSyllabusRepository_FindByUserNRPAndGroupByRegistrationID_Error(t *testing.T) {
	mockRepo := new(mocks.MockSyllabusRepository)

	ctx := context.Background()
	userNRP := "5123123123"

	mockRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(map[string][]entity.Syllabus{}, errors.New("error"))

	result, err := mockRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}
