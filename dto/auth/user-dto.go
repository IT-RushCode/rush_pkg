package auth

import (
	"time"
)

type UserRequestDTO struct {
	ID          uint      `json:"id"`
	FirstName   string    `json:"firstName" validate:"required"`
	LastName    string    `json:"lastName" validate:"required"`
	MiddleName  string    `json:"middleName" validate:"required"`
	Email       string    `json:"email" validate:"email"`
	PhoneNumber string    `json:"phoneNumber" validate:"required"`
	BirthDate   time.Time `json:"birthDate"`
	Status      bool      `json:"status"`
	Avatar      string    `json:"avatar"`
	UserName    string    `json:"userName" validate:"required"`
	Password    string    `json:"password"`
	IsSuperUser bool      `json:"isSuperUser"  validate:"required"`
	Roles       []uint    `json:"roles" validate:"required"`
}

type UserResponseDTO struct {
	ID          uint      `json:"id,omitempty"`
	FirstName   string    `json:"firstName,omitempty"`
	LastName    string    `json:"lastName,omitempty"`
	MiddleName  string    `json:"middleName,omitempty"`
	Email       string    `json:"email,omitempty"`
	PhoneNumber string    `json:"phoneNumber,omitempty"`
	BirthDate   time.Time `json:"birthDate,omitempty"`
	Status      bool      `json:"status,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
	UserName    string    `json:"userName,omitempty"`
	IsSuperUser bool      `json:"isSuperUser,omitempty"`
	CreatedAt   time.Time `json:"createAt,omitempty"`
	UpdatedAt   time.Time `json:"updateAt,omitempty"`

	Roles RolesWithPermissionsDTO `json:"roles,omitempty"`
}

type UsersResponseDTO []UserResponseDTO
