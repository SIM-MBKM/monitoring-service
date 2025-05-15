package service_mock

import (
	"context"
	"mime/multipart"
	"monitoring-service/dto"

	"github.com/stretchr/testify/mock"
)

type MockReportService struct {
	mock.Mock
}

func NewMockReportService() *MockReportService {
	return &MockReportService{}
}

func (m *MockReportService) Index(ctx context.Context) ([]dto.ReportResponse, error) {
	args := m.Called(ctx)

	return args.Get(0).([]dto.ReportResponse), args.Error(1)
}

func (m *MockReportService) Create(ctx context.Context, report dto.ReportRequest, file *multipart.FileHeader, token string) (dto.ReportResponse, error) {
	args := m.Called(ctx, report, file, token)

	return args.Get(0).(dto.ReportResponse), args.Error(1)
}

func (m *MockReportService) Update(ctx context.Context, id string, report dto.ReportRequest) (dto.ReportResponse, error) {
	args := m.Called(ctx, id, report)

	return args.Get(0).(dto.ReportResponse), args.Error(1)
}

func (m *MockReportService) FindByID(ctx context.Context, id string, token string) (dto.ReportResponse, error) {
	args := m.Called(ctx, id, token)

	return args.Get(0).(dto.ReportResponse), args.Error(1)
}

func (m *MockReportService) Destroy(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *MockReportService) FindByReportScheduleID(ctx context.Context, reportScheduleID string) ([]dto.ReportResponse, error) {
	args := m.Called(ctx, reportScheduleID)

	return args.Get(0).([]dto.ReportResponse), args.Error(1)
}

func (m *MockReportService) Approval(ctx context.Context, token string, report dto.ReportApprovalRequest) error {
	args := m.Called(ctx, token, report)

	return args.Error(0)
}
