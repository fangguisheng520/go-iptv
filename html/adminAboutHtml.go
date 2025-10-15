package html

import (
	"go-iptv/dto"
	"go-iptv/until"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

func About(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}

	var pageData = dto.AboutDto{
		LoginUser: username,
		Title:     "升级日志",
	}
	data, _ := os.ReadFile("/app/README.md")
	pageData.Content = string(data)
	re := regexp.MustCompile(`\./static`)
	pageData.Content = re.ReplaceAllString(pageData.Content, "/static")
	re = regexp.MustCompile(`\./ChangeLog.md`)
	pageData.Content = re.ReplaceAllString(pageData.Content, "/ChangeLog.md")

	c.HTML(http.StatusOK, "admin_about.html", pageData)
}
