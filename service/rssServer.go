package service

import (
	"go-iptv/dao"
	"go-iptv/models"
	"strings"
)

func GetRss(user string, id int64) string {
	var meal models.IptvMeals
	if err := dao.DB.Model(&models.IptvMeals{}).Where("id = ?", id).First(&meal).Error; err != nil {
		return ""
	}
	channelNameList := strings.Split(meal.Content, "_")

	var channels []models.IptvChannel
	dao.DB.Model(&models.IptvChannel{}).Where("category in (?)", channelNameList).Find(&channels)

	return ""
}
