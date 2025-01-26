package models

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID          int64     `json:"id" db:"id" swaggerignore:"true"`
	Username    string    `json:"username" db:"username" example:"johndoe"`
	Email       string    `json:"email" db:"email" example:"john@example.com"`
	Password    string    `json:"password" db:"password" example:"securepassword123"`
	FullName    string    `json:"full_name" db:"full_name" example:"John Doe"`
	IsAdmin     bool      `json:"is_admin" db:"is_admin" swaggerignore:"true"`
	IsSuspended bool      `json:"is_suspended" db:"is_suspended" swaggerignore:"true"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" swaggerignore:"true"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" swaggerignore:"true"`
}

type UserResponse struct {
	ID        int64     `json:"id" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	FullName  string    `json:"full_name" example:"John Doe"`
	IsAdmin   bool      `json:"is_admin" example:"false"`
	CreatedAt time.Time `json:"created_at" example:"2023-06-15T14:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-06-15T14:30:00Z"`
}

func (l *User) Validate() error {
	l.Email = strings.TrimSpace(l.Email)
	l.Password = strings.TrimSpace(l.Password)
	l.Username = strings.TrimSpace(l.Username)

	if l.Email == "" {
		return errors.New("email is required")
	}

	if l.Password == "" {
		return errors.New("password is required")
	}

	if l.Username == "" {
		return errors.New("username is required")
	}

	return nil
}

func (l *User) ValidateLogin() error {
	l.Email = strings.TrimSpace(l.Email)
	l.Password = strings.TrimSpace(l.Password)

	if l.Email == "" {
		return errors.New("email is required")
	}

	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type LoginResponse struct {
	Id    int    `json:"id" example:"1"`
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

func NewUserResponse() *UserResponse {
	return &UserResponse{}
}

func (ur *UserResponse) ToResponse(user *User) *UserResponse {
	ur.ID = user.ID
	ur.Username = user.Username
	ur.Email = user.Email
	ur.FullName = user.FullName
	ur.IsAdmin = user.IsAdmin
	ur.CreatedAt = user.CreatedAt
	ur.UpdatedAt = user.UpdatedAt
	return ur
}
