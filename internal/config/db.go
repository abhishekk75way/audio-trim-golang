package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(str string) error {
	db, err := gorm.Open(postgres.Open(str), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}
