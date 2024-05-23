package dto

// Роль
type RoleDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"default:true"`
	Permissions []uint `json:"permissions"`
}

// Роли
type RolesDTO []RoleDTO

// Роль с привилегиями
type RoleWithPermissionsDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"default:true"`

	Permissions PermissionsDTO
}

// Роли с привилегиями
type RolesWithPermissionsDTO []RoleWithPermissionsDTO
