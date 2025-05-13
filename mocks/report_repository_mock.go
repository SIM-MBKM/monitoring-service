package mocks

import (
	"context"
	"monitoring-service/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) Index(ctx context.Context, tx *gorm.DB) ([]entity.Report, error) {
	args := m.Called(ctx, tx)

	return args.Get(0).([]entity.Report), args.Error(1)
}

func (m *MockReportRepository) Create(ctx context.Context, report entity.Report, tx *gorm.DB) (entity.Report, error) {
	args := m.Called(ctx, report, tx)

	return args.Get(0).(entity.Report), args.Error(1)
}

func (m *MockReportRepository) Update(ctx context.Context, id string, report entity.Report, tx *gorm.DB) error {
	args := m.Called(ctx, id, report, tx)

	return args.Error(0)
}

func (m *MockReportRepository) FindByReportScheduleID(ctx context.Context, reportScheduleID string, tx *gorm.DB) ([]entity.Report, error) {
	args := m.Called(ctx, reportScheduleID, tx)

	return args.Get(0).([]entity.Report), args.Error(1)
}

func (m *MockReportRepository) Approval(ctx context.Context, id string, report entity.Report, tx *gorm.DB) error {
	args := m.Called(ctx, id, report, tx)

	return args.Error(0)
}
