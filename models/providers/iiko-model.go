package providers

import rpBase "github.com/IT-RushCode/rush_pkg/models"

// Интеграция с iiko
type IikoIntegration struct {
	ID       uint   `gorm:"primaryKey"`
	PointID  uint   `gorm:"index;not null"`
	APIKey   string `gorm:"type:varchar(255)"`
	APIURL   string `gorm:"type:varchar(255)"`
	Username string `gorm:"type:varchar(100);default:null"`
	Password string `gorm:"type:varchar(100);default:null"`
	rpBase.BaseModel
}

type IikoIntegrations []IikoIntegration

func (IikoIntegration) TableName() string {
	return "IikoIntegrations"
}
