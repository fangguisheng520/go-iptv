package html

import (
	"go-iptv/bootstrap"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func Client(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	cfg := dao.GetConfig()

	var pageData = dto.AdminClientDto{
		LoginUser:   username,
		Title:       "客户端设置",
		Build:       cfg.Build,
		App:         cfg.App,
		Tips:        cfg.Tips,
		UpSize:      until.GetFileSize("./app/" + cfg.Build.Name + ".apk"),
		BuildStatus: bootstrap.GetBuildStatus(), // 获取APK编译状态
	}

	if until.Exists("images/icon/icon.png") {
		pageData.IconUrl = "/icon/icon.png"
	}

	pageData.BjUrl, _ = until.GetPngFileNames("images/bj")

	c.HTML(200, "admin_client.html", pageData)
}
