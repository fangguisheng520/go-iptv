package models

type IptvChannel struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Url      string `gorm:"column:url" json:"url"`
	Category string `gorm:"column:category" json:"category"`
}

func (IptvChannel) TableName() string {
	return "iptv_channels"
}
