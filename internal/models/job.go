package models

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uint      `gorm:"index;not null"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	Status  string `gorm:"type:varchar(20);not null"`
	ZipPath *string
	Error   *string

	Files []JobFile `gorm:"foreignKey:JobID"`

	CreatedAt time.Time
}
