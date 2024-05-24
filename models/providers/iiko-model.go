package providers

import "time"

// Интеграция с iiko
type IikoIntegration struct {
	ID        uint `gorm:"primaryKey"`
	PointID   uint
	APIKey    string    `gorm:"type:varchar(255)"`
	APIURL    string    `gorm:"type:varchar(255)"`
	Username  string    `gorm:"type:varchar(100);default:null"`
	Password  string    `gorm:"type:varchar(100);default:null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type IikoIntegrations []IikoIntegration

func (IikoIntegration) TableName() string {
	return "IikoIntegrations"
}
