package payment

type YooKassaSettingDTO struct {
	ID        uint   `json:"id"`
	PointID   uint   `json:"pointId" validate:"required"`
	StoreID   string `json:"storeId" validate:"required"`
	SecretKey string `json:"secretKey" validate:"required"`
	IsTest    bool   `json:"isTest"`
}
