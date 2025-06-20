package dto

type (
	ReportRequest struct {
		ReportScheduleID string `form:"report_schedule_id" validate:"required"`
		Title            string `form:"title" validate:"required"`
		Content          string `form:"content"`
		ReportType       string `form:"report_type" validate:"oneof=WEEKLY_REPORT FINAL_REPORT"`
	}

	ReportUpdateRequest struct {
		Title      string `form:"title" validate:"required"`
		Content    string `form:"content" validate:"required"`
		ReportType string `form:"report_type" validate:"required,oneof=WEEKLY_REPORT FINAL_REPORT"`
	}

	ReportApprovalRequest struct {
		Status   string   `json:"status" validate:"required"`
		Feedback string   `json:"feedback"`
		IDs      []string `json:"ids"`
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
