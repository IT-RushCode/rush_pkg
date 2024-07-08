package auth

// Привилегия
type PermissionDTO struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      *bool  `json:"status,omitempty" default:"true"`
}

// Привилегии
type PermissionsDTO []PermissionDTO
