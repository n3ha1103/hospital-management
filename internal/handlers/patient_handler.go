package handlers

import (
	"net/http"
	"strconv"

	"hospital-management/internal/models"
	"hospital-management/internal/service"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientService service.PatientService
}

func NewPatientHandler(patientService service.PatientService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

// CreatePatient handles patient creation
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var patientReq models.PatientRequest
	if err := c.ShouldBindJSON(&patientReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if patientReq.FirstName == "" || patientReq.LastName == "" || patientReq.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "First name, last name, and email are required"})
		return
	}

	createdPatient, err := h.patientService.CreatePatient(&patientReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPatient)
}

// SearchPatients handles patient search
func (h *PatientHandler) SearchPatients(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	patients, err := h.patientService.SearchPatients(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patients)
}

// GetPatients handles getting all patients with pagination
func (h *PatientHandler) GetPatients(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	allPatients, err := h.patientService.GetAllPatients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total := len(allPatients)
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}
	patients := allPatients[start:end]

	response := gin.H{
		"patients": patients,
		"total":    total,
		"page":     page,
		"limit":    limit,
	}

	c.JSON(http.StatusOK, response)
}

// GetPatientByID handles getting a specific patient
func (h *PatientHandler) GetPatientByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	patient, err := h.patientService.GetPatientByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// UpdatePatient handles patient updates
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	var patient models.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Format date safely
	dob := ""
	if !patient.DateOfBirth.IsZero() {
		dob = patient.DateOfBirth.Format("2006-01-02")
	}

	patientReq := models.PatientRequest{
		FirstName:      patient.FirstName,
		LastName:       patient.LastName,
		Email:          models.StringValue(patient.Email),
		Phone:          patient.Phone,
		DateOfBirth:    dob,
		Gender:         patient.Gender,
		Address:        models.StringValue(patient.Address),
		MedicalHistory: models.StringValue(patient.MedicalHistory),
		Allergies:      models.StringValue(patient.Allergies),
		Medications:    models.StringValue(patient.Medications),
	}

	updatedPatient, err := h.patientService.UpdatePatient(uint(id), &patientReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPatient)
}

// DeletePatient handles patient deletion
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	if err := h.patientService.DeletePatient(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}
