package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Channels(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	var pageData = dto.AdminChannelsDto{
		LoginUser: username,
		Title:     "频道列表",
	}

	cfg := dao.GetConfig()
	autoUpdate := cfg.Channel.Auto
	if autoUpdate == 1 {
		pageData.AutoUpdate = true
	} else {
		pageData.AutoUpdate = false
	}

	pageData.UpdateInterval = cfg.Channel.Interval

	dao.DB.Model(&models.IptvCategory{}).Select("id, name, url, enable, autocategory, sort, type").Order("sort ASC").Find(&pageData.CategoryList)

	c.HTML(200, "admin_channels.html", pageData)
}
