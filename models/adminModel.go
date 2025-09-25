package models

type IptvAdmin struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserName string `gorm:"column:username" json:"username"`
	PassWord string `gorm:"column:password" json:"password"`
}

func (IptvAdmin) TableName() string {
	return "iptv_admin"
}
