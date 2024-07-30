package auth

// Привилегия
type PermissionDTO struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Привилегии
type PermissionsDTO []PermissionDTO
