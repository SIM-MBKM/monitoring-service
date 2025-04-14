package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	ReportSchedule struct {
		ID                uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
		UserID            string     `json:"user_id"`
		RegistrationID    string     `json:"registration_id"`
		AcademicAdvisorID string     `json:"academic_advisor_id"`
		ReportType        string     `json:"report_type"`
		Week              string     `json:"week"`
		StartDate         *time.Time `json:"start_date"`
		EndDate           *time.Time `json:"end_date"`
		Report            []Report   `json:"report" gorm:"foreignKey:ReportScheduleID"`
		BaseModel
	}
)
