package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebHandler struct{}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

// LoginPage renders the login page
func (h *WebHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Hospital Management - Login",
	})
}

// ReceptionistDashboard renders the receptionist dashboard
func (h *WebHandler) ReceptionistDashboard(c *gin.Context) {
	username := getUsername(c)
	c.HTML(http.StatusOK, "receptionist/dashboard.html", gin.H{
		"title":    "Receptionist Dashboard",
		"username": username,
	})
}

// DoctorDashboard renders the doctor dashboard
func (h *WebHandler) DoctorDashboard(c *gin.Context) {
	username := getUsername(c)
	c.HTML(http.StatusOK, "doctor/dashboard.html", gin.H{
		"title":    "Doctor Dashboard",
		"username": username,
	})
}

// ReceptionistPatients renders the patient management page for receptionists
func (h *WebHandler) ReceptionistPatients(c *gin.Context) {
	username := getUsername(c)
	c.HTML(http.StatusOK, "receptionist/patients.html", gin.H{
		"title":    "Patient Management",
		"username": username,
	})
}

// DoctorPatients renders the patient records view for doctors
func (h *WebHandler) DoctorPatients(c *gin.Context) {
	username := getUsername(c)
	c.HTML(http.StatusOK, "doctor/patients.html", gin.H{
		"title":    "Patient Records",
		"username": username,
	})
}

// getUsername safely retrieves the username from context
func getUsername(c *gin.Context) string {
	if val, exists := c.Get("username"); exists {
		if username, ok := val.(string); ok {
			return username
		}
	}
	return "Unknown"
}
