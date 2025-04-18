package dto

type (
	TranscriptRequest struct {
		RegistrationID string `form:"registration_id" validate:"required"`
		Title          string `form:"title" validate:"required"`
	}

	TranscriptResponse struct {
		ID             string `json:"id"`
		RegistrationID string `json:"registration_id"`
		Title          string `json:"title"`
		FileStorageID  string `json:"file_storage_id"`
	}
)
