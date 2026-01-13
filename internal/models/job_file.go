package models

import (
	"time"

	"github.com/google/uuid"
)

type JobFile struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	JobID uuid.UUID `gorm:"type:uuid;index"`

	InputPath  string
	OutputPath *string

	StartTime int
	EndTime   int

	Status string `gorm:"type:varchar(20);not null"`
	Error  *string

	CreatedAt time.Time
}
