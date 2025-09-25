package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"net/url"
	"strconv"
)

func SetSameIpUser(params url.Values) dto.ReturnJsonDto {

	sameIpUser := params.Get("sameip_user")

	if sameIpUser == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入值", Type: "danger"}
	}

	num, err := strconv.ParseInt(sameIpUser, 10, 64)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "请输入数字", Type: "danger"}
	}
	cfg := dao.GetConfig()

	cfg.App.MaxSameIPUser = num
	dao.SetConfig(cfg)

	return dto.ReturnJsonDto{Code: 1, Msg: "设置成功", Type: "success"}
}

func ClearVpn() dto.ReturnJsonDto {
	dao.DB.Model(&models.IptvUser{}).Where("vpn > ?", 0).Update("vpn", 0)
	return dto.ReturnJsonDto{Code: 1, Msg: "抓包记录已清空", Type: "success"}
}

func ClearIdChange() dto.ReturnJsonDto {
	dao.DB.Model(&models.IptvUser{}).Where("idchange = ?", 0).Update("idchange", 0)
	return dto.ReturnJsonDto{Code: 1, Msg: "设备ID更换记录已清空", Type: "success"}

}

func StopUse(params url.Values) dto.ReturnJsonDto {
	name := params.Get("name")
	dao.DB.Model(&models.IptvUser{}).Where("name = ?", name).Update("status", 0)
	return dto.ReturnJsonDto{Code: 1, Msg: "已禁用", Type: "success"}
}

func StartUse(params url.Values) dto.ReturnJsonDto {
	name := params.Get("name")
	dao.DB.Model(&models.IptvUser{}).Where("name = ? and status=0", name).Update("status", 1)
	return dto.ReturnJsonDto{Code: 1, Msg: "已启用", Type: "success"}
}
