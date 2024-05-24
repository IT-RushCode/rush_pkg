package providers

// Данные гугл/яндекс карты
type MapData struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`
	Coordinates []byte `gorm:"type:jsonb"`
	Latitude    float64
	Longitude   float64
}

type MapDatas []MapData

func (MapData) TableName() string {
	return "MapData"
}
