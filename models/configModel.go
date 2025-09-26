package models

type IptvConfig struct {
	ID    int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"column:name" json:"name"`
	Value string `gorm:"column:value" json:"value"`
}

func (IptvConfig) TableName() string {
	return "iptv_config"
}
