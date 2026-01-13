package models

import "time"

type User struct {
	ID               uint   `gorm:"primaryKey"`
	Email            string `gorm:"uniqueIndex;not null"`
	Password         string `gorm:"not null"`
	Role             string `gorm:"type:varchar(20);not null;default:'user'"`
	ResetToken       *string
	ResetTokenExpiry *time.Time
	CreatedAt        time.Time
}
