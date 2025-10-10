package api

import (
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func GetRssUrl(c *gin.Context) {
	_, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	c.Request.ParseForm()
	params := c.Request.PostForm
	id := params.Get("id")

	c.JSON(200, service.GetRssUrl(id))
}

func GetTXTRss(c *gin.Context) {
	if token, ok := c.GetQuery("token"); !ok {
		c.String(200, "token 参数不存在")
	} else {
		if token == "" {
			c.String(200, "token 参数不存在")
		}
		c.String(200, service.GetRss(token))
	}
}
