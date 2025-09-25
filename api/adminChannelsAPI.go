package api

import (
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Channels(c *gin.Context) {
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
		case "update_interval":
			res = service.UpdateInterval(params)
		case "updatelist":
			res = service.UpdateList(params)
		case "addlist":
			res = service.AddList(params)
		case "dellist":
			res = service.DelList(params)
		case "getchannels":
			c.String(200, service.AdminGetChannels(params))
			return
		case "forbiddenchannels":
			res = service.ForbiddenChannels(params)
		case "submit_addtype":
			res = service.SubmitAddType(params)
		case "submit_deltype":
			res = service.SubmitDelType(params)
		case "submit_modifytype":
			res = service.SubmitModifyType(params)
		case "submit_moveup":
			res = service.SubmitMoveUp(params)
		case "submit_movedown":
			res = service.SubmitMoveDown(params)
		case "submit_movetop":
			res = service.SubmitMoveTop(params)
		case "submitsave":
			res = service.SubmitSave(params)
		}
	}

	c.JSON(200, res)
}
