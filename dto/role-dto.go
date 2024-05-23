package dto

type RoleDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"default:true"`
	Permissions []uint `json:"permissions"`
}

type RolesDTO []RoleDTO

type RoleWithPermissionsDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"default:true"`

	Permissions PermissionsDTO
}

type RolesWithPermissionsDTO []RoleWithPermissionsDTO
