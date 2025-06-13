// repository/appointment_repository.go
package repository

import (
	"fmt"
	"hospital-management/internal/models"
	"time"

	"gorm.io/gorm"
)

type AppointmentRepositoryImpl struct {
	db *gorm.DB
}

// internal/repository/appointment_repository.go
func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &AppointmentRepositoryImpl{db: db}
}

func (r *AppointmentRepositoryImpl) Create(appointment *models.Appointment) (*models.Appointment, error) {
	now := time.Now()
	appointment.CreatedAt = now
	appointment.UpdatedAt = now

	if err := r.db.Create(appointment).Error; err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, nil
}

func (r *AppointmentRepositoryImpl) GetByID(id uint) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.Preload("Patient").Preload("Doctor").First(&appointment, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("appointment with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get appointment by id: %w", err)
	}

	return &appointment, nil
}

func (r *AppointmentRepositoryImpl) GetAll() ([]*models.Appointment, error) {
	var appointments []*models.Appointment

	err := r.db.
		Preload("Patient").
		Preload("Doctor").
		Order("date_time DESC").
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get all appointments: %w", err)
	}

	return appointments, nil
}

func (r *AppointmentRepositoryImpl) Update(appointment *models.Appointment) (*models.Appointment, error) {
	appointment.UpdatedAt = time.Now()

	if err := r.db.Model(&models.Appointment{}).
		Where("id = ?", appointment.ID).
		Updates(map[string]interface{}{
			"patient_id": appointment.PatientID,
			"doctor_id":  appointment.DoctorID,
			"date_time":  appointment.DateTime,
			"duration":   appointment.Duration,
			"status":     appointment.Status,
			"notes":      appointment.Notes,
			"diagnosis":  appointment.Diagnosis,
			"treatment":  appointment.Treatment,
			"updated_at": appointment.UpdatedAt,
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to update appointment: %w", err)
	}

	return appointment, nil
}

func (r *AppointmentRepositoryImpl) Delete(id uint) error {
	result := r.db.Delete(&models.Appointment{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete appointment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("appointment with id %d not found", id)
	}
	return nil
}

func (r *AppointmentRepositoryImpl) GetByPatientID(patientID uint) ([]*models.Appointment, error) {
	var appointments []*models.Appointment
	err := r.db.
		Preload("Doctor").
		Where("patient_id = ?", patientID).
		Order("date_time DESC").
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by patient id: %w", err)
	}

	return appointments, nil
}

func (r *AppointmentRepositoryImpl) GetByDoctorID(doctorID uint) ([]*models.Appointment, error) {
	var appointments []*models.Appointment
	err := r.db.
		Preload("Patient").
		Where("doctor_id = ?", doctorID).
		Order("date_time ASC").
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by doctor id: %w", err)
	}

	return appointments, nil
}

func (r *AppointmentRepositoryImpl) GetByDateRange(startDate, endDate time.Time) ([]*models.Appointment, error) {
	var appointments []*models.Appointment

	err := r.db.
		Preload("Patient").
		Preload("Doctor").
		Where("date_time BETWEEN ? AND ?", startDate, endDate).
		Order("date_time ASC").
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by date range: %w", err)
	}

	return appointments, nil
}

func (r *AppointmentRepositoryImpl) GetByStatus(status string) ([]*models.Appointment, error) {
	var appointments []*models.Appointment

	err := r.db.
		Preload("Patient").
		Preload("Doctor").
		Where("status = ?", status).
		Order("date_time ASC").
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by status: %w", err)
	}

	return appointments, nil
}

func (r *AppointmentRepositoryImpl) GetUpcomingAppointments(doctorID uint, days int) ([]*models.Appointment, error) {
	var appointments []*models.Appointment
	now := time.Now()
	future := now.Add(time.Duration(days) * 24 * time.Hour)

	err := r.db.
		Preload("Patient").
		Where("doctor_id = ? AND date_time >= ? AND date_time <= ? AND status = ?", doctorID, now, future, "scheduled").
		Order("date_time ASC").
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming appointments: %w", err)
	}

	return appointments, nil
}

func (r *AppointmentRepositoryImpl) GetTodaysAppointments(doctorID uint) ([]*models.Appointment, error) {
	var appointments []*models.Appointment

	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	err := r.db.
		Preload("Patient").
		Where("doctor_id = ? AND date_time >= ? AND date_time < ?", doctorID, start, end).
		Order("date_time ASC").
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get today's appointments: %w", err)
	}

	return appointments, nil
}

func (r *AppointmentRepositoryImpl) GetWithPagination(offset, limit int) ([]*models.Appointment, error) {
	var appointments []*models.Appointment

	err := r.db.
		Preload("Patient").
		Preload("Doctor").
		Order("date_time DESC").
		Limit(limit).
		Offset(offset).
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get appointments with pagination: %w", err)
	}

	return appointments, nil
}

func (r *AppointmentRepositoryImpl) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Appointment{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count appointments: %w", err)
	}
	return count, nil
}

func (r *AppointmentRepositoryImpl) CheckConflict(doctorID uint, dateTime time.Time, duration int, excludeID uint) (bool, error) {
	endTime := dateTime.Add(time.Duration(duration) * time.Minute)

	var count int64
	err := r.db.Model(&models.Appointment{}).
		Where("doctor_id = ? AND status = ? AND id != ?", doctorID, "scheduled", excludeID).
		Where(`
			(date_time <= ? AND date_time + interval '1 minute' * duration > ?) OR 
			(date_time < ? AND date_time + interval '1 minute' * duration >= ?) OR 
			(date_time >= ? AND date_time < ?)`,
			dateTime, dateTime,
			endTime, endTime,
			dateTime, endTime,
		).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check appointment conflict: %w", err)
	}

	return count > 0, nil
}
