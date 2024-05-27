package auth

// Роль
type RoleDTO struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      bool   `json:"status,omitempty" default:"true"`
	Permissions []uint `json:"permissions,omitempty"`
}

// Роли
type RolesDTO []RoleDTO

// Роль с привилегиями
type RoleWithPermissionsDTO struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      bool   `json:"status,omitempty" default:"true"`

	Permissions PermissionsDTO
}

// Роли с привилегиями
type RolesWithPermissionsDTO []RoleWithPermissionsDTO
