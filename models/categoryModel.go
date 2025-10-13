package models

type IptvCategory struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"unique;column:name" json:"name"`
	Enable       int64  `gorm:"column:enable;default:1" json:"enable"`
	Type         string `gorm:"default:add;column:type" json:"type"`
	Url          string `gorm:"column:url" json:"url"`
	LatestTime   string `gorm:"column:latesttime" json:"latesttime"`
	AutoCategory int64  `gorm:"column:autocategory" json:"autocategory"`
	Sort         int64  `gorm:"column:sort" json:"sort"`
}

func (IptvCategory) TableName() string {
	return "iptv_category"
}
