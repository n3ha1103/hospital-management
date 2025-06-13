package service

import (
	"fmt"
	"hospital-management/internal/models"
	"hospital-management/internal/repository"
	"time"
)

type AppointmentService interface {
	CreateAppointment(req *models.AppointmentRequest) (*models.Appointment, error)
	GetAppointmentByID(id uint) (*models.Appointment, error)
	UpdateAppointment(id uint, req *models.AppointmentUpdateRequest) (*models.Appointment, error)
	DeleteAppointment(id uint) error
	GetAppointmentsByPatient(patientID uint) ([]*models.Appointment, error)
	GetAppointmentsByDoctor(doctorID uint) ([]*models.Appointment, error)
	GetAllAppointments() ([]*models.Appointment, error)
}

type appointmentService struct {
	appointmentRepo repository.AppointmentRepository
	patientRepo     repository.PatientRepository
	userRepo        repository.UserRepository
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository, patientRepo repository.PatientRepository, userRepo repository.UserRepository) AppointmentService {
	return &appointmentService{
		appointmentRepo: appointmentRepo,
		patientRepo:     patientRepo,
		userRepo:        userRepo,
	}
}

func (s *appointmentService) CreateAppointment(req *models.AppointmentRequest) (*models.Appointment, error) {
	// Validate patient exists - convert uint to int
	_, err := s.patientRepo.GetByID(int(req.PatientID))
	if err != nil {
		return nil, fmt.Errorf("patient not found: %w", err)
	}

	// Validate doctor exists - convert uint to int
	doctor, err := s.userRepo.GetByID(uint(req.DoctorID))
	if err != nil {
		return nil, fmt.Errorf("doctor not found: %w", err)
	}

	// Validate doctor role
	if doctor.Role != "doctor" {
		return nil, fmt.Errorf("user is not a doctor")
	}

	// Parse DateTime from string to time.Time
	parsedDateTime, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		return nil, fmt.Errorf("invalid date_time format: %w", err)
	}

	appointment := &models.Appointment{
		PatientID: req.PatientID,
		DoctorID:  req.DoctorID,
		DateTime:  parsedDateTime,
		Duration:  req.Duration,
		Status:    "scheduled",
		Notes:     models.StringPtr(req.Notes),
		CreatedBy: 0, // Set appropriately if needed
	}

	createdAppointment, err := s.appointmentRepo.Create(appointment)
	if err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return createdAppointment, nil
}

func (s *appointmentService) GetAppointmentByID(id uint) (*models.Appointment, error) {
	appointment, err := s.appointmentRepo.GetByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}
	return appointment, nil
}

func (s *appointmentService) UpdateAppointment(id uint, req *models.AppointmentUpdateRequest) (*models.Appointment, error) {
	// Get existing appointment
	appointment, err := s.appointmentRepo.GetByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("appointment not found: %w", err)
	}

	// Update fields
	if req.DateTime != "" {
		parsedDateTime, err := time.Parse(time.RFC3339, req.DateTime)
		if err == nil {
			appointment.DateTime = parsedDateTime
		}
	}
	if req.Duration != 0 {
		appointment.Duration = req.Duration
	}
	if req.Status != "" {
		appointment.Status = req.Status
	}
	if req.Notes != "" {
		appointment.Notes = models.StringPtr(req.Notes)
	}
	if req.Diagnosis != "" {
		appointment.Diagnosis = models.StringPtr(req.Diagnosis)
	}
	if req.Treatment != "" {
		appointment.Treatment = models.StringPtr(req.Treatment)
	}

	updatedAppointment, err := s.appointmentRepo.Update(appointment)
	if err != nil {
		return nil, fmt.Errorf("failed to update appointment: %w", err)
	}

	return updatedAppointment, nil
}

func (s *appointmentService) DeleteAppointment(id uint) error {
	err := s.appointmentRepo.Delete(uint(id))
	if err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}
	return nil
}

func (s *appointmentService) GetAppointmentsByPatient(patientID uint) ([]*models.Appointment, error) {
	appointments, err := s.appointmentRepo.GetByPatientID(uint(patientID))
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments for patient: %w", err)
	}
	return appointments, nil
}

func (s *appointmentService) GetAppointmentsByDoctor(doctorID uint) ([]*models.Appointment, error) {
	appointments, err := s.appointmentRepo.GetByDoctorID(uint(doctorID))
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments for doctor: %w", err)
	}
	return appointments, nil
}

func (s *appointmentService) GetAllAppointments() ([]*models.Appointment, error) {
	appointments, err := s.appointmentRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all appointments: %w", err)
	}
	return appointments, nil
}
