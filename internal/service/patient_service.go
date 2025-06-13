// internal/service/patient_service.go
package service

import (
	"fmt"
	"hospital-management/internal/models"
	"hospital-management/internal/repository"
	"time"
)

type PatientService interface {
	CreatePatient(req *models.PatientRequest) (*models.Patient, error)
	GetPatientByID(id uint) (*models.Patient, error)
	GetPatientByPhone(phone string) (*models.Patient, error)
	UpdatePatient(id uint, req *models.PatientRequest) (*models.Patient, error)
	DeletePatient(id uint) error
	GetAllPatients() ([]*models.Patient, error)
	SearchPatients(query string) ([]*models.Patient, error)
}

type patientService struct {
	patientRepo repository.PatientRepository
}

func NewPatientService(patientRepo repository.PatientRepository) PatientService {
	return &patientService{
		patientRepo: patientRepo,
	}
}

func (s *patientService) CreatePatient(req *models.PatientRequest) (*models.Patient, error) {
	// Check if patient with phone already exists
	existingPatient, _ := s.patientRepo.GetByPhone(req.Phone)
	if existingPatient != nil {
		return nil, fmt.Errorf("patient with phone number already exists")
	}

	// Parse DateOfBirth from string to time.Time
	parsedDOB, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("invalid date_of_birth format: %w", err)
	}

	patient := &models.Patient{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          models.StringPtr(req.Email),
		Phone:          req.Phone,
		DateOfBirth:    parsedDOB,
		Gender:         req.Gender,
		Address:        models.StringPtr(req.Address),
		MedicalHistory: models.StringPtr(req.MedicalHistory),
		Allergies:      models.StringPtr(req.Allergies),
		Medications:    models.StringPtr(req.Medications),
		CreatedBy:      0, // Set appropriately if needed
	}

	createdPatient, err := s.patientRepo.Create(patient)
	if err != nil {
		return nil, fmt.Errorf("failed to create patient: %w", err)
	}

	return createdPatient, nil
}

func (s *patientService) GetPatientByID(id uint) (*models.Patient, error) {
	patient, err := s.patientRepo.GetByID(int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}
	return patient, nil
}

func (s *patientService) GetPatientByPhone(phone string) (*models.Patient, error) {
	patient, err := s.patientRepo.GetByPhone(phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient by phone: %w", err)
	}
	return patient, nil
}

func (s *patientService) UpdatePatient(id uint, req *models.PatientRequest) (*models.Patient, error) {
	// Get existing patient
	patient, err := s.patientRepo.GetByID(int(id))
	if err != nil {
		return nil, fmt.Errorf("patient not found: %w", err)
	}

	// Update fields (only update if non-empty)
	if req.FirstName != "" {
		patient.FirstName = req.FirstName
	}
	if req.LastName != "" {
		patient.LastName = req.LastName
	}
	if req.Email != "" {
		patient.Email = models.StringPtr(req.Email)
	}
	if req.Phone != "" {
		patient.Phone = req.Phone
	}
	if req.DateOfBirth != "" {
		parsedDOB, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			patient.DateOfBirth = parsedDOB
		}
	}
	if req.Gender != "" {
		patient.Gender = req.Gender
	}
	if req.Address != "" {
		patient.Address = models.StringPtr(req.Address)
	}
	if req.MedicalHistory != "" {
		patient.MedicalHistory = models.StringPtr(req.MedicalHistory)
	}
	if req.Allergies != "" {
		patient.Allergies = models.StringPtr(req.Allergies)
	}
	if req.Medications != "" {
		patient.Medications = models.StringPtr(req.Medications)
	}

	updatedPatient, err := s.patientRepo.Update(patient)
	if err != nil {
		return nil, fmt.Errorf("failed to update patient: %w", err)
	}

	return updatedPatient, nil
}

func (s *patientService) DeletePatient(id uint) error {
	err := s.patientRepo.Delete(int(id))
	if err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}
	return nil
}

func (s *patientService) GetAllPatients() ([]*models.Patient, error) {
	patients, err := s.patientRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all patients: %w", err)
	}
	return patients, nil
}

func (s *patientService) SearchPatients(query string) ([]*models.Patient, error) {
	patients, err := s.patientRepo.Search(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search patients: %w", err)
	}
	return patients, nil
}
