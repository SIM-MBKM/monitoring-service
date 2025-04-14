package dto

type (
	ReportRequest struct {
		ReportScheduleID string `json:"report_schedule_id" validate:"required"`
		Title            string `json:"title" validate:"required"`
		Content          string `json:"content" validate:"required"`
		ReportType       string `json:"report_type" validate:"required"`
		Feedback         string `json:"feedback"`
	}

	ReportResponse struct {
		ID                    string `json:"id"`
		ReportScheduleID      string `json:"report_schedule_id"`
		FileStorageID         string `json:"file_storage_id"`
		Title                 string `json:"title"`
		Content               string `json:"content"`
		ReportType            string `json:"report_type"`
		Feedback              string `json:"feedback"`
		AcademicAdvisorStatus string `json:"academic_advisor_status"`
	}
)
