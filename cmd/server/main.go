package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"hospital-management/internal/config"
	"hospital-management/internal/handlers"
	"hospital-management/internal/models"
	"hospital-management/internal/repository"
	"hospital-management/internal/service"
)

func main() {
	// Load environment variables from .env file if present
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database connection using GORM
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := db.AutoMigrate(
		&models.User{},
		&models.Patient{},
		&models.Appointment{},
	); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)

	authService := service.NewAuthService(userRepo)
	patientService := service.NewPatientService(patientRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo, patientRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	patientHandler := handlers.NewPatientHandler(patientService)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)

	// Setup Gin router and API routes
	router := gin.Default()
	api := router.Group("/api/v1")

	// Auth routes
	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)
	// api.POST("/logout", authHandler.Logout) // Uncomment if implemented

	// Patient routes
	// Patient routes
	patients := api.Group("/patients")
	patients.GET("", patientHandler.GetPatients)
	patients.POST("", patientHandler.CreatePatient)
	patients.GET(":id", patientHandler.GetPatientByID) // ✅ Fix here
	patients.PUT(":id", patientHandler.UpdatePatient)
	patients.DELETE(":id", patientHandler.DeletePatient)

	// Appointment routes
	appointments := api.Group("/appointments")
	appointments.GET("", appointmentHandler.GetAppointments)
	appointments.POST("", appointmentHandler.CreateAppointment)
	appointments.GET(":id", appointmentHandler.GetAppointmentByID)
	appointments.PUT(":id", appointmentHandler.UpdateAppointment)
	appointments.DELETE(":id", appointmentHandler.DeleteAppointment)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
