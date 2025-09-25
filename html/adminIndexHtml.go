package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"time"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

	username, ok := until.GetAuthName(c)
	if !ok {
		c.Redirect(302, "/admin/login")
		return
	}

	var pageData = dto.AdminIndexDto{
		LoginUser: username,
		Title:     "首页",
	}

	today := time.Now().Truncate(24 * time.Hour).Unix()

	dao.DB.Model(&models.IptvUser{}).Count(&pageData.UserTotal)
	dao.DB.Model(&models.IptvUser{}).Where("vpn > ?", 0).Count(&pageData.UserUNormal)
	dao.DB.Model(&models.IptvUser{}).Where("lasttime > ?", today).Count(&pageData.UserToday)
	dao.DB.Model(&models.IptvUser{}).Where("status > ? AND authortime > ?", 0, today).Count(&pageData.UserTodayAuth)
	dao.DB.Model(&models.IptvCategory{}).Where("autocategory IS NULL OR autocategory!='on'").Count(&pageData.ChannelTypeCount)
	dao.DB.Model(&models.IptvChannel{}).Count(&pageData.ChannelCount)
	dao.DB.Model(&models.IptvEpg{}).Count(&pageData.EpgCount)
	dao.DB.Model(&models.IptvMeals{}).Count(&pageData.MealsCount)

	var categoryList []models.IptvCategory
	dao.DB.Model(&models.IptvCategory{}).Where("type = ?", "add").Find(&categoryList)
	for i := range categoryList {
		var count int64
		var channelType dto.ChannelType
		dao.DB.Model(&models.IptvChannel{}).Where("category = ?", categoryList[i].Name).Count(&count)
		channelType.ChannelCount = count
		channelType.Num = int64(i + 1)
		channelType.Name = categoryList[i].Name
		pageData.ChannelTypeList = append(pageData.ChannelTypeList, channelType)
	}

	c.HTML(200, "admin_index.html", pageData)
}
