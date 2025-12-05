package models

import "gorm.io/gorm"

type AppVersion struct {
	ID             uint   `gorm:"primaryKey;comment:Первичный ключ версии" json:"id"`
	AndroidVersion string `gorm:"column:android_version;comment:Версия Android" json:"android_version"`
	IOSVersion     string `gorm:"column:ios_version;comment:Версия iOS" json:"ios_version"`
	AndroidURL     string `gorm:"column:android_url;comment:Ссылка на Android приложение" json:"android_url"`
	IOSURL         string `gorm:"column:ios_url;comment:Ссылка на iOS приложение" json:"ios_url"`
	ReqIOS         *bool  `gorm:"column:req_ios;comment:Обязательное обновление для iOS" json:"req_ios"`
	ReqAndroid     *bool  `gorm:"column:req_android;comment:Обязательное обновление для Android" json:"req_android"`
	TextUpIOS      string `gorm:"column:text_up_ios;comment:Текст окна обновления для iOS" json:"text_up_ios"`
	TextUpAndroid  string `gorm:"column:text_up_android;comment:Текст окна обновления для Android" json:"text_up_android"`
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
