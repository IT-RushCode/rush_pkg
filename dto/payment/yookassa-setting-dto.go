package payment

type YooKassaSettingDTO struct {
	ID             uint   `json:"id"`
	PointID        uint   `json:"pointId"`
	StoreID        string `json:"storeId"`
	SecretKey      string `json:"secretKey"`
	ApplicationKey string `json:"applicationKey"`
}
