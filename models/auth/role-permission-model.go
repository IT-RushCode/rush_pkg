package auth

// Привилегия роли
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;autoIncrement:false"`
	PermissionID uint `gorm:"primaryKey;autoIncrement:false"`
	Role         Role
	Permission   Permission
}

// Привилегии роли
type RolePermissions []RolePermission

func (RolePermission) TableName() string {
	return "RolePermissions"
}
