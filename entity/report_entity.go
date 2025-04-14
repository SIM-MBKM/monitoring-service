package entity

import "github.com/google/uuid"

type (
	Report struct {
		ID                    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		ReportScheduleID      string    `json:"report_schedule_id"`
		Title                 string    `json:"title"`
		Content               string    `json:"content"`
		ReportType            string    `json:"report_type"`
		FileStorageID         string    `json:"file_storage_id"`
		Feedback              string    `json:"feedback"`
		AcademicAdvisorStatus string    `json:"academic_advisor_status"`
		BaseModel
	}
)
