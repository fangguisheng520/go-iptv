package api

import (
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Fenlei(c *gin.Context) {
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
		case "update_interval":
			res = service.UpdateInterval(params)
		case "updatelist":
			res = service.UpdateList(params)
		case "updatelistall":
			res = service.UpdateListAll()
		case "addlist":
			res = service.AddList(params)
		case "dellist":
			res = service.DelList_Fenlei(params)
		case "getchannels_pindao":
			res = service.CaGetChannels_pindao(params)
		case "delca":
			res = service.DelCa_Fenlei(params)
		case "moveup":
			res = service.SubmitMoveUp_Fenlei(params)
		case "movedown":
			res = service.SubmitMoveDown_Fenlei(params)
		case "movetop":
			res = service.SubmitMoveTop_Fenlei(params)
		case "saveChannels":
			res = service.SubmitSave_Fenlei(params)
		case "saveChannelsOne_fenlei":
			res = service.SaveChannelsOne_Fenlei(params)
		case "fenleiStatus":
			res = service.FenleiChangeStatus(params)
		case "categoryListStatus":
			res = service.CategoryListChangeStatus(params)
		case "channelsStatus":
			res = service.ChannelsChangeStatus_Fenlei(params)
		case "saveCa":
			res = service.SaveFenlei(params)
		case "addPindao":
			res = service.AddPindao(params)
		}
	}

	c.JSON(200, res)
}

//func UploadPayList(c *gin.Context) {
//	_, ok := until.GetAuthName(c)
//	if !ok {
//		c.JSON(200, dto.NewAdminRedirectDto())
//		return
//	}
//	c.JSON(200, service.UploadPayList(c))
//}
