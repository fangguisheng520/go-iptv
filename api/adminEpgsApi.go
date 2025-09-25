package api

import (
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Epgs(c *gin.Context) {
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
		case "editepg":
			res = service.GetEpgData(params)
		case "edit_save_epg":
			res = service.SaveEpg(params, 1)
		case "add_save_epg":
			res = service.SaveEpg(params, 0)
		case "change_status":
			res = service.ChangeStatus(params)
		case "delepg":
			res = service.DeleteEpg(params)
		case "bindchannel":
			res = service.BindChannel()
		case "clearbind":
			res = service.ClearBind()
			// case "saveepgapi_chk":
			// 	res = service.SaveEpgApi(params)
		}

	}
	c.JSON(200, res)
}
