package models

import "time"

type Appointment struct {
	ID        uint      `json:"id" db:"id"`
	PatientID uint      `json:"patient_id" db:"patient_id" validate:"required"`
	DoctorID  uint      `json:"doctor_id" db:"doctor_id" validate:"required"`
	DateTime  time.Time `json:"date_time" db:"date_time" validate:"required"`
	Duration  int       `json:"duration" db:"duration" validate:"required,min=15,max=240"` // in minutes
	Status    string    `json:"status" db:"status" validate:"required,oneof=scheduled completed cancelled"`
	Notes     *string   `json:"notes" db:"notes"`
	Diagnosis *string   `json:"diagnosis" db:"diagnosis"`
	Treatment *string   `json:"treatment" db:"treatment"`
	CreatedBy uint      `json:"created_by" db:"created_by" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Populated by joins
	Patient *Patient `gorm:"foreignKey:PatientID"`
	Doctor  *User    `gorm:"foreignKey:DoctorID"`
}

// Request/Response types for Appointments
type AppointmentRequest struct {
	PatientID uint   `json:"patient_id" validate:"required"`
	DoctorID  uint   `json:"doctor_id" validate:"required"`
	DateTime  string `json:"date_time" validate:"required"` // Will be parsed to time.Time
	Duration  int    `json:"duration" validate:"required,min=15,max=240"`
	Notes     string `json:"notes"`
}

type AppointmentUpdateRequest struct {
	DateTime  string `json:"date_time"`
	Duration  int    `json:"duration" validate:"omitempty,min=15,max=240"`
	Status    string `json:"status" validate:"omitempty,oneof=scheduled completed cancelled"`
	Notes     string `json:"notes"`
	Diagnosis string `json:"diagnosis"`
	Treatment string `json:"treatment"`
}

type AppointmentResponse struct {
	ID        uint             `json:"id"`
	PatientID uint             `json:"patient_id"`
	DoctorID  uint             `json:"doctor_id"`
	DateTime  string           `json:"date_time"`
	Duration  int              `json:"duration"`
	Status    string           `json:"status"`
	Notes     *string          `json:"notes"`
	Diagnosis *string          `json:"diagnosis"`
	Treatment *string          `json:"treatment"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
	Patient   *PatientResponse `json:"patient,omitempty"`
	Doctor    *UserResponse    `json:"doctor,omitempty"`
}

// Helper functions
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func UintPtr(u uint) *uint {
	if u == 0 {
		return nil
	}
	return &u
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func UintValue(u *uint) uint {
	if u == nil {
		return 0
	}
	return *u
}
