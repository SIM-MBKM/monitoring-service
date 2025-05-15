package service_mock

import (
	"context"
	"mime/multipart"
	"monitoring-service/dto"

	"github.com/stretchr/testify/mock"
)

type MockSyllabusService struct {
	mock.Mock
}

func NewMockSyllabusService() *MockSyllabusService {
	return &MockSyllabusService{}
}

func (m *MockSyllabusService) Index(ctx context.Context) ([]dto.SyllabusResponse, error) {
	args := m.Called(ctx)

	return args.Get(0).([]dto.SyllabusResponse), args.Error(1)
}

func (m *MockSyllabusService) Create(ctx context.Context, syllabus dto.SyllabusRequest, file *multipart.FileHeader, token string) (dto.SyllabusResponse, error) {
	args := m.Called(ctx, syllabus, file, token)

	return args.Get(0).(dto.SyllabusResponse), args.Error(1)
}

func (m *MockSyllabusService) Update(ctx context.Context, id string, syllabus dto.SyllabusRequest) error {
	args := m.Called(ctx, id, syllabus)

	return args.Error(0)
}

func (m *MockSyllabusService) FindByID(ctx context.Context, id string, token string) (dto.SyllabusResponse, error) {
	args := m.Called(ctx, id, token)

	return args.Get(0).(dto.SyllabusResponse), args.Error(1)
}

func (m *MockSyllabusService) Destroy(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *MockSyllabusService) FindByRegistrationID(ctx context.Context, registrationID string) (dto.SyllabusResponse, error) {
	args := m.Called(ctx, registrationID)

	return args.Get(0).(dto.SyllabusResponse), args.Error(1)
}

func (m *MockSyllabusService) FindAllByRegistrationID(ctx context.Context, registrationID string) ([]dto.SyllabusResponse, error) {
	args := m.Called(ctx, registrationID)

	return args.Get(0).([]dto.SyllabusResponse), args.Error(1)
}

func (m *MockSyllabusService) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, token string, pagReq dto.PaginationRequest, filter dto.SyllabusAdvisorFilterRequest) (dto.SyllabusAdvisorResponse, dto.PaginationResponse, error) {
	args := m.Called(ctx, token, pagReq, filter)

	return args.Get(0).(dto.SyllabusAdvisorResponse), args.Get(1).(dto.PaginationResponse), args.Error(2)
}

func (m *MockSyllabusService) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.SyllabusByStudentResponse, error) {
	args := m.Called(ctx, token)

	return args.Get(0).(dto.SyllabusByStudentResponse), args.Error(1)
}
