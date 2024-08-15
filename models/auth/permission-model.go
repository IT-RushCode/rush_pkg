package auth

// Привилегия
type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description string `gorm:"type:varchar(255);default:null"`
}

// Привилегии
type Permissions []Permission

func (Permission) TableName() string {
	return "Permissions"
}
