package html

import (
	"go-iptv/dto"
	"go-iptv/until"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func About(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.Redirect(302, "/admin/login")
		return
	}

	var pageData = dto.AboutDto{
		LoginUser: username,
		Title:     "升级日志",
	}
	data, _ := os.ReadFile("/app/README.md")
	pageData.Content = string(data)

	c.HTML(http.StatusOK, "admin_about.html", pageData)
}
