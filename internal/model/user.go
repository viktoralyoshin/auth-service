package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	Id           uuid.UUID
	Username     string
	Email        string
	Password     string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
