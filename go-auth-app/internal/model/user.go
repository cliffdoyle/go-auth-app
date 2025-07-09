package model

import "gorm.io/gorm"

// Role defines the user roles in the system
type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

// User defines the user model for the database
type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" gorm:"unique" validate:"required,email"`
	Password string `json:"-"` // Omit from JSON responses
	Role     Role   `json:"role" gorm:"type:varchar(10)"`
}