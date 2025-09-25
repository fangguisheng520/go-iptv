package models

type IptvChannel struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"" json:"name"`
	Url      string `gorm:"" json:"url"`
	Category string `gorm:"" json:"category"`
}

func (IptvChannel) TableName() string {
	return "iptv_channels"
}
