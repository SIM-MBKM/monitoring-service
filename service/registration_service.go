package service

import (
	"fmt"
	"log"

	baseService "github.com/SIM-MBKM/mod-service/src/service"
)

type RegistrationManagementService struct {
	baseService *baseService.Service
}

const (
	GET_REGISTRATION_BY_ID_ENDPOINT = "registration-management/api/v1/registration/%s"
)

func NewRegistrationManagementService(baseURI string, asyncURIs []string) *RegistrationManagementService {
	return &RegistrationManagementService{
		baseService: baseService.NewService(baseURI, asyncURIs),
	}
}

// create function to get user by id
func (s *RegistrationManagementService) GetRegistrationByID(method string, id string, token string) map[string]interface{} {
	endpoint := fmt.Sprintf(GET_REGISTRATION_BY_ID_ENDPOINT, id)
	res, err := s.baseService.Request(method, endpoint, nil, token)
	if err != nil {
		return nil
	}

	registration, ok := res["data"].(map[string]interface{})
	if !ok {
		return nil
	}

	var registrationData map[string]interface{}
	registrationData = map[string]interface{}{
		"id":                     registration["id"],
		"user_id":                registration["user_id"],
		"user_nrp":               registration["user_nrp"],
		"user_name":              registration["user_name"],
		"academic_advisor":       registration["academic_advisor"],
		"academic_advisor_email": registration["academic_advisor_email"],
		"activity_name":          registration["activity_name"],
	}

	log.Println("REGISTRATION DATA: ", registrationData)
	return registrationData
}
