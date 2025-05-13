package mocks

import (
	"context"
	"monitoring-service/dto"
	"monitoring-service/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockSyllabusRepository struct {
	mock.Mock
}

func (m *MockSyllabusRepository) Index(ctx context.Context, tx *gorm.DB) ([]entity.Syllabus, error) {
	args := m.Called(ctx, tx)

	return args.Get(0).([]entity.Syllabus), args.Error(1)
}

func (m *MockSyllabusRepository) Create(ctx context.Context, syllabus entity.Syllabus, tx *gorm.DB) (entity.Syllabus, error) {
	args := m.Called(ctx, syllabus, tx)

	return args.Get(0).(entity.Syllabus), args.Error(1)
}

func (m *MockSyllabusRepository) Update(ctx context.Context, id string, syllabus entity.Syllabus, tx *gorm.DB) error {
	args := m.Called(ctx, id, syllabus, tx)

	return args.Error(0)
}

func (m *MockSyllabusRepository) FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.Syllabus, error) {
	args := m.Called(ctx, id, tx)

	return args.Get(0).(entity.Syllabus), args.Error(1)
}

func (m *MockSyllabusRepository) Destroy(ctx context.Context, id string, tx *gorm.DB) error {
	args := m.Called(ctx, id, tx)

	return args.Error(0)
}

func (m *MockSyllabusRepository) FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) (entity.Syllabus, error) {
	args := m.Called(ctx, registrationID, tx)

	return args.Get(0).(entity.Syllabus), args.Error(1)
}

func (m *MockSyllabusRepository) FindAllByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.Syllabus, error) {
	args := m.Called(ctx, registrationID, tx)

	return args.Get(0).([]entity.Syllabus), args.Error(1)
}

func (m *MockSyllabusRepository) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNRPFilter string) (map[string]entity.Syllabus, int64, error) {
	args := m.Called(ctx, advisorEmail, tx, pagReq, userNRPFilter)

	return args.Get(0).(map[string]entity.Syllabus), args.Get(1).(int64), args.Error(2)
}

func (m *MockSyllabusRepository) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, userNRP string, tx *gorm.DB) (map[string][]entity.Syllabus, error) {
	args := m.Called(ctx, userNRP, tx)

	return args.Get(0).(map[string][]entity.Syllabus), args.Error(1)
}
