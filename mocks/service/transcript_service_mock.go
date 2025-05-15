package service_mock

import (
	"context"
	"mime/multipart"
	"monitoring-service/dto"

	"github.com/stretchr/testify/mock"
)

type MockTranscriptService struct {
	mock.Mock
}

func NewMockTranscriptService() *MockTranscriptService {
	return &MockTranscriptService{}
}

func (m *MockTranscriptService) Index(ctx context.Context) ([]dto.TranscriptResponse, error) {
	args := m.Called(ctx)

	return args.Get(0).([]dto.TranscriptResponse), args.Error(1)
}

func (m *MockTranscriptService) Create(ctx context.Context, transcript dto.TranscriptRequest, file *multipart.FileHeader, token string) (dto.TranscriptResponse, error) {
	args := m.Called(ctx, transcript, file, token)

	return args.Get(0).(dto.TranscriptResponse), args.Error(1)
}

func (m *MockTranscriptService) Update(ctx context.Context, id string, transcript dto.TranscriptRequest) error {
	args := m.Called(ctx, id, transcript)

	return args.Error(0)
}

func (m *MockTranscriptService) FindByID(ctx context.Context, id string, token string) (dto.TranscriptResponse, error) {
	args := m.Called(ctx, id, token)

	return args.Get(0).(dto.TranscriptResponse), args.Error(1)
}

func (m *MockTranscriptService) Destroy(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *MockTranscriptService) FindByRegistrationID(ctx context.Context, registrationID string) (dto.TranscriptResponse, error) {
	args := m.Called(ctx, registrationID)

	return args.Get(0).(dto.TranscriptResponse), args.Error(1)
}

func (m *MockTranscriptService) FindAllByRegistrationID(ctx context.Context, registrationID string) ([]dto.TranscriptResponse, error) {
	args := m.Called(ctx, registrationID)

	return args.Get(0).([]dto.TranscriptResponse), args.Error(1)
}

func (m *MockTranscriptService) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, token string, pagReq dto.PaginationRequest, filter dto.TranscriptAdvisorFilterRequest) (dto.TranscriptAdvisorResponse, dto.PaginationResponse, error) {
	args := m.Called(ctx, token, pagReq, filter)

	return args.Get(0).(dto.TranscriptAdvisorResponse), args.Get(1).(dto.PaginationResponse), args.Error(2)
}

func (m *MockTranscriptService) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.TranscriptByStudentResponse, error) {
	args := m.Called(ctx, token)

	return args.Get(0).(dto.TranscriptByStudentResponse), args.Error(1)
}
