package auth

import (
	"time"
)

type UserRequestDTO struct {
	ID          uint      `json:"id,omitempty"`
	LastName    string    `json:"lastName,omitempty" validate:"required"`
	FirstName   string    `json:"firstName,omitempty" validate:"required"`
	MiddleName  string    `json:"middleName,omitempty" validate:"required"`
	PhoneNumber string    `json:"phoneNumber,omitempty" validate:"required,phone,len=12"`
	BirthDate   time.Time `json:"birthDate,omitempty"`
	Status      bool      `json:"status,omitempty"`
	AvatarUrl   string    `json:"avatarUrl,omitempty"`
	Email       string    `json:"email,omitempty" validate:"required,email"`
	UserName    string    `json:"userName,omitempty"`
	Password    string    `json:"password,omitempty"`
	IsStaff     bool      `json:"isStaff,omitempty"  validate:"required"`
	Roles       []uint    `json:"roles,omitempty" validate:"required"`
}

type UserResponseDTO struct {
	ID                      uint      `json:"id"`
	LastName                string    `json:"lastName"`
	FirstName               string    `json:"firstName"`
	MiddleName              string    `json:"middleName"`
	Email                   string    `json:"email"`
	EmailConfirmed          *bool     `json:"emailConfirmed"`
	PhoneNumber             string    `json:"phoneNumber"`
	PhoneConfirmed          *bool     `json:"phoneConfirmed"`
	BirthDate               time.Time `json:"birthDate"`
	Status                  *bool     `json:"status"`
	AvatarUrl               string    `json:"avatarUrl"`
	UserName                string    `json:"userName"`
	IsStaff                 *bool     `json:"IsStaff"`
	LastActivity            time.Time `json:"lastAcitvity"`
	ChangePasswordWhenLogin *bool     `json:"changePasswordWhenLogin"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	Roles RolesResponseDTO `json:"roles"`
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
}

// DTO для изменения пароля пользователя
type ChangePasswordRequestDTO struct {
	OldPassword     string `json:"oldPassword" validate:"required,min=8"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}
