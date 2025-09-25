package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Notice(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.Redirect(302, "/admin/login")
		return
	}

	cfg := dao.GetConfig()

	var pageData = dto.AdminNoticeDto{
		LoginUser: username,
		Title:     "系统公告",
		Ad:        cfg.Ad,
	}

	c.HTML(200, "admin_notice.html", pageData)
}
