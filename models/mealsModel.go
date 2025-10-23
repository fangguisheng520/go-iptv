package models

type IptvMeals struct {
	ID      int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name    string `gorm:"column:name" json:"name"`
	Content string `gorm:"column:content" json:"content"`
	Status  int64  `gorm:"column:status" json:"status"`
}

func (IptvMeals) TableName() string {
	return "iptv_meals"
}

type IptvMealsShow struct {
	ID      int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name    string `gorm:"column:name" json:"name"`
	Content string `gorm:"column:content" json:"content"`
	Status  int64  `gorm:"column:status" json:"status"`
	CaName  string `gorm:"-" json:"caname"`
}

func (IptvMealsShow) TableName() string {
	return "iptv_meals"
}
