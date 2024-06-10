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
	Status      *bool     `json:"status"`
	AvatarUrl   string    `json:"avatarUrl"`
	UserName    string    `json:"userName" validate:"required"`
	Password    string    `json:"password"`
	IsPersonal  *bool     `json:"isPersonal"  validate:"required"`
	Roles       []uint    `json:"roles" validate:"required"`
}

type UserResponseDTO struct {
	ID           uint      `json:"id"`
	LastName     string    `json:"lastName"`
	FirstName    string    `json:"firstName"`
	MiddleName   string    `json:"middleName"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phoneNumber"`
	BirthDate    time.Time `json:"birthDate"`
	Status       *bool     `json:"status"`
	AvatarUrl    string    `json:"avatarUrl"`
	UserName     string    `json:"userName"`
	IsPersonal   *bool     `json:"IsPersonal"`
	LastActivity time.Time `json:"lastAcitvity"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	Roles RolesWithPermissionsDTO `json:"roles"`
}

type UsersResponseDTO []UserResponseDTO

type UserPhoneDataDTO struct {
	ID          uint       `json:"id"`
	PhoneNumber string     `json:"phoneNumber" validate:"required,phone,len=12"`
	LastName    string     `json:"lastName" validate:"required"`
	FirstName   string     `json:"firstName" validate:"required"`
	MiddleName  string     `json:"middleName"`
	Email       string     `json:"email" validate:"email"`
	BirthDate   *time.Time `json:"birthDate"`
	UserName    string     `json:"userName,omitempty"`
}
