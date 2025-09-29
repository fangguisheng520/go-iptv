package bootstrap

import (
	"go-iptv/dao"
	"go-iptv/models"
)

func InitDB() {
	dao.DB.AutoMigrate(&models.IptvAdmin{})
	dao.DB.AutoMigrate(&models.IptvUser{})
	dao.DB.AutoMigrate(&models.IptvChannel{})
	dao.DB.AutoMigrate(&models.IptvCategory{})
	dao.DB.AutoMigrate(&models.IptvEpg{})
	dao.DB.AutoMigrate(&models.IptvMeals{})
	dao.DB.AutoMigrate(&models.IptvMovie{})
}
