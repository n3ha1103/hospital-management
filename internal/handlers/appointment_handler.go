package handlers

import (
	"net/http"
	"strconv"
	"time"

	"hospital-management/internal/models"
	"hospital-management/internal/service"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	appointmentService service.AppointmentService
}

func NewAppointmentHandler(appointmentService service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var appointmentReq models.AppointmentRequest
	if err := c.ShouldBindJSON(&appointmentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	createdAppointment, err := h.appointmentService.CreateAppointment(&appointmentReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdAppointment)
}

func (h *AppointmentHandler) GetAppointments(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	patientIDStr := c.Query("patient_id")
	dateStr := c.Query("date")

	var appointments []*models.Appointment
	var err error

	if doctorIDStr != "" {
		if doctorID, parseErr := strconv.ParseUint(doctorIDStr, 10, 32); parseErr == nil {
			appointments, err = h.appointmentService.GetAppointmentsByDoctor(uint(doctorID))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
			return
		}
	} else if patientIDStr != "" {
		if patientID, parseErr := strconv.ParseUint(patientIDStr, 10, 32); parseErr == nil {
			appointments, err = h.appointmentService.GetAppointmentsByPatient(uint(patientID))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
			return
		}
	} else if dateStr != "" {
		allAppointments, err := h.appointmentService.GetAllAppointments()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, appt := range allAppointments {
			if appt.DateTime.Format("2006-01-02") == dateStr {
				appointments = append(appointments, appt)
			}
		}
	} else {
		appointments, err = h.appointmentService.GetAllAppointments()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *AppointmentHandler) GetAppointmentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appointment, err := h.appointmentService.GetAppointmentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var appointmentUpdateReq models.AppointmentUpdateRequest
	if err := c.ShouldBindJSON(&appointmentUpdateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedAppointment, err := h.appointmentService.UpdateAppointment(uint(id), &appointmentUpdateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}

func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	if err := h.appointmentService.DeleteAppointment(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AppointmentHandler) GetDoctorSchedule(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("doctorId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	appointments, err := h.appointmentService.GetAppointmentsByDoctor(uint(doctorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetDate, parseErr := time.Parse("2006-01-02", dateStr)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var dayAppointments []*models.Appointment
	for _, appointment := range appointments {
		if appointment.DateTime.Format("2006-01-02") == targetDate.Format("2006-01-02") {
			dayAppointments = append(dayAppointments, appointment)
		}
	}

	c.JSON(http.StatusOK, dayAppointments)
}

func (h *AppointmentHandler) UpdateAppointmentStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
		Notes  string `json:"notes,omitempty"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if statusUpdate.Status != "confirmed" && statusUpdate.Status != "cancelled" && statusUpdate.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Use 'confirmed', 'cancelled', or 'completed'"})
		return
	}

	updateReq := models.AppointmentUpdateRequest{
		Status: statusUpdate.Status,
		Notes:  statusUpdate.Notes,
	}
	updatedAppointment, err := h.appointmentService.UpdateAppointment(uint(id), &updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}

func (h *AppointmentHandler) GetUpcomingAppointments(c *gin.Context) {
	allAppointments, err := h.appointmentService.GetAllAppointments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	var upcoming []*models.Appointment
	for _, appt := range allAppointments {
		if appt.DateTime.After(now) {
			upcoming = append(upcoming, appt)
		}
	}

	c.JSON(http.StatusOK, upcoming)
}

func (h *AppointmentHandler) RescheduleAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var rescheduleReq struct {
		DateTime string `json:"date_time"`
	}
	if err := c.ShouldBindJSON(&rescheduleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updateReq := models.AppointmentUpdateRequest{
		DateTime: rescheduleReq.DateTime,
	}
	updatedAppointment, err := h.appointmentService.UpdateAppointment(uint(id), &updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}
