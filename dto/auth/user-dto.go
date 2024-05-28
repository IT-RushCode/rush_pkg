package auth

import (
	"time"
)

type UserRequestDTO struct {
	ID          uint      `json:"id"`
	LastName    string    `json:"lastName" validate:"required"`
	FirstName   string    `json:"firstName" validate:"required"`
	MiddleName  string    `json:"middleName" validate:"required"`
	Email       string    `json:"email" validate:"email"`
	PhoneNumber string    `json:"phoneNumber" validate:"required,phone,len=12"`
	BirthDate   time.Time `json:"birthDate"`
	Status      bool      `json:"status"`
	Avatar      string    `json:"avatar"`
	UserName    string    `json:"userName" validate:"required"`
	Password    string    `json:"password"`
	IsSuperUser bool      `json:"isSuperUser"  validate:"required"`
	Roles       []uint    `json:"roles" validate:"required"`
}

type UserResponseDTO struct {
	ID           uint       `json:"id,omitempty"`
	LastName     string     `json:"lastName,omitempty"`
	FirstName    string     `json:"firstName,omitempty"`
	MiddleName   string     `json:"middleName,omitempty"`
	Email        string     `json:"email,omitempty"`
	PhoneNumber  string     `json:"phoneNumber,omitempty"`
	BirthDate    time.Time  `json:"birthDate,omitempty"`
	Status       bool       `json:"status,omitempty"`
	Avatar       string     `json:"avatar,omitempty"`
	UserName     string     `json:"userName,omitempty"`
	IsSuperUser  bool       `json:"isSuperUser,omitempty"`
	LastActivity time.Time  `json:"lastAcitvity,omitempty"`
	CreatedAt    *time.Time `json:"createAt,omitempty"`
	UpdatedAt    *time.Time `json:"updateAt,omitempty"`

	Roles RolesWithPermissionsDTO `json:"roles,omitempty"`
}

type UsersResponseDTO []UserResponseDTO

type UserPhoneDataDTO struct {
	ID          uint      `json:"id,omitempty"`
	PhoneNumber string    `json:"phoneNumber,omitempty" validate:"required,phone,len=12"`
	LastName    string    `json:"lastName,omitempty" validate:"required"`
	FirstName   string    `json:"firstName,omitempty" validate:"required"`
	MiddleName  string    `json:"middleName,omitempty"`
	Email       string    `json:"email,omitempty" validate:"email"`
	BirthDate   time.Time `json:"birthDate,omitempty"`
}
