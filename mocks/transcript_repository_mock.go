package mocks

import (
	"context"
	"monitoring-service/dto"
	"monitoring-service/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockTranscriptRepository struct {
	mock.Mock
}

func (m *MockTranscriptRepository) Index(ctx context.Context, tx *gorm.DB) ([]entity.Transcript, error) {
	args := m.Called(ctx, tx)

	return args.Get(0).([]entity.Transcript), args.Error(1)
}

func (m *MockTranscriptRepository) Create(ctx context.Context, transcript entity.Transcript, tx *gorm.DB) (entity.Transcript, error) {
	args := m.Called(ctx, transcript, tx)

	return args.Get(0).(entity.Transcript), args.Error(1)
}

func (m *MockTranscriptRepository) Update(ctx context.Context, id string, transcript entity.Transcript, tx *gorm.DB) error {
	args := m.Called(ctx, id, transcript, tx)

	return args.Error(0)
}

func (m *MockTranscriptRepository) FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.Transcript, error) {
	args := m.Called(ctx, id, tx)

	return args.Get(0).(entity.Transcript), args.Error(1)
}

func (m *MockTranscriptRepository) Destroy(ctx context.Context, id string, tx *gorm.DB) error {
	args := m.Called(ctx, id, tx)

	return args.Error(0)
}

func (m *MockTranscriptRepository) FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) (entity.Transcript, error) {
	args := m.Called(ctx, registrationID, tx)

	return args.Get(0).(entity.Transcript), args.Error(1)
}

func (m *MockTranscriptRepository) FindAllByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.Transcript, error) {
	args := m.Called(ctx, registrationID, tx)

	return args.Get(0).([]entity.Transcript), args.Error(1)
}

func (m *MockTranscriptRepository) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNRPFilter string) (map[string]entity.Transcript, int64, error) {
	args := m.Called(ctx, advisorEmail, tx, pagReq, userNRPFilter)

	return args.Get(0).(map[string]entity.Transcript), args.Get(1).(int64), args.Error(2)
}

func (m *MockTranscriptRepository) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, userNRP string, tx *gorm.DB) (map[string][]entity.Transcript, error) {
	args := m.Called(ctx, userNRP, tx)

	return args.Get(0).(map[string][]entity.Transcript), args.Error(1)
}
