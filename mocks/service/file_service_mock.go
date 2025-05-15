package service_mock

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MockFileService struct {
	mock.Mock
}

func NewMockFileService() *MockFileService {
	return &MockFileService{}
}

func (m *MockFileService) Upload(file *multipart.FileHeader, bucket string, folder string, filename string) (map[string]interface{}, error) {
	args := m.Called(file, bucket, folder, filename)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
