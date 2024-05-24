package auth

// Роли пользователей
type UserRole struct {
	UserID uint `gorm:"primaryKey;autoIncrement:false"`
	RoleID uint `gorm:"primaryKey;autoIncrement:false"`
	User   User
	Role   Role
}

type UserRoles []UserRole

func (UserRole) TableName() string {
	return "UserRoles"
}
