package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Exception(c *gin.Context) {
	idStr := c.DefaultQuery("id", "")

	if idStr != "" {
		if until.IsSafe(idStr) {
			var user models.IptvUser
			result := dao.DB.Where("name = ?", idStr).First(&user)
			if result.Error == nil {
				user.VPN += 1
				dao.DB.Save(&user)
				return
			}
		}
	}

	username, ok := until.GetAuthName(c)
	if !ok {
		c.Redirect(302, "/admin/login")
		return
	}
	var pageData = dto.AdminExceptionDto{
		LoginUser: username,
		Title:     "异常用户",
	}

	cfg := dao.GetConfig()
	pageData.MaxSameipUser = int(cfg.App.MaxSameIPUser)

	dao.DB.Model(&models.IptvUser{}).Select("status,name,model,vpn,idchange,marks,exp").Where("vpn>0 or idchange>0").Find(&pageData.Users)

	pageData.Users = until.CheckUserDay(pageData.Users)

	c.HTML(http.StatusOK, "admin_exception.html", pageData)
}
