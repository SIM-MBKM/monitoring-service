package service

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"monitoring-service/helper"
	"monitoring-service/repository"
	"reflect"
	"time"

	storageService "github.com/SIM-MBKM/filestorage/storage"
	"github.com/google/uuid"
)

type syllabusService struct {
	syllabusRepo          repository.SyllabusRepository
	fileService           *FileService
	userManagementService *UserManagementService
	registrationService   *RegistrationManagementService
}

type SyllabusService interface {
	Index(ctx context.Context) ([]dto.SyllabusResponse, error)
	Create(ctx context.Context, syllabus dto.SyllabusRequest, file *multipart.FileHeader, token string) (dto.SyllabusResponse, error)
	Update(ctx context.Context, id string, syllabus dto.SyllabusRequest) error
	FindByID(ctx context.Context, id string, token string) (dto.SyllabusResponse, error)
	Destroy(ctx context.Context, id string) error
	FindByRegistrationID(ctx context.Context, registrationID string) (dto.SyllabusResponse, error)
	FindAllByRegistrationID(ctx context.Context, registrationID string) ([]dto.SyllabusResponse, error)
	FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, token string, pagReq dto.PaginationRequest, filter dto.SyllabusAdvisorFilterRequest) (dto.SyllabusAdvisorResponse, dto.PaginationResponse, error)
	FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.SyllabusByStudentResponse, error)
}

func NewSyllabusService(
	syllabusRepo repository.SyllabusRepository,
	userManagementBaseURI string,
	registrationBaseURI string,
	asyncURIs []string,
	config *storageService.Config,
	tokenManager *storageService.CacheTokenManager,
) SyllabusService {
	return &syllabusService{
		syllabusRepo:          syllabusRepo,
		fileService:           NewFileService(config, tokenManager),
		userManagementService: NewUserManagementService(userManagementBaseURI, asyncURIs),
		registrationService:   NewRegistrationManagementService(registrationBaseURI, asyncURIs),
	}
}

func (s *syllabusService) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, token string, pagReq dto.PaginationRequest, filter dto.SyllabusAdvisorFilterRequest) (dto.SyllabusAdvisorResponse, dto.PaginationResponse, error) {
	user := s.userManagementService.GetUserData("GET", token)
	advisorEmail, ok := user["email"].(string)
	if !ok {
		return dto.SyllabusAdvisorResponse{}, dto.PaginationResponse{}, errors.New("advisor email not found")
	}

	syllabuses, totalCount, err := s.syllabusRepo.FindByAdvisorEmailAndGroupByUserNRP(ctx, advisorEmail, nil, &pagReq, filter.UserNRP)
	if err != nil {
		return dto.SyllabusAdvisorResponse{}, dto.PaginationResponse{}, err
	}

	syllabusResponses := make(map[string][]dto.SyllabusResponse)
	for userNRP, syllabus := range syllabuses {
		// Get registration details to add activity name
		registration := s.registrationService.GetRegistrationByID("GET", syllabus.RegistrationID, token)

		var activityName string
		if registration != nil {
			activityNameVal, ok := registration["activity_name"].(string)
			if ok {
				activityName = activityNameVal
			}
		}

		syllabusResponses[userNRP] = append(syllabusResponses[userNRP], dto.SyllabusResponse{
			ID:                   syllabus.ID.String(),
			UserID:               syllabus.UserID,
			UserNRP:              syllabus.UserNRP,
			ActivityName:         activityName,
			AcademicAdvisorID:    syllabus.AcademicAdvisorID,
			AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
			RegistrationID:       syllabus.RegistrationID,
			Title:                syllabus.Title,
			FileStorageID:        syllabus.FileStorageID,
		})
	}

	// Generate pagination metadata
	paginationResponse := helper.MetaDataPagination(totalCount, pagReq)

	return dto.SyllabusAdvisorResponse{
		Syllabuses: syllabusResponses,
	}, paginationResponse, nil
}

// Index retrieves all syllabuses
func (s *syllabusService) Index(ctx context.Context) ([]dto.SyllabusResponse, error) {
	syllabuses, err := s.syllabusRepo.Index(ctx, nil)
	if err != nil {
		return nil, err
	}

	var syllabusResponses []dto.SyllabusResponse
	for _, syllabus := range syllabuses {
		syllabusResponses = append(syllabusResponses, dto.SyllabusResponse{
			ID:                   syllabus.ID.String(),
			UserID:               syllabus.UserID,
			UserNRP:              syllabus.UserNRP,
			AcademicAdvisorID:    syllabus.AcademicAdvisorID,
			AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
			RegistrationID:       syllabus.RegistrationID,
			Title:                syllabus.Title,
			FileStorageID:        syllabus.FileStorageID,
		})
	}

	return syllabusResponses, nil
}

// Create creates a new syllabus
func (s *syllabusService) Create(ctx context.Context, syllabus dto.SyllabusRequest, file *multipart.FileHeader, token string) (dto.SyllabusResponse, error) {
	syllabusCheck, err := s.syllabusRepo.FindByRegistrationID(ctx, syllabus.RegistrationID, nil)
	if err != nil {
		if err.Error() != "record not found" {
			return dto.SyllabusResponse{}, err
		}
	}

	if syllabusCheck.ID != uuid.Nil {
		return dto.SyllabusResponse{}, errors.New("syllabus already exists")
	}

	if file == nil {
		return dto.SyllabusResponse{}, errors.New("file is required")
	}

	// Upload file to storage
	result, err := s.fileService.storage.GcsUpload(file, "sim_mbkm", "", "")
	if err != nil {
		return dto.SyllabusResponse{}, err
	}

	// Verify user has access to this registration
	user := s.userManagementService.GetUserData("GET", token)
	if user == nil {
		log.Println("ERROR GETTING USER DATA: ", err)
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	registration := s.registrationService.GetRegistrationByID("GET", syllabus.RegistrationID, token)
	if registration == nil {
		log.Println("ERROR GETTING REGISTRATION: ", err)
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	userID, ok := registration["user_id"].(string)
	if !ok {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	if userID != user["id"] {
		log.Println("USER ID DOES NOT MATCH: ", userID, user["id"])
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	// Create syllabus entity
	var syllabusEntity entity.Syllabus
	syllabusEntity.ID = uuid.New()
	syllabusEntity.UserID = userID
	syllabusEntity.UserNRP = registration["user_nrp"].(string)
	syllabusEntity.AcademicAdvisorID = registration["academic_advisor"].(string)
	syllabusEntity.AcademicAdvisorEmail = registration["academic_advisor_email"].(string)
	syllabusEntity.RegistrationID = syllabus.RegistrationID
	syllabusEntity.Title = syllabus.Title
	syllabusEntity.FileStorageID = result.FileID

	// Set timestamps
	now := time.Now()
	syllabusEntity.CreatedAt = &now
	syllabusEntity.UpdatedAt = &now

	// Save to database
	syllabusResponse, err := s.syllabusRepo.Create(ctx, syllabusEntity, nil)
	if err != nil {
		return dto.SyllabusResponse{}, err
	}

	return dto.SyllabusResponse{
		ID:                   syllabusResponse.ID.String(),
		UserID:               syllabusResponse.UserID,
		UserNRP:              syllabusResponse.UserNRP,
		AcademicAdvisorID:    syllabusResponse.AcademicAdvisorID,
		AcademicAdvisorEmail: syllabusResponse.AcademicAdvisorEmail,
		RegistrationID:       syllabusResponse.RegistrationID,
		Title:                syllabusResponse.Title,
		FileStorageID:        syllabusResponse.FileStorageID,
	}, nil
}

// Update updates an existing syllabus
func (s *syllabusService) Update(ctx context.Context, id string, subject dto.SyllabusRequest) error {
	res, err := s.syllabusRepo.FindByID(ctx, id, nil)
	if err != nil {
		return err
	}

	// Create syllabusEntity with original ID
	syllabusEntity := entity.Syllabus{
		ID: res.ID,
	}

	// Get reflection values
	resValue := reflect.ValueOf(res)
	reqValue := reflect.ValueOf(subject)
	entityValue := reflect.ValueOf(&syllabusEntity).Elem()

	// Iterate through fields of the request type
	for i := 0; i < reqValue.Type().NumField(); i++ {
		reqField := reqValue.Type().Field(i)
		reqFieldValue := reqValue.Field(i)

		// Find corresponding field in the entity
		entityField := entityValue.FieldByName(reqField.Name)

		// Find corresponding field in the original result
		resField := resValue.FieldByName(reqField.Name)

		// Check if the field exists and can be set
		if entityField.IsValid() && entityField.CanSet() {
			// If the request field is not zero, use its value
			if !reflect.DeepEqual(reqFieldValue.Interface(), reflect.Zero(reqFieldValue.Type()).Interface()) {
				entityField.Set(reqFieldValue)
			} else {
				// Otherwise, use the original value
				entityField.Set(resField)
			}
		}
	}

	// Update timestamps
	now := time.Now()
	syllabusEntity.UpdatedAt = &now

	// Perform the update
	err = s.syllabusRepo.Update(ctx, id, syllabusEntity, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindByID retrieves a syllabus by its ID
func (s *syllabusService) FindByID(ctx context.Context, id string, token string) (dto.SyllabusResponse, error) {
	syllabus, err := s.syllabusRepo.FindByID(ctx, id, nil)
	if err != nil {
		return dto.SyllabusResponse{}, err
	}

	user := s.userManagementService.GetUserData("GET", token)
	if user == nil {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	userRole, ok := user["role"].(string)
	if !ok {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	if userRole == "DOSEN PEMBIMBING" {
		if syllabus.AcademicAdvisorEmail != user["email"] {
			return dto.SyllabusResponse{}, errors.New("unauthorized")
		}
	} else if userRole == "MAHASISWA" {
		if syllabus.UserNRP != user["nrp"] {
			return dto.SyllabusResponse{}, errors.New("unauthorized")
		}
	} else if userRole != "ADMIN" && userRole != "LO-MBKM" {
		return dto.SyllabusResponse{}, errors.New("unauthorized")
	}

	return dto.SyllabusResponse{
		ID:                   syllabus.ID.String(),
		UserID:               syllabus.UserID,
		UserNRP:              syllabus.UserNRP,
		AcademicAdvisorID:    syllabus.AcademicAdvisorID,
		AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
		RegistrationID:       syllabus.RegistrationID,
		Title:                syllabus.Title,
		FileStorageID:        syllabus.FileStorageID,
	}, nil
}

// Destroy deletes a syllabus
func (s *syllabusService) Destroy(ctx context.Context, id string) error {
	err := s.syllabusRepo.Destroy(ctx, id, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindByRegistrationID retrieves syllabuses by registration ID
func (s *syllabusService) FindByRegistrationID(ctx context.Context, registrationID string) (dto.SyllabusResponse, error) {
	syllabus, err := s.syllabusRepo.FindByRegistrationID(ctx, registrationID, nil)
	if err != nil {
		return dto.SyllabusResponse{}, err
	}

	var syllabusResponses dto.SyllabusResponse

	syllabusResponses = dto.SyllabusResponse{
		ID:                   syllabus.ID.String(),
		UserID:               syllabus.UserID,
		UserNRP:              syllabus.UserNRP,
		AcademicAdvisorID:    syllabus.AcademicAdvisorID,
		AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
		RegistrationID:       syllabus.RegistrationID,
		Title:                syllabus.Title,
		FileStorageID:        syllabus.FileStorageID,
	}

	return syllabusResponses, nil
}

// FindAllByRegistrationID retrieves all syllabuses by registration ID
func (s *syllabusService) FindAllByRegistrationID(ctx context.Context, registrationID string) ([]dto.SyllabusResponse, error) {
	syllabuses, err := s.syllabusRepo.FindAllByRegistrationID(ctx, registrationID, nil)
	if err != nil {
		return nil, err
	}

	var syllabusResponses []dto.SyllabusResponse
	for _, syllabus := range syllabuses {
		syllabusResponses = append(syllabusResponses, dto.SyllabusResponse{
			ID:                   syllabus.ID.String(),
			UserID:               syllabus.UserID,
			UserNRP:              syllabus.UserNRP,
			AcademicAdvisorID:    syllabus.AcademicAdvisorID,
			AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
			RegistrationID:       syllabus.RegistrationID,
			Title:                syllabus.Title,
			FileStorageID:        syllabus.FileStorageID,
		})
	}

	return syllabusResponses, nil
}

func (s *syllabusService) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.SyllabusByStudentResponse, error) {
	user := s.userManagementService.GetUserData("GET", token)

	userNRP, ok := user["nrp"].(string)
	if !ok {
		return dto.SyllabusByStudentResponse{}, errors.New("user NRP not found")
	}

	syllabuses, err := s.syllabusRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)
	if err != nil {
		log.Println("ERROR FINDING SYLLABUSES BY USER NRP AND GROUP BY REGISTRATION ID: ", err)
		return dto.SyllabusByStudentResponse{}, err
	}

	syllabusResponses := make(map[string][]dto.SyllabusResponse)

	for registrationID, syllabusList := range syllabuses {
		registration := s.registrationService.GetRegistrationByID("GET", registrationID, token)
		registrationActivityName, ok := registration["activity_name"].(string)
		if !ok {
			return dto.SyllabusByStudentResponse{}, errors.New("registration activity name not found")
		}

		approvalStatus, ok := registration["approval_status"].(bool)
		if !ok {
			return dto.SyllabusByStudentResponse{}, errors.New("registration approval status not found")
		}

		if !approvalStatus {
			log.Println("REGISTRATION APPROVAL STATUS IS FALSE: ", registrationID)
			continue
		}

		var syllabusResponseList []dto.SyllabusResponse
		for _, syllabus := range syllabusList {
			response := dto.SyllabusResponse{
				ID:                   syllabus.ID.String(),
				UserID:               syllabus.UserID,
				UserNRP:              syllabus.UserNRP,
				AcademicAdvisorID:    syllabus.AcademicAdvisorID,
				AcademicAdvisorEmail: syllabus.AcademicAdvisorEmail,
				RegistrationID:       syllabus.RegistrationID,
				Title:                syllabus.Title,
				FileStorageID:        syllabus.FileStorageID,
			}

			syllabusResponseList = append(syllabusResponseList, response)
		}

		syllabusResponses[registrationActivityName] = syllabusResponseList
	}

	return dto.SyllabusByStudentResponse{
		Syllabuses: syllabusResponses,
	}, nil
}
