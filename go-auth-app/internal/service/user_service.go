package service

import (
	"errors"
	"fmt"

	"github.com/cliffdoyle/go-auth-app/internal/auth"
	"github.com/cliffdoyle/go-auth-app/internal/model"
	"github.com/cliffdoyle/go-auth-app/internal/repository"
	"gorm.io/gorm"
)

type SignupPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserService interface {
	RegisterUser(payload SignupPayload) (*model.User, error)
	CreateAdmin(payload SignupPayload) (*model.User, error)
	LoginUser(payload LoginPayload, jwtSecret string, jwtExpiry int) (string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(payload SignupPayload) (*model.User, error) {
	_, err := s.repo.FindUserByEmail(payload.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	// var userCount int64
	// // Simple check: if the first user (ID=1) doesn't exist, this is the first registration.
	// _, err = s.repo.FindUserByID(1)
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	userCount = 0
	// } else {
	// 	userCount = 1
	// }

	// role := model.UserRole
	// if userCount == 0 {
	// 	role = model.AdminRole
	// }

	user := &model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     model.UserRole,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

// NEW CreateAdmin function
func (s *userService) CreateAdmin(payload SignupPayload) (*model.User, error) {
	// Check if user already exists
	_, err := s.repo.FindUserByEmail(payload.Email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", payload.Email)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	// Explicitly create a user with the AdminRole
	user := &model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     model.AdminRole,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *userService) LoginUser(payload LoginPayload, jwtSecret string, jwtExpiry int) (string, error) {
	user, err := s.repo.FindUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if !auth.CheckPasswordHash(payload.Password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return auth.GenerateJWT(user, jwtSecret, jwtExpiry)
}
