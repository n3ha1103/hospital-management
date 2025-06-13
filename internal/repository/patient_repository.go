package repository

import (
	"fmt"
	"hospital-management/internal/models"

	"gorm.io/gorm"
)

// PatientRepository defines the interface for patient data operations.
type PatientRepository interface {
	Create(patient *models.Patient) (*models.Patient, error)
	GetByID(id int) (*models.Patient, error)
	GetByPhone(phone string) (*models.Patient, error)
	Update(patient *models.Patient) (*models.Patient, error)
	Delete(id int) error
	GetAll() ([]*models.Patient, error)
	Search(query string) ([]*models.Patient, error)
}

// PatientRepositoryImpl implements PatientRepository using GORM.
type PatientRepositoryImpl struct {
	db *gorm.DB
}

// NewPatientRepository creates a new PatientRepository.
func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &PatientRepositoryImpl{
		db: db,
	}
}

// Create inserts a new patient record into the database.
func (r *PatientRepositoryImpl) Create(patient *models.Patient) (*models.Patient, error) {
	if err := r.db.Create(patient).Error; err != nil {
		return nil, fmt.Errorf("failed to create patient: %w", err)
	}
	return patient, nil
}

// GetByID retrieves a patient by their ID.
func (r *PatientRepositoryImpl) GetByID(id int) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.First(&patient, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("patient with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}
	return &patient, nil
}

// GetByPhone retrieves a patient by their phone number.
func (r *PatientRepositoryImpl) GetByPhone(phone string) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.Where("phone = ?", phone).First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("patient with phone %s not found", phone)
		}
		return nil, fmt.Errorf("failed to get patient by phone: %w", err)
	}
	return &patient, nil
}

// Update modifies an existing patient record.
func (r *PatientRepositoryImpl) Update(patient *models.Patient) (*models.Patient, error) {
	if err := r.db.Save(patient).Error; err != nil {
		return nil, fmt.Errorf("failed to update patient: %w", err)
	}
	return patient, nil
}

// Delete removes a patient record by ID.
func (r *PatientRepositoryImpl) Delete(id int) error {
	result := r.db.Delete(&models.Patient{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete patient: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("patient with id %d not found", id)
	}
	return nil
}

// GetAll retrieves all patients, ordered by creation date descending.
func (r *PatientRepositoryImpl) GetAll() ([]*models.Patient, error) {
	var patients []*models.Patient
	if err := r.db.Order("created_at desc").Find(&patients).Error; err != nil {
		return nil, fmt.Errorf("failed to get patients: %w", err)
	}
	return patients, nil
}

// Search finds patients by first name, last name, email, or phone (case-insensitive).
func (r *PatientRepositoryImpl) Search(query string) ([]*models.Patient, error) {
	var patients []*models.Patient
	searchPattern := "%" + query + "%"
	if err := r.db.Where(
		"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern,
	).Order("created_at desc").Find(&patients).Error; err != nil {
		return nil, fmt.Errorf("failed to search patients: %w", err)
	}
	return patients, nil
}
