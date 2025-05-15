package service_mock

import "github.com/stretchr/testify/mock"

type MockUserManagementService struct {
	mock.Mock
}

func NewMockUserManagementService() *MockUserManagementService {
	return &MockUserManagementService{}
}

func (m *MockUserManagementService) GetUserByID(method string, token string, id string) map[string]interface{} {
	args := m.Called(method, token, id)

	return args.Get(0).(map[string]interface{})
}

func (m *MockUserManagementService) GetUserData(method string, token string) map[string]interface{} {
	args := m.Called(method, token)

	return args.Get(0).(map[string]interface{})
}

func (m *MockUserManagementService) GetUserByFilter(data map[string]interface{}, method string, token string) []map[string]interface{} {
	args := m.Called(data, method, token)

	return args.Get(0).([]map[string]interface{})
}

func (m *MockUserManagementService) GetUserRole(method string, token string) map[string]interface{} {
	args := m.Called(method, token)

	return args.Get(0).(map[string]interface{})
}

func (m *MockUserManagementService) GetDosenDataByEmail(data map[string]interface{}, method string, token string) map[string]interface{} {
	args := m.Called(data, method, token)

	return args.Get(0).(map[string]interface{})
}
