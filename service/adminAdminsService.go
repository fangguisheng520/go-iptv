package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"net/url"
)

func Admins(params url.Values) dto.ReturnJsonDto {
	username := params.Get("username")
	oldPassword := params.Get("oldpassword")
	newpassword := params.Get("newpassword")
	newpassword2 := params.Get("newpassword_2")

	if username == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "用户名不能为空", Type: "danger"}
	}
	if oldPassword == "" && (newpassword != "" || newpassword2 != "") {
		return dto.ReturnJsonDto{Code: 0, Msg: "旧密码不能为空"}
	}

	if newpassword != newpassword2 && newpassword != "" && newpassword2 != "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "两次新密码不一致", Type: "danger"}
	}

	if oldPassword == newpassword && newpassword != "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "新密码不能与旧密码相同", Type: "danger"}
	}

	if !until.IsSafe(username) {
		return dto.ReturnJsonDto{Code: 0, Msg: "用户名不合法", Type: "danger"}
	}

	var adminData models.IptvAdmin
	dao.DB.Model(&models.IptvAdmin{}).Where("id = ?", 1).First(&adminData)
	if adminData.PassWord != until.HashPassword(oldPassword) {
		return dto.ReturnJsonDto{Code: 0, Msg: "旧密码错误", Type: "danger"}
	}

	dao.DB.Model(&models.IptvAdmin{}).Where("id = ?", 1).Updates(map[string]interface{}{
		"password": until.HashPassword(newpassword),
		"username": username,
	})

	// TODO
	return dto.ReturnJsonDto{Code: 1, Msg: "修改成功", Type: "success"}

	// TODO
}
