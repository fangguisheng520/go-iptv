package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/until"
	"net/url"
	"os"
	"strconv"
)

func ImportLicense(params url.Values) dto.ReturnJsonDto {
	lickey := params.Get("lickey")
	if lickey == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "lickey不能为空", Type: "danger"}
	}

	oldLickey := until.ReadFile("/config/license.lic")

	err := os.WriteFile("/config/license.lic", []byte(lickey), 0644)
	if err != nil {
		os.WriteFile("/config/license.lic", []byte(oldLickey), 0644)
		return dto.ReturnJsonDto{Code: 0, Msg: "文件写入失败: " + err.Error(), Type: "danger"}
	}

	res, err := dao.WS.SendWS(dao.Request{Action: "reloadLic"})
	if err != nil {
		os.WriteFile("/config/license.lic", []byte(oldLickey), 0644)
		return dto.ReturnJsonDto{Code: 0, Msg: "授权失败: " + err.Error(), Type: "danger"}
	} else if res.Code == 1 {
		//授权成功
	} else if res.Code != 1 {
		os.WriteFile("/config/license.lic", []byte(oldLickey), 0644)
		return dto.ReturnJsonDto{Code: 0, Msg: "授权失败: " + res.Msg, Type: "danger"}
	}
	return dto.ReturnJsonDto{Code: 1, Msg: "授权成功", Type: "success"}
}

func Proxy(params url.Values) dto.ReturnJsonDto {
	cfg := dao.GetConfig()
	if dao.Lic.Tpye == 0 {
		cfg.Proxy.Status = 0

		dao.SetConfig(cfg)
		dao.WS.SendWS(dao.Request{Action: "stopProxy"})
		return dto.ReturnJsonDto{Code: 0, Msg: "未授权", Type: "danger"}
	}

	port := params.Get("port")
	proxy := params.Get("proxy")

	if proxy == "1" || proxy == "true" || proxy == "on" {
		if port == "" {
			return dto.ReturnJsonDto{Code: 0, Msg: "port不能为空", Type: "danger"}
		}

		portInt64, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "port为数字", Type: "danger"}
		}
		if portInt64 < 82 || portInt64 > 65535 {
			return dto.ReturnJsonDto{Code: 0, Msg: "port为82-65535", Type: "danger"}
		}
		res, err := dao.WS.SendWS(dao.Request{Action: "startProxy"})
		if err != nil {
			cfg.Proxy.Status = 0
			cfg.Proxy.Port = portInt64

			dao.SetConfig(cfg)
			return dto.ReturnJsonDto{Code: 0, Msg: "启动失败: " + err.Error(), Type: "danger"}
		} else if res.Code == 1 {
			cfg.Proxy.Status = 1
			cfg.Proxy.Port = portInt64

			dao.SetConfig(cfg)
			go until.CleanAutoCacheAll() // 清理缓存
			return dto.ReturnJsonDto{Code: 1, Msg: "启动成功", Type: "success"}
		} else if res.Code != 1 {
			cfg.Proxy.Status = 0
			cfg.Proxy.Port = portInt64

			dao.SetConfig(cfg)
			return dto.ReturnJsonDto{Code: 0, Msg: "启动失败: " + res.Msg, Type: "danger"}
		} else {
			cfg.Proxy.Status = 0
			cfg.Proxy.Port = portInt64

			dao.SetConfig(cfg)
			dao.WS.SendWS(dao.Request{Action: "stopProxy"})
			go until.CleanAutoCacheAll()
			return dto.ReturnJsonDto{Code: 0, Msg: "启动失败: " + res.Msg, Type: "danger"}
		}
	} else {
		cfg.Proxy.Status = 0

		dao.SetConfig(cfg)
		dao.WS.SendWS(dao.Request{Action: "stopProxy"})
		go until.CleanAutoCacheAll()
		return dto.ReturnJsonDto{Code: 1, Msg: "停止成功", Type: "success"}
	}

}
