package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"strings"

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

	dao.DB.Model(&models.IptvCategoryList{}).Find(&pageData.CategoryList)
	dao.DB.Model(&models.IptvCategory{}).Order("sort ASC").Find(&pageData.Categorys)
	dao.DB.Model(&models.IptvEpg{}).Where("status = 1").Find(&pageData.Epgs)

	logoList := until.GetLogos()
	for i, v := range pageData.Epgs {
		epgName := strings.SplitN(v.Name, "-", 2)[1]
		for _, logo := range logoList {
			logoName := strings.Split(logo, ".")[0]
			if strings.EqualFold(epgName, logoName) {
				pageData.Epgs[i].Logo = "/logo/" + logo
			}
		}
	}

	c.HTML(200, "admin_channels.html", pageData)
}
