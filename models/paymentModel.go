package models

type IptvPayment struct {
	UserID  int64  `gorm:"" json:"userid"`
	OrderID string `gorm:"" json:"order_id"`
	Meal    int    `gorm:"" json:"meal"`
	Days    int    `gorm:"" json:"days"`
	Status  int    `gorm:"" json:"status"`
}

func (IptvPayment) TableName() string {
	return "iptv_payment"
}
