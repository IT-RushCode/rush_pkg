package auth

// Роль
type RoleRequestDTO struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Permissions []uint `json:"permissions,omitempty"`
}

// Роли
type RolesRequestDTO []RoleRequestDTO

type RoleResponseDTO struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type RolesResponseDTO []RoleResponseDTO

// Роль с привилегиями
type RoleWithPermissionsDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	Permissions PermissionsDTO `json:"permissions"`
}

// Роли с привилегиями
type RolesWithPermissionsDTO []RoleWithPermissionsDTO
