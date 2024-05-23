package dto

type RoleDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions []uint `json:"permissions"`
}

type RolesDTO []RoleDTO

type RoleWithPermissionsDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions PermissionsDTO
}

type RolesWithPermissionsDTO []RoleWithPermissionsDTO
