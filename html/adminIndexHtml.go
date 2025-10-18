package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}

	var pageData = dto.AdminIndexDto{
		LoginUser: username,
		Title:     "首页",
	}

	today := time.Now().Truncate(24 * time.Hour).Unix()

	dao.DB.Model(&models.IptvUser{}).Count(&pageData.UserTotal)
	dao.DB.Model(&models.IptvUser{}).Where("lasttime > ?", today).Count(&pageData.UserToday)
	dao.DB.Model(&models.IptvCategory{}).Where("autocategory IS NULL OR autocategory!='on'").Count(&pageData.ChannelTypeCount)
	dao.DB.Model(&models.IptvChannel{}).Count(&pageData.ChannelCount)
	dao.DB.Model(&models.IptvEpg{}).Count(&pageData.EpgCount)
	dao.DB.Model(&models.IptvMeals{}).Count(&pageData.MealsCount)

	var categoryList []models.IptvCategory
	dao.DB.Model(&models.IptvCategory{}).Where("type != ?", "import").Find(&categoryList)
	for i := range categoryList {
		var count int64
		var channelType dto.ChannelType
		if categoryList[i].Sort == -2 {
			text := until.GetCCTVChannelList(false)
			text = strings.TrimSpace(text) // 去掉结尾多余换行
			parts := strings.Split(text, "\n")
			if len(parts) == 1 && parts[0] == "" {
				channelType.ChannelCount = 0
			} else {
				channelType.ChannelCount = int64(len(parts))
			}
			channelType.Num = int64(i + 1)
			channelType.Name = categoryList[i].Name
			channelType.RawCount = categoryList[i].Rawcount
			pageData.ChannelTypeList = append(pageData.ChannelTypeList, channelType)
			continue
		} else if categoryList[i].Sort == -1 {
			text := until.GetProvinceChannelList(false)
			text = strings.TrimSpace(text) // 去掉结尾多余换行
			parts := strings.Split(text, "\n")
			if len(parts) == 1 && parts[0] == "" {
				channelType.ChannelCount = 0
			} else {
				channelType.ChannelCount = int64(len(parts))
			}

			channelType.Num = int64(i + 1)
			channelType.Name = categoryList[i].Name
			channelType.RawCount = categoryList[i].Rawcount
			pageData.ChannelTypeList = append(pageData.ChannelTypeList, channelType)
			continue
		}

		dao.DB.Model(&models.IptvChannel{}).Where("category = ?", categoryList[i].Name).Count(&count)
		channelType.ChannelCount = count
		channelType.Num = int64(i + 1)
		channelType.Name = categoryList[i].Name
		channelType.RawCount = categoryList[i].Rawcount
		pageData.ChannelTypeList = append(pageData.ChannelTypeList, channelType)
	}

	c.HTML(200, "admin_index.html", pageData)
}
