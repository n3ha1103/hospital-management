package repository

import (
	"hospital-management/internal/models"
	"time"
)

type AppointmentRepository interface {
	Create(appointment *models.Appointment) (*models.Appointment, error)
	GetByID(id uint) (*models.Appointment, error)
	GetAll() ([]*models.Appointment, error)
	GetByPatientID(patientID uint) ([]*models.Appointment, error)
	GetByDoctorID(doctorID uint) ([]*models.Appointment, error)
	GetByDateRange(start, end time.Time) ([]*models.Appointment, error)
	Update(appointment *models.Appointment) (*models.Appointment, error)
	Delete(id uint) error
}

type DoctorRepository interface {
	Create(doctor *models.Doctor) (*models.Doctor, error)
	GetByID(id uint) (*models.Doctor, error)
	GetByEmail(email string) (*models.Doctor, error)
	Update(doctor *models.Doctor) (*models.Doctor, error)
	Delete(id uint) error
	GetAll() ([]*models.Doctor, error)
}
type AuthService interface {
	Login(req *models.LoginRequest) (string, error)
	Register(req *models.RegisterRequest) (*models.User, error)
}

type PatientService interface {
	CreatePatient(req *models.PatientRequest) (*models.Patient, error)
	GetPatientByID(id uint) (*models.Patient, error)
	GetAllPatients() ([]*models.Patient, error)
	UpdatePatient(id uint, req *models.PatientRequest) (*models.Patient, error)
	DeletePatient(id uint) error
}
