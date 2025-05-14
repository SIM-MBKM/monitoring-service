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

func createMockReportSchedule() entity.ReportSchedule {
	now := time.Now()

	return entity.ReportSchedule{
		ID:                   uuid.New(),
		UserNRP:              "5123123123",
		RegistrationID:       uuid.NewString(),
		StartDate:            &now,
		EndDate:              &now,
		AcademicAdvisorEmail: "advisor@example.com",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}
}

func TestReportScheduleRepository_Index(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	mockSchedules := map[string][]entity.ReportSchedule{
		"5123123123": {createMockReportSchedule(), createMockReportSchedule()},
	}

	mockRepo.On("Index", ctx, mock.Anything).Return(mockSchedules, nil)

	result, err := mockRepo.Index(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSchedules, result)
}

func TestReportScheduleRepository_Index_Error(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()

	mockRepo.On("Index", ctx, mock.Anything).Return(map[string][]entity.ReportSchedule{}, errors.New("error"))

	result, err := mockRepo.Index(ctx, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestReportScheduleRepository_Create(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	mockSchedule := createMockReportSchedule()

	mockRepo.On("Create", ctx, mock.AnythingOfType("entity.ReportSchedule"), mock.Anything).Return(mockSchedule, nil)

	result, err := mockRepo.Create(ctx, mockSchedule, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSchedule.ID, result.ID)
}

func TestReportScheduleRepository_Create_Error(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	mockSchedule := createMockReportSchedule()

	mockRepo.On("Create", ctx, mock.AnythingOfType("entity.ReportSchedule"), mock.Anything).Return(entity.ReportSchedule{}, errors.New("error"))

	result, err := mockRepo.Create(ctx, mockSchedule, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.ReportSchedule{}, result)
}

func TestReportScheduleRepository_Update(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	mockSchedule := createMockReportSchedule()
	id := mockSchedule.ID.String()

	mockRepo.On("Update", ctx, id, mock.AnythingOfType("entity.ReportSchedule"), mock.Anything).Return(nil)

	err := mockRepo.Update(ctx, id, mockSchedule, nil)

	assert.NoError(t, err)
}

func TestReportScheduleRepository_Update_Error(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	mockSchedule := createMockReportSchedule()
	id := mockSchedule.ID.String()

	mockRepo.On("Update", ctx, id, mock.AnythingOfType("entity.ReportSchedule"), mock.Anything).Return(errors.New("error"))

	err := mockRepo.Update(ctx, id, mockSchedule, nil)

	assert.Error(t, err)
}

func TestReportScheduleRepository_FindByID(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	mockSchedule := createMockReportSchedule()
	id := mockSchedule.ID.String()

	mockRepo.On("FindByID", ctx, id, mock.Anything).Return(mockSchedule, nil)

	result, err := mockRepo.FindByID(ctx, id, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSchedule.ID, result.ID)
}

func TestReportScheduleRepository_FindByID_Error(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	id := uuid.New().String()

	mockRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.ReportSchedule{}, errors.New("error"))

	result, err := mockRepo.FindByID(ctx, id, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.ReportSchedule{}, result)
}

func TestReportScheduleRepository_Destroy(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	id := uuid.New().String()

	mockRepo.On("Destroy", ctx, id, mock.Anything).Return(nil)

	err := mockRepo.Destroy(ctx, id, nil)

	assert.NoError(t, err)
}

func TestReportScheduleRepository_Destroy_Error(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	id := uuid.New().String()

	mockRepo.On("Destroy", ctx, id, mock.Anything).Return(errors.New("error"))

	err := mockRepo.Destroy(ctx, id, nil)

	assert.Error(t, err)
}

func TestReportScheduleRepository_FindByRegistrationID(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()
	mockSchedules := []entity.ReportSchedule{createMockReportSchedule(), createMockReportSchedule()}

	mockRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).Return(mockSchedules, nil)

	result, err := mockRepo.FindByRegistrationID(ctx, registrationID, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSchedules, result)
}

func TestReportScheduleRepository_FindByRegistrationID_Error(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	registrationID := uuid.New().String()

	mockRepo.On("FindByRegistrationID", ctx, registrationID, mock.Anything).Return([]entity.ReportSchedule{}, errors.New("error"))

	result, err := mockRepo.FindByRegistrationID(ctx, registrationID, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestReportScheduleRepository_FindByUserID(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	userNRP := "5123123123"
	mockSchedules := []entity.ReportSchedule{createMockReportSchedule(), createMockReportSchedule()}

	mockRepo.On("FindByUserID", ctx, userNRP, mock.Anything).Return(mockSchedules, nil)

	result, err := mockRepo.FindByUserID(ctx, userNRP, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSchedules, result)
}

func TestReportScheduleRepository_FindByUserID_Error(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	userNRP := "5123123123"

	mockRepo.On("FindByUserID", ctx, userNRP, mock.Anything).Return([]entity.ReportSchedule{}, errors.New("error"))

	result, err := mockRepo.FindByUserID(ctx, userNRP, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestReportScheduleRepository_FindByUserNRPAndGroupByRegistrationID(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	userNRP := "5123123123"
	mockSchedules := map[string][]entity.ReportSchedule{
		"reg-id-1": {createMockReportSchedule(), createMockReportSchedule()},
	}

	mockRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(mockSchedules, nil)

	result, err := mockRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockSchedules, result)
}

func TestReportScheduleRepository_FindByUserNRPAndGroupByRegistrationID_Error(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	userNRP := "5123123123"

	mockRepo.On("FindByUserNRPAndGroupByRegistrationID", ctx, userNRP, mock.Anything).Return(map[string][]entity.ReportSchedule{}, errors.New("error"))

	result, err := mockRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestReportScheduleRepository_FindByAdvisorEmailAndGroupByUserID(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	advisorEmail := "advisor@example.com"
	pagReq := &dto.PaginationRequest{Limit: 10, Offset: 0}
	userNRP := ""
	mockSchedules := map[string][]entity.ReportSchedule{
		"5123123123": {createMockReportSchedule(), createMockReportSchedule()},
	}
	totalCount := int64(1)

	mockRepo.On("FindByAdvisorEmailAndGroupByUserID", ctx, advisorEmail, mock.Anything, pagReq, userNRP).Return(mockSchedules, totalCount, nil)

	result, count, err := mockRepo.FindByAdvisorEmailAndGroupByUserID(ctx, advisorEmail, nil, pagReq, userNRP)

	assert.NoError(t, err)
	assert.Equal(t, mockSchedules, result)
	assert.Equal(t, totalCount, count)
}

func TestReportScheduleRepository_FindByAdvisorEmailAndGroupByUserID_Error(t *testing.T) {
	mockRepo := new(mocks.MockReportScheduleRepository)

	ctx := context.Background()
	advisorEmail := "advisor@example.com"
	pagReq := &dto.PaginationRequest{Limit: 10, Offset: 0}
	userNRP := ""

	mockRepo.On("FindByAdvisorEmailAndGroupByUserID", ctx, advisorEmail, mock.Anything, pagReq, userNRP).Return(map[string][]entity.ReportSchedule{}, int64(0), errors.New("error"))

	result, count, err := mockRepo.FindByAdvisorEmailAndGroupByUserID(ctx, advisorEmail, nil, pagReq, userNRP)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), count)
}
