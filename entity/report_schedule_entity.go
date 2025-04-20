package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	ReportSchedule struct {
		ID                   uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
		UserID               string     `json:"user_id"`
		UserNRP              string     `json:"user_nrp"`
		RegistrationID       string     `json:"registration_id"`
		AcademicAdvisorID    string     `json:"academic_advisor_id"`
		AcademicAdvisorEmail string     `json:"academic_advisor_email"`
		ReportType           string     `json:"report_type"`
		Week                 int        `json:"week"`
		StartDate            *time.Time `json:"start_date"`
		EndDate              *time.Time `json:"end_date"`
		Report               []Report   `json:"report" gorm:"foreignKey:ReportScheduleID"`
		BaseModel
	}
)
