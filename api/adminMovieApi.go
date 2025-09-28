package api

import (
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Movie(c *gin.Context) {
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
			res = service.ChangeMovieStatus(params)
		case "delmovie":
			res = service.DeleteMovie(params)
		case "submitmovie":
			res = service.EditMovie(params)
		}
	}
	c.JSON(200, res)
}
