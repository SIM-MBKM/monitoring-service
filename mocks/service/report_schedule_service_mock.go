package service_mock

import (
	"context"
	"monitoring-service/dto"

	"github.com/stretchr/testify/mock"
)

type MockReportScheduleService struct {
	mock.Mock
}

func NewMockReportScheduleService() *MockReportScheduleService {
	return &MockReportScheduleService{}
}

func (m *MockReportScheduleService) Index(ctx context.Context, token string) (dto.ReportScheduleByAdvisorResponse, error) {
	args := m.Called(ctx, token)

	return args.Get(0).(dto.ReportScheduleByAdvisorResponse), args.Error(1)
}

func (m *MockReportScheduleService) Create(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (dto.ReportScheduleResponse, error) {
	args := m.Called(ctx, reportSchedule, token)

	return args.Get(0).(dto.ReportScheduleResponse), args.Error(1)
}

func (m *MockReportScheduleService) Update(ctx context.Context, id string, subject dto.ReportScheduleRequest, token string) error {
	args := m.Called(ctx, id, subject, token)

	return args.Error(0)
}

func (m *MockReportScheduleService) Destroy(ctx context.Context, id string, token string) error {
	args := m.Called(ctx, id, token)

	return args.Error(0)
}

func (m *MockReportScheduleService) FindByRegistrationID(ctx context.Context, registrationID string) ([]dto.ReportScheduleResponse, error) {
	args := m.Called(ctx, registrationID)

	return args.Get(0).([]dto.ReportScheduleResponse), args.Error(1)
}

func (m *MockReportScheduleService) ReportScheduleAccess(ctx context.Context, reportSchedule dto.ReportScheduleRequest, token string) (bool, error) {
	args := m.Called(ctx, reportSchedule, token)

	return args.Get(0).(bool), args.Error(1)
}

func (m *MockReportScheduleService) FindByUserID(ctx context.Context, token string) ([]dto.ReportScheduleResponse, error) {
	args := m.Called(ctx, token)

	return args.Get(0).([]dto.ReportScheduleResponse), args.Error(1)
}

func (m *MockReportScheduleService) FindByAdvisorEmail(ctx context.Context, token string, pagReq dto.PaginationRequest, reportScheduleRequest dto.ReportScheduleAdvisorRequest) (dto.ReportScheduleByAdvisorResponse, dto.PaginationResponse, error) {
	args := m.Called(ctx, token, pagReq, reportScheduleRequest)

	return args.Get(0).(dto.ReportScheduleByAdvisorResponse), args.Get(1).(dto.PaginationResponse), args.Error(2)
}

func (m *MockReportScheduleService) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.ReportScheduleByStudentResponse, error) {
	args := m.Called(ctx, token)

	return args.Get(0).(dto.ReportScheduleByStudentResponse), args.Error(1)
}
