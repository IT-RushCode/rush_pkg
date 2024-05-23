package dto

type PermissionDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"default:true"`
}

type PermissionsDTO []PermissionDTO
