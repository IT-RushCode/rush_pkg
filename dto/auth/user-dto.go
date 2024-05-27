package auth

import (
	"time"
)

type UserRequestDTO struct {
	ID          uint      `json:"id" `
	FirstName   string    `json:"firstName" validate:"required"`
	LastName    string    `json:"lastName" validate:"required"`
	MiddleName  string    `json:"middleName" validate:"required"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber" validate:"required"`
	BirthDate   time.Time `json:"birthDate"`
	Status      bool      `json:"status" default:"false"`
	Avatar      string    `json:"avatar"`
	UserName    string    `json:"userName" validate:"required"`
	Password    string    `json:"password"`
	IsSuperUser bool      `json:"isSuperUser"  validate:"required"`
	Roles       []uint    `json:"roles" validate:"required"`
}

type UserResponseDTO struct {
	ID          uint      `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	MiddleName  string    `json:"middleName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	BirthDate   time.Time `json:"birthDate"`
	Status      bool      `json:"status"`
	Avatar      string    `json:"avatar"`
	UserName    string    `json:"userName"`
	IsSuperUser bool      `json:"isSuperUser"`
	CreatedAt   time.Time `json:"createAt"`
	UpdatedAt   time.Time `json:"updateAt"`

	Roles RolesWithPermissionsDTO `json:"roles"`
}

type UsersResponseDTO []UserResponseDTO
