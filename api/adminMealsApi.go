package api

import (
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Meals(c *gin.Context) {
	_, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	c.Request.ParseForm()
	params := c.Request.PostForm
	var res dto.ReturnJsonDto

	for k := range params {
		switch k {
		case "change_status":
			res = service.MealsChangeStatus(params)
		case "editmeal":
			res = service.MealsEdit(params, 1)
		case "addmeal":
			res = service.MealsEdit(params, 0)
		case "delmeal":
			res = service.MealsDel(params)
		case "submitmeal":
			res = service.MealsSubmit(params)
		}
	}
	c.JSON(200, res)
}
