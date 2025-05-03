package dto

type (
	TranscriptRequest struct {
		RegistrationID string `form:"registration_id" validate:"required"`
		Title          string `form:"title" validate:"required"`
	}

	TranscriptAdvisorFilterRequest struct {
		UserNRP string `json:"user_nrp"`
	}

	TranscriptResponse struct {
		ID                   string `json:"id"`
		UserID               string `json:"user_id"`
		ActivityName         string `json:"activity_name"`
		UserNRP              string `json:"user_nrp"`
		AcademicAdvisorID    string `json:"academic_advisor_id"`
		AcademicAdvisorEmail string `json:"academic_advisor_email"`
		RegistrationID       string `json:"registration_id"`
		Title                string `json:"title"`
		FileStorageID        string `json:"file_storage_id"`
	}

	TranscriptAdvisorResponse struct {
		Transcripts map[string][]TranscriptResponse `json:"transcripts"`
	}

	TranscriptByStudentResponse struct {
		Transcripts map[string][]TranscriptResponse `json:"transcripts"`
	}
)
