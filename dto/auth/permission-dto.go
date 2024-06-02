package auth

// Привилегия
type PermissionDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      *bool  `json:"status" default:"true"`
}

// Привилегии
type PermissionsDTO []PermissionDTO
