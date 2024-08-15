package auth

// Привилегия роли
type RolePermission struct {
	RoleID       uint       `gorm:"primaryKey;index;not null;autoIncrement:false"`
	PermissionID uint       `gorm:"primaryKey;index;not null;autoIncrement:false"`
	Role         Role       `gorm:"foreignKey:RoleID"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
}

// Привилегии роли
type RolePermissions []RolePermission

func (RolePermission) TableName() string {
	return "RolePermissions"
}
