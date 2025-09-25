package models

type IptvEpg struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string `gorm:"" json:"name"`
	Content string `gorm:"" json:"content"`
	Status  int    `gorm:"" json:"status"`
	Remarks string `gorm:"" json:"remarks"`
}

func (IptvEpg) TableName() string {
	return "iptv_epg"
}
