package api

import (
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Client(c *gin.Context) {
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
		case "deleteIcon":
			res = service.DeleteFile(params, "icon")
		case "deleteBj":
			res = service.DeleteFile(params, "bj")
		case "decoder":
			res = service.DecoderSelect(params)
		case "buffTimeOut":
			res = service.SetBuffTimeOut(params)
		case "submittrialdays":
			res = service.SetTrialDays(params)
		case "submitappinfo":
			res = service.SetAppInfo(params)
		case "submittipset":
			res = service.SetTipSet(params)
		}
	}
	c.JSON(200, res)
}

func ClientUploadIcon(c *gin.Context) {
	_, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	c.JSON(200, service.UploadFile(c, "icon"))
}

func ClientUploadBj(c *gin.Context) {
	_, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	c.JSON(200, service.UploadFile(c, "bj"))
}
