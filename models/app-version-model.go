package models

import "gorm.io/gorm"

type AppVersion struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	AndroidVersion string `gorm:"column:android_version" json:"android_version"`
	IOSVersion     string `gorm:"column:ios_version" json:"ios_version"`
	AndroidURL     string `gorm:"column:android_url" json:"android_url"`
	IOSURL         string `gorm:"column:ios_url" json:"ios_url"`
	ReqIOS         *bool  `gorm:"column:req_ios" json:"req_ios"`
	ReqAndroid     *bool  `gorm:"column:req_android" json:"req_android"`
	TextUpIOS      string `gorm:"column:text_up_ios" json:"text_up_ios"`
	TextUpAndroid  string `gorm:"column:text_up_android" json:"text_up_android"`
	BaseModel
}

type AppVersions []AppVersion

func (AppVersion) TableName() string {
	return "AppVersions"
}

func (m *AppVersion) BeforeCreate(db *gorm.DB) (err error) {
	if err := CheckSequence(m.TableName(), db); err != nil {
		return err
	}
	return nil
}
