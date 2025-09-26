package models

type IptvEpg struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string `gorm:"column:name" json:"name"`
	Content string `gorm:"column:content" json:"content"`
	Status  int    `gorm:"column:status" json:"status"`
	Remarks string `gorm:"column:remarks" json:"remarks"`
}

func (IptvEpg) TableName() string {
	return "iptv_epg"
}
