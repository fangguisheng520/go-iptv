package api

import (
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Exception(c *gin.Context) {
	_, ok := until.GetAuthName(c)
	if !ok {
		c.Redirect(302, "/admin/login")
		return
	}
	c.Request.ParseForm()
	params := c.Request.PostForm

	var res dto.ReturnJsonDto

	for k := range params {
		switch k {
		case "submitsameip_user":
			res = service.SetSameIpUser(params)
		case "clearvpn":
			res = service.ClearVpn()
		case "clearidchange":
			res = service.ClearIdChange()
		case "stopuse":
			res = service.StopUse(params)
		case "startuse":
			res = service.StartUse(params)
		}
	}

	c.JSON(200, res)
}
