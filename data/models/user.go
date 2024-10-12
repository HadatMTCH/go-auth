package models

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	FullName  string    `json:"full_name" db:"full_name"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Id    int    `json:"id"`
	Token string `json:"token"`
}

func NewUserResponse() *UserResponse {
	return &UserResponse{}
}

func (ur *UserResponse) ToResponse(user *User) {
	ur.ID = user.ID
	ur.Username = user.Username
	ur.Email = user.Email
	ur.FullName = user.FullName
	ur.IsAdmin = user.IsAdmin
	ur.CreatedAt = user.CreatedAt
	ur.UpdatedAt = user.UpdatedAt
}
