package models

import "time"

type Patient struct {
	ID             uint      `json:"id" db:"id"`
	FirstName      string    `json:"first_name" db:"first_name" validate:"required"`
	LastName       string    `json:"last_name" db:"last_name" validate:"required"`
	Email          *string   `json:"email" db:"email" validate:"omitempty,email"`
	Phone          string    `json:"phone" db:"phone" validate:"required"`
	DateOfBirth    time.Time `json:"date_of_birth" db:"date_of_birth" validate:"required"`
	Gender         string    `json:"gender" db:"gender" validate:"required,oneof=male female other"`
	Address        *string   `json:"address" db:"address"`
	MedicalHistory *string   `json:"medical_history" db:"medical_history"`
	Allergies      *string   `json:"allergies" db:"allergies"`
	Medications    *string   `json:"medications" db:"medications"`
	CreatedBy      uint      `json:"created_by" db:"created_by" validate:"required"`
	UpdatedBy      *uint     `json:"updated_by" db:"updated_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Helper method to get full name
func (p *Patient) GetFullName() string {
	return p.FirstName + " " + p.LastName
}

// Request/Response types for API
type PatientRequest struct {
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	Email          string `json:"email" validate:"omitempty,email"`
	Phone          string `json:"phone" validate:"required"`
	DateOfBirth    string `json:"date_of_birth" validate:"required"` // Will be parsed to time.Time
	Gender         string `json:"gender" validate:"required,oneof=male female other"`
	Address        string `json:"address"`
	MedicalHistory string `json:"medical_history"`
	Allergies      string `json:"allergies"`
	Medications    string `json:"medications"`
}

type PatientResponse struct {
	ID             uint    `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	FullName       string  `json:"full_name"`
	Email          *string `json:"email"`
	Phone          string  `json:"phone"`
	DateOfBirth    string  `json:"date_of_birth"` // Formatted as string for API
	Gender         string  `json:"gender"`
	Address        *string `json:"address"`
	MedicalHistory *string `json:"medical_history"`
	Allergies      *string `json:"allergies"`
	Medications    *string `json:"medications"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// TableName returns the table name for Patient model
func (Patient) TableName() string {
	return "patients"
}
