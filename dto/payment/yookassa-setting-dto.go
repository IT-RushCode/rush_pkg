package payment

type YooKassaSettingDTO struct {
	ID        uint   `json:"id"`
	PointID   uint   `json:"pointId" validate:"required"`
	StoreID   string `json:"storeId" validate:"required_if_true=Status"`
	SecretKey string `json:"secretKey" validate:"required_if_true=Status"`
	IsTest    bool   `json:"isTest"`
	Status    bool   `json:"status"`
}
