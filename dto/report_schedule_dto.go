package dto

type (
	ReportScheduleRequest struct {
		RegistrationID    string `json:"registration_id" validate:"required"`
		AcademicAdvisorID string `json:"academic_advisor_id" validate:"required"`
		ReportType        string `json:"report_type" validate:"required"`
		Week              string `json:"week" validate:"required"`
		StartDate         string `json:"start_date" validate:"required"`
		EndDate           string `json:"end_date" validate:"required"`
	}

	ReportScheduleResponse struct {
		ID                string `json:"id"`
		RegistrationID    string `json:"registration_id"`
		AcademicAdvisorID string `json:"academic_advisor_id"`
		ReportType        string `json:"report_type"`
		Week              string `json:"week"`
		StartDate         string `json:"start_date"`
		EndDate           string `json:"end_date"`
		Report            ReportResponse
	}
)
