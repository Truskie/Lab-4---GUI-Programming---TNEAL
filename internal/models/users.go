package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required,min=3,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	FirstName string    `json:"first_name,omitempty" validate:"max=50"`
	LastName  string    `json:"last_name,omitempty" validate:"max=50"`
	IsActive  bool      `json:"is_active"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Validation method
func (u *User) Validate() error {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	// Username validation
	if u.Username == "" {
		return errors.New("username required")
	}
	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(u.Username) > 50 {
		return errors.New("username canot be longer than 50 characters")
	}

	// setting username format
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(u.Username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	// validating for email
	if u.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	if len(u.Email) > 100 {
		return errors.New("email cannot be longer than 100 characters")
	}

	// validating for names
	if len(u.FirstName) > 50 {
		return errors.New("first name cannot be longer than 50 characters")
	}
	if len(u.LastName) > 50 {
		return errors.New("last name cannot be longer than 50 characters")
	}

	return nil
}

// PATCH
func (u *User) ValidatePartial() error {
	if u.Username != "" {
		if len(u.Username) < 3 {
			return errors.New("username must be at least 3 characters")
		}
		if len(u.Username) > 50 {
			return errors.New("username cannot be longer than 50 characters")
		}
		usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		if !usernameRegex.MatchString(u.Username) {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}

	if u.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(u.Email) {
			return errors.New("invalid email format")
		}
		if len(u.Email) > 100 {
			return errors.New("email cannot be longer than 100 characters")
		}
	}

	if len(u.FirstName) > 50 {
		return errors.New("first name cannot be longer than 50 characters")
	}
	if len(u.LastName) > 50 {
		return errors.New("last name cannot be longer than 50 characters")
	}

	return nil
}
