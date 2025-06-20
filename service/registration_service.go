package service

import (
	"fmt"
	"log"
	"strings"
	"sync"

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
	// split token
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 {
		return nil
	}

	token = tokenParts[1]

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
		"approval_status":        registration["approval_status"],
	}

	return registrationData
}

// Tambahkan method baru untuk batch processing
func (s *RegistrationManagementService) GetRegistrationsByIDs(method string, ids []string, token string) (map[string]map[string]interface{}, error) {
	// Split token
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}
	token = tokenParts[1]

	// Create result map
	results := make(map[string]map[string]interface{})
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Limit concurrent requests (adjust based on your API rate limits)
	maxConcurrent := 10
	semaphore := make(chan struct{}, maxConcurrent)

	// Error collection
	var errors []error
	var errorMu sync.Mutex

	for _, id := range ids {
		wg.Add(1)
		go func(registrationID string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Call individual API
			endpoint := fmt.Sprintf(GET_REGISTRATION_BY_ID_ENDPOINT, registrationID)
			res, err := s.baseService.Request(method, endpoint, nil, token)

			if err != nil {
				errorMu.Lock()
				errors = append(errors, fmt.Errorf("failed to get registration %s: %v", registrationID, err))
				errorMu.Unlock()
				return
			}

			registration, ok := res["data"].(map[string]interface{})
			if !ok {
				errorMu.Lock()
				errors = append(errors, fmt.Errorf("invalid response format for registration %s", registrationID))
				errorMu.Unlock()
				return
			}

			registrationData := map[string]interface{}{
				"id":                     registration["id"],
				"user_id":                registration["user_id"],
				"user_nrp":               registration["user_nrp"],
				"user_name":              registration["user_name"],
				"academic_advisor":       registration["academic_advisor"],
				"academic_advisor_email": registration["academic_advisor_email"],
				"activity_name":          registration["activity_name"],
				"approval_status":        registration["approval_status"],
			}

			// Store result safely
			mu.Lock()
			results[registrationID] = registrationData
			mu.Unlock()

		}(id)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check if there were any errors
	if len(errors) > 0 {
		// Log errors but don't fail the entire request
		for _, err := range errors {
			log.Printf("Registration API Error: %v", err)
		}
	}

	return results, nil
}
