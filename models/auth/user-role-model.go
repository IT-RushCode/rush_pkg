	package auth

	// Роли пользователей
	type UserRole struct {
		UserID uint `gorm:"primaryKey;autoIncrement:false"`
		RoleID uint `gorm:"primaryKey;autoIncrement:false"`
		User   User `gorm:"foreignKey:UserID"`
		Role   Role `gorm:"foreignKey:RoleID"`
	}

	type UserRoles []UserRole

	func (UserRole) TableName() string {
		return "UserRoles"
	}
