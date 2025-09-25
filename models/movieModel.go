package models

type IptvMovie struct {
	ID    int64  `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"name" json:"name"`
	Api   string `gorm:"api" json:"api"`
	State int64  `gorm:"state"`
}

func (IptvMovie) TableName() string {
	return "iptv_movie"
}
