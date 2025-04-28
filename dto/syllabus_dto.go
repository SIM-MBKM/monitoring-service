package dto

type (
	SyllabusRequest struct {
		RegistrationID string `form:"registration_id" validate:"required"`
		Title          string `form:"title" validate:"required"`
	}

	SyllabusResponse struct {
		ID                   string `json:"id"`
		UserID               string `json:"user_id"`
		UserNRP              string `json:"user_nrp"`
		AcademicAdvisorID    string `json:"academic_advisor_id"`
		AcademicAdvisorEmail string `json:"academic_advisor_email"`
		RegistrationID       string `json:"registration_id"`
		Title                string `json:"title"`
		FileStorageID        string `json:"file_storage_id"`
	}

	SyllabusAdvisorResponse struct {
		Syllabuses map[string][]SyllabusResponse `json:"syllabuses"`
	}

	SyllabusByStudentResponse struct {
		Syllabuses map[string][]SyllabusResponse `json:"syllabuses"`
	}
)
