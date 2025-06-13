package database

import (
	"hospital-management/internal/config"
	"hospital-management/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Patient{}, &models.Appointment{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
