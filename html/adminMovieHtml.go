package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Movie(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.Redirect(302, "/admin/login")
		return
	}
	var pageData = dto.AdminMovieDto{
		LoginUser: username,
		Title:     "点播管理",
	}

	dao.DB.Model(&models.IptvMovie{}).Find(&pageData.Movies)

	c.HTML(200, "admin_movie.html", pageData)
}
