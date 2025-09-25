package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"time"
)

func AdminLogin(username, password string) dto.ReturnJsonDto {
	// 获取记住密码选项

	var userDb models.IptvAdmin

	err := dao.DB.Model(&models.IptvAdmin{}).Where("username = ? and password = ?", username, password).First(&userDb).Error
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "用户名或密码错误", Type: "danger"}
	}

	// 生成 JWT token，有效期为7天
	tokenString, err := until.GenerateJWT(username, 7*24*time.Hour)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "生成Token失败", Type: "danger"}
	}
	// 设置 Cookie，过期时间为7天
	// c.SetCookie("token", tokenString, 7*24*3600, "/", "", false, true)
	return dto.ReturnJsonDto{Code: 1, Msg: "登录成功", Type: "success", Data: tokenString}
}
