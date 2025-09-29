package models

type IptvCategory struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"unique;column:name" json:"name"`
	Enable       int    `gorm:"column:enable;default:1" json:"enable"`
	Type         string `gorm:"default:add;column:type" json:"type"`
	Url          string `gorm:"column:url" json:"url"`
	LatestTime   string `gorm:"column:latesttime" json:"latesttime"`
	AutoCategory int    `gorm:"column:autocategory" json:"autocategory"`
	Sort         int    `gorm:"column:sort" json:"sort"`
}

func (IptvCategory) TableName() string {
	return "iptv_category"
}
