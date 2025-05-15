package service_mock

import "github.com/stretchr/testify/mock"

type MockRegistrationService struct {
	mock.Mock
}

func NewMockRegistrationService() *MockRegistrationService {
	return &MockRegistrationService{}
}

func (m *MockRegistrationService) GetRegistrationByID(method string, id string, token string) map[string]interface{} {
	args := m.Called(method, id, token)

	return args.Get(0).(map[string]interface{})
}
