package dto

type (
	ReportRequest struct {
		ReportScheduleID string `json:"report_schedule_id" validate:"required"`
		Title            string `json:"title" validate:"required"`
		Content          string `json:"content" validate:"required"`
		ReportType       string `json:"report_type" validate:"required"`
		Feedback         string `json:"feedback"`
	}
)
