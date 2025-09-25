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
		c.Redirect(302, "/admin/login")
		return
	}
	var pageData = dto.AdminMealsDto{
		LoginUser: username,
		Title:     "套餐管理",
	}

	dao.DB.Model(&models.IptvMeals{}).Find(&pageData.Meals)

	c.HTML(200, "admin_meals.html", pageData)
}
