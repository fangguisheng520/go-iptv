package models

type IptvCategory struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"unique;column:name" json:"name"`
	Enable       int64  `gorm:"column:enable;default:1" json:"enable"`
	Type         string `gorm:"default:hand;column:type" json:"type"`
	Url          string `gorm:"column:url" json:"url"`
	UA           string `gorm:"column:ua" json:"ua"`
	LatestTime   string `gorm:"column:latesttime" json:"latesttime"`
	AutoCategory int64  `gorm:"column:autocategory" json:"autocategory"`
	Repeat       int64  `gorm:"column:repeat" json:"repeat"`
	Sort         int64  `gorm:"column:sort" json:"sort"`
	Rawcount     int64  `gorm:"column:rawcount;default:0" json:"rawcount"`
}

func (IptvCategory) TableName() string {
	return "iptv_category"
}
