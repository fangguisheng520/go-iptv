package service

import (
	"fmt"
	"go-iptv/bootstrap"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/until"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context, imgType string) dto.ReturnJsonDto {
	if imgType == "icon" {
		file, err := c.FormFile("iconfile")
		if err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "获取文件失败:" + err.Error(), Type: "danger"}
		}
		dst := "./images/icon/icon.png"
		if err := c.SaveUploadedFile(file, dst); err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "保存文件失败:" + err.Error(), Type: "danger"}
		}
		return dto.ReturnJsonDto{Code: 1, Msg: "上传成功", Type: "success", Data: map[string]interface{}{"url": "/icon/icon.png"}}
	} else {
		file, err := c.FormFile("bjfile")
		if err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "获取文件失败:" + err.Error(), Type: "danger"}
		}
		pngName := until.Md5(url.QueryEscape(fmt.Sprintf("%s%d", file.Filename, time.Now().Unix())))

		dst := "/app/images/bj/" + pngName + ".png"
		if err := c.SaveUploadedFile(file, dst); err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "保存文件失败:" + err.Error(), Type: "danger"}
		}
		return dto.ReturnJsonDto{Code: 1, Msg: "上传成功", Type: "success", Data: map[string]interface{}{"name": pngName}}
	}
}

func DeleteFile(params url.Values, imgType string) dto.ReturnJsonDto {
	// iconFile := params.Get("iconfile")
	if imgType == "icon" {
		iconFile := "./images/icon/icon.png"
		if err := os.Remove(iconFile); err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "删除失败:" + err.Error(), Type: "danger"}
		}
		return dto.ReturnJsonDto{Code: 1, Msg: "删除成功", Type: "success"}
	} else {
		bjName := params.Get("deleteBj")
		if !until.IsSafeImgName(bjName) {
			return dto.ReturnJsonDto{Code: 0, Msg: "文件名不合法", Type: "danger"}
		}
		bjFile := "/app/images/bj/" + bjName + ".png"
		if err := os.Remove(bjFile); err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "删除失败", Type: "danger"}
		}
		return dto.ReturnJsonDto{Code: 1, Msg: "删除成功", Type: "success"}
	}
}

func DecoderSelect(params url.Values) dto.ReturnJsonDto {
	decoder := params.Get("decoder")
	if decoder != "0" && decoder != "1" && decoder != "2" {
		return dto.ReturnJsonDto{Code: 0, Msg: "解码器选择失败", Type: "danger"}
	} else {
		cfg := dao.GetConfig()
		decoderInt, err := strconv.ParseInt(decoder, 10, 64)
		if err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "解码器选择失败", Type: "danger"}
		}
		cfg.App.Decoder = decoderInt
		dao.SetConfig(cfg)
		return dto.ReturnJsonDto{Code: 1, Msg: "解码器选择成功", Type: "success"}
	}
}

func SetBuffTimeOut(params url.Values) dto.ReturnJsonDto {
	buffTimeOut := params.Get("buffTimeOut")
	if buffTimeOut != "5" && buffTimeOut != "10" && buffTimeOut != "15" && buffTimeOut != "20" && buffTimeOut != "25" && buffTimeOut != "30" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	} else {
		cfg := dao.GetConfig()
		buffTimeOutInt, err := strconv.ParseInt(buffTimeOut, 10, 64)
		if err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
		}
		cfg.App.BuffTimeout = buffTimeOutInt
		dao.SetConfig(cfg)
		return dto.ReturnJsonDto{Code: 1, Msg: "超时设置成功", Type: "success"}
	}
}

func SetNeedAuthor(params url.Values) dto.ReturnJsonDto {
	needauthor := params.Get("needauthor")
	if needauthor != "1" && needauthor != "0" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	} else {
		cfg := dao.GetConfig()
		needauthorInt, err := strconv.ParseInt(needauthor, 10, 64)
		if err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
		}
		cfg.App.NeedAuthor = needauthorInt
		dao.SetConfig(cfg)
		return dto.ReturnJsonDto{Code: 1, Msg: "授权设置成功", Type: "success"}
	}
}

func SetAppInfo(params url.Values) dto.ReturnJsonDto {

	buildStatus := bootstrap.GetBuildStatus()
	if buildStatus == 1 {
		return dto.ReturnJsonDto{Code: 0, Msg: "正在打包中，请稍后再试", Type: "danger"}
	}
	appName := params.Get("app_appname")
	appPackag := params.Get("app_packagename")
	appVersion := params.Get("app_version")
	appSign := params.Get("app_sign")
	upSet := params.Get("up_sets")
	upText := params.Get("up_text")

	if appName == "" || appPackag == "" || appVersion == "" || appSign == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	} else {
		cfg := dao.GetConfig()
		cfg.Build.Name = appName
		cfg.Build.Package = appPackag

		if cfg.Build.Version == appVersion {
			return dto.ReturnJsonDto{Code: 0, Msg: "版本号不能相同", Type: "danger"}
		}
		cfg.Build.Version = appVersion
		appSignInt, err := strconv.ParseInt(appSign, 10, 64)
		if err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "签名参数错误", Type: "danger"}
		}
		if appSignInt < 1 || appSignInt > 65535 {
			return dto.ReturnJsonDto{Code: 0, Msg: "签名参数超过范围", Type: "danger"}
		}

		if upSet == "on" || upSet == "1" || upSet == "true" {
			cfg.App.Update.Set = 1
		} else {
			cfg.App.Update.Set = 0
		}
		cfg.App.Update.Text = upText
		cfg.Build.Sign = appSignInt
		// cfg.App.Update.Url = strings.TrimSuffix(cfg.ServerUrl, "/") + "/app/" + cfg.Build.Name + ".apk"
		dao.SetConfig(cfg)
		// go bootstrap.BuildAPK() // 启动编译
		return dto.ReturnJsonDto{Code: 1, Msg: "APK编译中...", Type: "success"}
	}
}

func SetTipSet(params url.Values) dto.ReturnJsonDto {
	tiploading := params.Get("tiploading")
	tipuserexpired := params.Get("tipuserexpired")
	tipuserforbidden := params.Get("tipuserforbidden")
	tipusernoreg := params.Get("tipusernoreg")

	cfg := dao.GetConfig()

	cfg.Tips.Loading = tiploading
	cfg.Tips.UserExpired = tipuserexpired
	cfg.Tips.UserForbidden = tipuserforbidden
	cfg.Tips.UserNoReg = tipusernoreg

	dao.SetConfig(cfg)
	return dto.ReturnJsonDto{Code: 1, Msg: "设置成功", Type: "success"}
}
