package mocks

import (
	"context"
	"monitoring-service/dto"
	"monitoring-service/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockReportScheduleRepository struct {
	mock.Mock
}

func (m *MockReportScheduleRepository) Index(ctx context.Context, tx *gorm.DB) (map[string][]entity.ReportSchedule, error) {
	args := m.Called(ctx, tx)

	return args.Get(0).(map[string][]entity.ReportSchedule), args.Error(1)
}

func (m *MockReportScheduleRepository) Create(ctx context.Context, reportSchedule entity.ReportSchedule, tx *gorm.DB) (entity.ReportSchedule, error) {
	args := m.Called(ctx, reportSchedule, tx)

	return args.Get(0).(entity.ReportSchedule), args.Error(1)
}

func (m *MockReportScheduleRepository) Update(ctx context.Context, id string, reportSchedule entity.ReportSchedule, tx *gorm.DB) error {
	args := m.Called(ctx, id, reportSchedule, tx)

	return args.Error(0)
}

func (m *MockReportScheduleRepository) FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.ReportSchedule, error) {
	args := m.Called(ctx, id, tx)

	return args.Get(0).(entity.ReportSchedule), args.Error(1)
}

func (m *MockReportScheduleRepository) Destroy(ctx context.Context, id string, tx *gorm.DB) error {
	args := m.Called(ctx, id, tx)

	return args.Error(0)
}

func (m *MockReportScheduleRepository) FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.ReportSchedule, error) {
	args := m.Called(ctx, registrationID, tx)

	return args.Get(0).([]entity.ReportSchedule), args.Error(1)
}

func (m *MockReportScheduleRepository) FindByUserID(ctx context.Context, userNRP string, tx *gorm.DB) ([]entity.ReportSchedule, error) {
	args := m.Called(ctx, userNRP, tx)

	return args.Get(0).([]entity.ReportSchedule), args.Error(1)
}

func (m *MockReportScheduleRepository) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, userNRP string, tx *gorm.DB) (map[string][]entity.ReportSchedule, error) {
	args := m.Called(ctx, userNRP, tx)

	return args.Get(0).(map[string][]entity.ReportSchedule), args.Error(1)
}

func (m *MockReportScheduleRepository) FindByAdvisorEmailAndGroupByUserID(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNrp string) (map[string][]entity.ReportSchedule, int64, error) {
	args := m.Called(ctx, advisorEmail, tx, pagReq, userNrp)

	return args.Get(0).(map[string][]entity.ReportSchedule), args.Get(1).(int64), args.Error(2)
}
