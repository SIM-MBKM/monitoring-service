package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockBaseRepository struct {
	mock.Mock
}

func (m *MockBaseRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	args := m.Called(ctx)
	return args.Get(0).(*gorm.DB), args.Error(1)
}

func (m *MockBaseRepository) CommitTx(ctx context.Context, tx *gorm.DB) (*gorm.DB, error) {
	args := m.Called(ctx, tx)

	return args.Get(0).(*gorm.DB), args.Error(1)
}

func (m *MockBaseRepository) RollbackTx(ctx context.Context, tx *gorm.DB) {
	m.Called(ctx, tx)
}
