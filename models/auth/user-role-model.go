package auth

// Роли пользователей
type UserRole struct {
	UserID uint `gorm:"primaryKey;index;not null;autoIncrement:false"`
	RoleID uint `gorm:"primaryKey;index;not null;autoIncrement:false"`
	User   User `gorm:"foreignKey:UserID"`
	Role   Role `gorm:"foreignKey:RoleID"`
}

type UserRoles []UserRole

func (UserRole) TableName() string {
	return "UserRoles"
}
