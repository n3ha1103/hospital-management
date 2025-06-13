// internal/service/auth_service.go
package service

import (
	"fmt"
	"hospital-management/internal/models"
	"hospital-management/internal/repository"
	"hospital-management/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	Register(req *models.RegisterRequest) (*models.User, error)
	ValidateToken(tokenString string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email or username
	var user *models.User
	var err error

	// Try to find user by email first, then by username if email lookup fails
	user, err = s.userRepo.GetByEmail(req.Email)
	if err != nil {
		// If not found by email, try username (assuming req.Email could be username)
		user, err = s.userRepo.GetByUsername(req.Email)
		if err != nil {
			return nil, fmt.Errorf("invalid credentials")
		}
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token with 2 parameters (userID and role)
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Username:  user.Name,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			FullName:  user.GetFullName(),
			Email:     user.Email,
			Role:      user.Role,
		},
	}, nil
}

func (s *authService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if user already exists by email
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email already exists")
	}

	// Check if username already exists (using Name field)
	if req.Name != "" {
		existingUser, _ = s.userRepo.GetByUsername(req.Name)
		if existingUser != nil {
			return nil, fmt.Errorf("username already exists")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      req.Role,
		FirstName: "", // Not provided in RegisterRequest
		LastName:  "", // Not provided in RegisterRequest
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Don't return password in response
	createdUser.Password = ""

	return createdUser, nil
}

func (s *authService) ValidateToken(tokenString string) (*models.User, error) {
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user from database - convert uint to int
	user, err := s.userRepo.GetByID(uint(claims.UserID))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Don't return password
	user.Password = ""

	return user, nil
}
