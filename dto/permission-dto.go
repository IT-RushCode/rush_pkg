package dto

type PermissionDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionsDTO []PermissionDTO
