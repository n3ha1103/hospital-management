package models

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"not null"`
	Email          string         `json:"email" gorm:"unique;not null"`
	Phone          string         `json:"phone"`
	Specialization string         `json:"specialization" gorm:"not null"`
	Experience     int            `json:"experience"` // Years of experience
	LicenseNumber  string         `json:"license_number" gorm:"unique;not null"`
	Department     string         `json:"department"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:DoctorID"`
}
