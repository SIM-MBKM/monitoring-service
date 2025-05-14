package repository_test

import (
	"context"
	"errors"
	"monitoring-service/entity"
	repository_mock "monitoring-service/mocks/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func createMockReport() entity.Report {
	now := time.Now()

	return entity.Report{
		ID:                    uuid.New(),
		ReportScheduleID:      uuid.NewString(),
		Title:                 "Test Report",
		Content:               "Test Content",
		ReportType:            "test",
		FileStorageID:         uuid.NewString(),
		Feedback:              "Test Feedback",
		AcademicAdvisorStatus: "APPROVED",
		BaseModel: entity.BaseModel{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}
}

func TestReportRepository_Index(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	reports := []entity.Report{createMockReport(), createMockReport()}
	mockRepo.On("Index", ctx, mock.Anything).Return(reports, nil)

	result, err := mockRepo.Index(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, reports, result)
}

func TestReportRepository_Index_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()

	mockRepo.On("Index", ctx, mock.Anything).Return([]entity.Report{}, errors.New("error"))

	result, err := mockRepo.Index(ctx, nil)

	assert.Error(t, err)
	assert.Equal(t, []entity.Report{}, result)
}

func TestReportRepository_Index_NotFound(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()

	mockRepo.On("Index", ctx, mock.Anything).Return([]entity.Report{}, gorm.ErrRecordNotFound)

	result, err := mockRepo.Index(ctx, nil)

	assert.Error(t, err)
	assert.Equal(t, []entity.Report{}, result)
}

func TestReportRepository_Create(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	mockReport := createMockReport()
	mockRepo.On("Create", ctx, mock.AnythingOfType("entity.Report"), mock.Anything).Return(mockReport, nil)

	result, err := mockRepo.Create(ctx, mockReport, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockReport.ID, result.ID)
}

func TestReportRepository_Create_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	mockReport := createMockReport()
	mockRepo.On("Create", ctx, mock.AnythingOfType("entity.Report"), mock.Anything).Return(entity.Report{}, errors.New("error"))

	result, err := mockRepo.Create(ctx, mockReport, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.Report{}, result)
}

func TestReportRepository_Update(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	mockReport := createMockReport()
	id := mockReport.ID.String()
	mockRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Report"), mock.Anything).Return(nil)

	err := mockRepo.Update(ctx, id, mockReport, nil)

	assert.NoError(t, err)
}

func TestReportRepository_Update_ErrorDB(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	mockReport := createMockReport()
	id := mockReport.ID.String()
	mockRepo.On("Update", ctx, id, mock.AnythingOfType("entity.Report"), mock.Anything).Return(errors.New("Database Error"))

	err := mockRepo.Update(ctx, id, mockReport, nil)

	assert.Error(t, err)
}

func TestReportRepository_FindByID(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	mockReport := createMockReport()
	id := mockReport.ID.String()
	mockRepo.On("FindByID", ctx, id, mock.Anything).Return(mockReport, nil)

	result, err := mockRepo.FindByID(ctx, id, nil)

	assert.NoError(t, err)
	assert.Equal(t, mockReport.ID, result.ID)
}

func TestReportRepository_FindByID_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	id := createMockReport().ID.String()
	mockRepo.On("FindByID", ctx, id, mock.Anything).Return(entity.Report{}, errors.New("error"))

	result, err := mockRepo.FindByID(ctx, id, nil)

	assert.Error(t, err)
	assert.Equal(t, entity.Report{}, result)
}

func TestReportRepository_Destroy(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	id := createMockReport().ID.String()
	mockRepo.On("Destroy", ctx, id, mock.Anything).Return(nil)

	err := mockRepo.Destroy(ctx, id, nil)

	assert.NoError(t, err)
}

func TestReportRepository_Destroy_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	id := createMockReport().ID.String()
	mockRepo.On("Destroy", ctx, id, mock.Anything).Return(errors.New("error"))

	err := mockRepo.Destroy(ctx, id, nil)

	assert.Error(t, err)
}

func TestReportRepository_FindByReportScheduleID(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	reportScheduleID := createMockReport().ReportScheduleID
	reports := []entity.Report{createMockReport()}
	mockRepo.On("FindByReportScheduleID", ctx, reportScheduleID, mock.Anything).Return(reports, nil)

	result, err := mockRepo.FindByReportScheduleID(ctx, reportScheduleID, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
}

func TestReportRepository_FindByReportScheduleID_NotFound(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	reportScheduleID := createMockReport().ReportScheduleID
	mockRepo.On("FindByReportScheduleID", ctx, reportScheduleID, mock.Anything).Return([]entity.Report{}, gorm.ErrRecordNotFound)

	result, err := mockRepo.FindByReportScheduleID(ctx, reportScheduleID, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestReportRepository_FindByReportScheduleID_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	reportScheduleID := createMockReport().ReportScheduleID
	mockRepo.On("FindByReportScheduleID", ctx, reportScheduleID, mock.Anything).Return([]entity.Report{}, errors.New("error"))

	result, err := mockRepo.FindByReportScheduleID(ctx, reportScheduleID, nil)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestReportRepository_Approval(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	mockReport := createMockReport()
	id := mockReport.ID.String()
	mockRepo.On("Approval", ctx, id, mock.AnythingOfType("entity.Report"), mock.Anything).Return(nil)

	err := mockRepo.Approval(ctx, id, mockReport, nil)

	assert.NoError(t, err)
}

func TestReportRepository_Approval_NotFound(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	mockReport := createMockReport()
	id := mockReport.ID.String()
	mockRepo.On("Approval", ctx, id, mock.AnythingOfType("entity.Report"), mock.Anything).Return(gorm.ErrRecordNotFound)

	err := mockRepo.Approval(ctx, id, mockReport, nil)

	assert.Error(t, err)
}

func TestReportRepository_Approval_Error(t *testing.T) {
	mockRepo := new(repository_mock.MockReportRepository)

	ctx := context.Background()
	mockReport := createMockReport()
	id := mockReport.ID.String()
	mockRepo.On("Approval", ctx, id, mock.AnythingOfType("entity.Report"), mock.Anything).Return(errors.New("error"))

	err := mockRepo.Approval(ctx, id, mockReport, nil)

	assert.Error(t, err)
}
