package models

type IptvConfig struct {
	ID    int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"" json:"name"`
	Value string `gorm:"" json:"value"`
}

func (IptvConfig) TableName() string {
	return "iptv_config"
}
