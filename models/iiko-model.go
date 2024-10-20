package models

import "gorm.io/gorm"

// Интеграция с iiko
type IikoIntegration struct {
	ID       uint   `gorm:"primaryKey"`
	PointID  uint   `gorm:"index;not null"`
	APIKey   string `gorm:"type:varchar(255)"`
	APIURL   string `gorm:"type:varchar(255)"`
	Username string `gorm:"type:varchar(100);default:null"`
	Password string `gorm:"type:varchar(100);default:null"`
	BaseModel
}

type IikoIntegrations []IikoIntegration

func (IikoIntegration) TableName() string {
	return "IikoIntegrations"
}

func (a *IikoIntegration) BeforeCreate(tx *gorm.DB) (err error) {
	if err := CheckSequence("IikoIntegrations", tx); err != nil {
		return err
	}
	return nil
}
