package api

import (
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Admins(c *gin.Context) {
	_, ok := until.GetAuthName(c)
	if !ok {
		c.Redirect(302, "/admin/login")
		return
	}
	c.Request.ParseForm()
	params := c.Request.PostForm

	c.JSON(200, service.Admins(params))
}
