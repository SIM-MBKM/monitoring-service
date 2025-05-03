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

type transcriptService struct {
	transcriptRepo        repository.TranscriptRepository
	fileService           *FileService
	userManagementService *UserManagementService
	registrationService   *RegistrationManagementService
}

type TranscriptService interface {
	Index(ctx context.Context) ([]dto.TranscriptResponse, error)
	Create(ctx context.Context, transcript dto.TranscriptRequest, file *multipart.FileHeader, token string) (dto.TranscriptResponse, error)
	Update(ctx context.Context, id string, transcript dto.TranscriptRequest) error
	FindByID(ctx context.Context, id string, token string) (dto.TranscriptResponse, error)
	Destroy(ctx context.Context, id string) error
	FindByRegistrationID(ctx context.Context, registrationID string) (dto.TranscriptResponse, error)
	FindAllByRegistrationID(ctx context.Context, registrationID string) ([]dto.TranscriptResponse, error)
	FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, token string, pagReq dto.PaginationRequest, filter dto.TranscriptAdvisorFilterRequest) (dto.TranscriptAdvisorResponse, dto.PaginationResponse, error)
	FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.TranscriptByStudentResponse, error)
}

func NewTranscriptService(
	transcriptRepo repository.TranscriptRepository,
	userManagementBaseURI string,
	registrationBaseURI string,
	asyncURIs []string,
	config *storageService.Config,
	tokenManager *storageService.CacheTokenManager,
) TranscriptService {
	return &transcriptService{
		transcriptRepo:        transcriptRepo,
		fileService:           NewFileService(config, tokenManager),
		userManagementService: NewUserManagementService(userManagementBaseURI, asyncURIs),
		registrationService:   NewRegistrationManagementService(registrationBaseURI, asyncURIs),
	}
}

func (s *transcriptService) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, token string, pagReq dto.PaginationRequest, filter dto.TranscriptAdvisorFilterRequest) (dto.TranscriptAdvisorResponse, dto.PaginationResponse, error) {
	user := s.userManagementService.GetUserData("GET", token)
	advisorEmail, ok := user["email"].(string)
	if !ok {
		return dto.TranscriptAdvisorResponse{}, dto.PaginationResponse{}, errors.New("advisor email not found")
	}

	transcripts, totalCount, err := s.transcriptRepo.FindByAdvisorEmailAndGroupByUserNRP(ctx, advisorEmail, nil, &pagReq, filter.UserNRP)
	if err != nil {
		return dto.TranscriptAdvisorResponse{}, dto.PaginationResponse{}, err
	}

	transcriptResponses := make(map[string][]dto.TranscriptResponse)
	for userNRP, transcript := range transcripts {
		// get registration by registration id
		registration := s.registrationService.GetRegistrationByID("GET", transcript.RegistrationID, token)
		if registration == nil {
			return dto.TranscriptAdvisorResponse{}, dto.PaginationResponse{}, errors.New("registration not found")
		}

		registrationActivityName, ok := registration["activity_name"].(string)
		if !ok {
			return dto.TranscriptAdvisorResponse{}, dto.PaginationResponse{}, errors.New("registration activity name not found")
		}

		transcriptResponses[userNRP] = append(transcriptResponses[userNRP], dto.TranscriptResponse{
			ID:                   transcript.ID.String(),
			UserID:               transcript.UserID,
			UserNRP:              transcript.UserNRP,
			ActivityName:         registrationActivityName,
			AcademicAdvisorID:    transcript.AcademicAdvisorID,
			AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
			RegistrationID:       transcript.RegistrationID,
			Title:                transcript.Title,
			FileStorageID:        transcript.FileStorageID,
		})
	}

	// Generate pagination metadata
	paginationResponse := helper.MetaDataPagination(totalCount, pagReq)

	return dto.TranscriptAdvisorResponse{
		Transcripts: transcriptResponses,
	}, paginationResponse, nil
}

// Index retrieves all transcripts
func (s *transcriptService) Index(ctx context.Context) ([]dto.TranscriptResponse, error) {
	transcripts, err := s.transcriptRepo.Index(ctx, nil)
	if err != nil {
		return nil, err
	}

	var transcriptResponses []dto.TranscriptResponse
	for _, transcript := range transcripts {
		transcriptResponses = append(transcriptResponses, dto.TranscriptResponse{
			ID:                   transcript.ID.String(),
			UserID:               transcript.UserID,
			UserNRP:              transcript.UserNRP,
			AcademicAdvisorID:    transcript.AcademicAdvisorID,
			AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
			RegistrationID:       transcript.RegistrationID,
			Title:                transcript.Title,
			FileStorageID:        transcript.FileStorageID,
		})
	}

	return transcriptResponses, nil
}

// Create creates a new transcript
func (s *transcriptService) Create(ctx context.Context, transcript dto.TranscriptRequest, file *multipart.FileHeader, token string) (dto.TranscriptResponse, error) {
	transcriptCheck, err := s.transcriptRepo.FindByRegistrationID(ctx, transcript.RegistrationID, nil)
	if err != nil {
		if err.Error() != "record not found" {
			return dto.TranscriptResponse{}, err
		}
	}

	if transcriptCheck.ID != uuid.Nil {
		return dto.TranscriptResponse{}, errors.New("transcript already exists")
	}

	if file == nil {
		return dto.TranscriptResponse{}, errors.New("file is required")
	}

	// Upload file to storage
	result, err := s.fileService.storage.GcsUpload(file, "sim_mbkm", "", "")
	if err != nil {
		return dto.TranscriptResponse{}, err
	}

	// Verify user has access to this registration
	user := s.userManagementService.GetUserData("GET", token)
	if user == nil {
		log.Println("ERROR GETTING USER DATA: ", err)
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	registration := s.registrationService.GetRegistrationByID("GET", transcript.RegistrationID, token)
	if registration == nil {
		log.Println("ERROR GETTING REGISTRATION: ", err)
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	userID, ok := registration["user_id"].(string)
	if !ok {
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	if userID != user["id"] {
		log.Println("USER ID DOES NOT MATCH: ", userID, user["id"])
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}
	// Create transcript entity
	var transcriptEntity entity.Transcript
	transcriptEntity.ID = uuid.New()
	transcriptEntity.UserID = userID
	transcriptEntity.UserNRP = registration["user_nrp"].(string)
	transcriptEntity.AcademicAdvisorID = registration["academic_advisor"].(string)
	transcriptEntity.AcademicAdvisorEmail = registration["academic_advisor_email"].(string)
	transcriptEntity.RegistrationID = transcript.RegistrationID
	transcriptEntity.Title = transcript.Title
	transcriptEntity.FileStorageID = result.FileID

	// Set timestamps
	now := time.Now()
	transcriptEntity.CreatedAt = &now
	transcriptEntity.UpdatedAt = &now

	// Save to database
	transcriptResponse, err := s.transcriptRepo.Create(ctx, transcriptEntity, nil)
	if err != nil {
		return dto.TranscriptResponse{}, err
	}

	return dto.TranscriptResponse{
		ID:                   transcriptResponse.ID.String(),
		UserID:               transcriptResponse.UserID,
		UserNRP:              transcriptResponse.UserNRP,
		AcademicAdvisorID:    transcriptResponse.AcademicAdvisorID,
		AcademicAdvisorEmail: transcriptResponse.AcademicAdvisorEmail,
		RegistrationID:       transcriptResponse.RegistrationID,
		Title:                transcriptResponse.Title,
		FileStorageID:        transcriptResponse.FileStorageID,
	}, nil
}

// Update updates an existing transcript
func (s *transcriptService) Update(ctx context.Context, id string, subject dto.TranscriptRequest) error {
	res, err := s.transcriptRepo.FindByID(ctx, id, nil)
	if err != nil {
		return err
	}

	// Create transcriptEntity with original ID
	transcriptEntity := entity.Transcript{
		ID: res.ID,
	}

	// Get reflection values
	resValue := reflect.ValueOf(res)
	reqValue := reflect.ValueOf(subject)
	entityValue := reflect.ValueOf(&transcriptEntity).Elem()

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
	transcriptEntity.UpdatedAt = &now

	// Perform the update
	err = s.transcriptRepo.Update(ctx, id, transcriptEntity, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindByID retrieves a transcript by its ID
func (s *transcriptService) FindByID(ctx context.Context, id string, token string) (dto.TranscriptResponse, error) {
	transcript, err := s.transcriptRepo.FindByID(ctx, id, nil)
	if err != nil {
		return dto.TranscriptResponse{}, err
	}

	user := s.userManagementService.GetUserData("GET", token)
	if user == nil {
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	userRole, ok := user["role"].(string)
	if !ok {
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	if userRole == "DOSEN PEMBIMBING" {
		if transcript.AcademicAdvisorEmail != user["email"] {
			return dto.TranscriptResponse{}, errors.New("unauthorized")
		}
	} else if userRole == "MAHASISWA" {
		if transcript.UserNRP != user["nrp"] {
			return dto.TranscriptResponse{}, errors.New("unauthorized")
		}
	} else if userRole != "ADMIN" && userRole != "LO-MBKM" {
		return dto.TranscriptResponse{}, errors.New("unauthorized")
	}

	return dto.TranscriptResponse{
		ID:                   transcript.ID.String(),
		UserID:               transcript.UserID,
		UserNRP:              transcript.UserNRP,
		AcademicAdvisorID:    transcript.AcademicAdvisorID,
		AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
		RegistrationID:       transcript.RegistrationID,
		Title:                transcript.Title,
		FileStorageID:        transcript.FileStorageID,
	}, nil
}

// Destroy deletes a transcript
func (s *transcriptService) Destroy(ctx context.Context, id string) error {
	err := s.transcriptRepo.Destroy(ctx, id, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindByRegistrationID retrieves transcripts by registration ID
func (s *transcriptService) FindByRegistrationID(ctx context.Context, registrationID string) (dto.TranscriptResponse, error) {
	transcript, err := s.transcriptRepo.FindByRegistrationID(ctx, registrationID, nil)
	if err != nil {
		return dto.TranscriptResponse{}, err
	}

	var transcriptResponses dto.TranscriptResponse

	transcriptResponses = dto.TranscriptResponse{
		ID:                   transcript.ID.String(),
		UserID:               transcript.UserID,
		UserNRP:              transcript.UserNRP,
		AcademicAdvisorID:    transcript.AcademicAdvisorID,
		AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
		RegistrationID:       transcript.RegistrationID,
		Title:                transcript.Title,
		FileStorageID:        transcript.FileStorageID,
	}

	return transcriptResponses, nil
}

func (s *transcriptService) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, token string) (dto.TranscriptByStudentResponse, error) {
	user := s.userManagementService.GetUserData("GET", token)

	userNRP, ok := user["nrp"].(string)
	if !ok {
		return dto.TranscriptByStudentResponse{}, errors.New("user NRP not found")
	}

	transcripts, err := s.transcriptRepo.FindByUserNRPAndGroupByRegistrationID(ctx, userNRP, nil)
	if err != nil {
		log.Println("ERROR FINDING TRANSCRIPTS BY USER NRP AND GROUP BY REGISTRATION ID: ", err)
		return dto.TranscriptByStudentResponse{}, err
	}

	transcriptResponses := make(map[string][]dto.TranscriptResponse)

	for registrationID, transcriptList := range transcripts {
		registration := s.registrationService.GetRegistrationByID("GET", registrationID, token)
		registrationActivityName, ok := registration["activity_name"].(string)
		if !ok {
			return dto.TranscriptByStudentResponse{}, errors.New("registration activity name not found")
		}

		var transcriptResponseList []dto.TranscriptResponse
		for _, transcript := range transcriptList {
			response := dto.TranscriptResponse{
				ID:                   transcript.ID.String(),
				UserID:               transcript.UserID,
				UserNRP:              transcript.UserNRP,
				AcademicAdvisorID:    transcript.AcademicAdvisorID,
				AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
				RegistrationID:       transcript.RegistrationID,
				Title:                transcript.Title,
				FileStorageID:        transcript.FileStorageID,
			}

			transcriptResponseList = append(transcriptResponseList, response)
		}

		transcriptResponses[registrationActivityName] = transcriptResponseList
	}

	return dto.TranscriptByStudentResponse{
		Transcripts: transcriptResponses,
	}, nil
}

// FindAllByRegistrationID retrieves all transcripts by registration ID
func (s *transcriptService) FindAllByRegistrationID(ctx context.Context, registrationID string) ([]dto.TranscriptResponse, error) {
	transcripts, err := s.transcriptRepo.FindAllByRegistrationID(ctx, registrationID, nil)
	if err != nil {
		return nil, err
	}

	var transcriptResponses []dto.TranscriptResponse
	for _, transcript := range transcripts {
		transcriptResponses = append(transcriptResponses, dto.TranscriptResponse{
			ID:                   transcript.ID.String(),
			UserID:               transcript.UserID,
			UserNRP:              transcript.UserNRP,
			AcademicAdvisorID:    transcript.AcademicAdvisorID,
			AcademicAdvisorEmail: transcript.AcademicAdvisorEmail,
			RegistrationID:       transcript.RegistrationID,
			Title:                transcript.Title,
			FileStorageID:        transcript.FileStorageID,
		})
	}

	return transcriptResponses, nil
}
