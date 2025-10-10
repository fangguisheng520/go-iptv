package models

type IptvChannel struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Url      string `gorm:"column:url" json:"url"`
	Category string `gorm:"column:category" json:"category"`
	Sort     int64  `gorm:"column:sort" json:"sort"`
}

func (IptvChannel) TableName() string {
	return "iptv_channels"
}
