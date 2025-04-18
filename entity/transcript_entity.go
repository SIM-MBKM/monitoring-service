package entity

import "github.com/google/uuid"

type Transcript struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID               string    `json:"user_id" gorm:"type:varchar(255)"`
	UserNRP              string    `json:"user_nrp" gorm:"type:varchar(255)"`
	AcademicAdvisorID    string    `json:"academic_advisor_id" gorm:"type:varchar(255)"`
	AcademicAdvisorEmail string    `json:"academic_advisor_email" gorm:"type:varchar(255)"`
	RegistrationID       string    `json:"registration_id" gorm:"type:varchar(255);not null"`
	Title                string    `json:"title" gorm:"type:varchar(255);not null"`
	FileStorageID        string    `json:"file_storage_id" gorm:"type:varchar(255);not null"`
	BaseModel
}
