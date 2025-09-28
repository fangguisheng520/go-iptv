package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Meals(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	var pageData = dto.AdminMealsDto{
		LoginUser: username,
		Title:     "套餐管理",
	}

	dao.DB.Model(&models.IptvMeals{}).Find(&pageData.Meals)
	dao.DB.Model(&models.IptvCategory{}).Where("type = ?", "add").Count(&pageData.ChannelNum)

	c.HTML(200, "admin_meals.html", pageData)
}
