package entity

import "github.com/google/uuid"

type Transcript struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	RegistrationID string    `gorm:"type:varchar(255);not null"`
	Title          string    `gorm:"type:varchar(255);not null"`
	FileStorageID  string    `gorm:"type:varchar(255);not null"`
	BaseModel
}
