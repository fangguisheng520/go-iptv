package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Admins(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.Redirect(302, "/admin/login")
		return
	}

	var pageData = dto.AdminsDto{
		LoginUser: username,
		Title:     "管理员设置",
	}

	dao.DB.Model(&models.IptvAdmin{}).Where("id = 1").First(&pageData.Admins)

	c.HTML(http.StatusOK, "admin_admins.html", pageData)
}
