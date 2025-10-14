package models

type IptvEpg struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string `gorm:"column:name" json:"name"`
	Content string `gorm:"column:content" json:"content"`
	Status  int64  `gorm:"column:status" json:"status"`
	Remarks string `gorm:"column:remarks" json:"remarks"`
	Logo    string `gorm:"-" json:"logo"`
}

func (IptvEpg) TableName() string {
	return "iptv_epg"
}

type IptvEpgList struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Remarks     string `gorm:"column:remarks" json:"remarks"`
	Url         string `gorm:"column:url" json:"url"`
	LastTime    int64  `gorm:"column:lasttime" json:"lasttime"`
	LastTimeStr string `gorm:"-" json:"lasttimeStr"`
	Status      int64  `gorm:"column:status" json:"status"`
}

func (IptvEpgList) TableName() string {
	return "iptv_epg_list"
}
