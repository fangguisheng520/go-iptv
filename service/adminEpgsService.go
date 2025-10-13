package service

import (
	"encoding/xml"
	"errors"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetEpgData(params url.Values) dto.ReturnJsonDto {
	//编辑
	epgId := params.Get("editepg")
	if epgId == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "EPG id不能为空", Type: "danger"}
	}

	var epg models.IptvEpg
	if err := dao.DB.Where("id = ?", epgId).First(&epg).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "查询EPG 失败", Type: "danger"}
	}

	var channeList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Select("distinct name").Order("category,id").Find(&channeList).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "查询频道失败", Type: "danger"}
	}

	epgList := strings.Split(epg.Remarks, "|")

	var dataList []dto.EpgsReturnDto

	for _, v := range channeList {
		var data dto.EpgsReturnDto
		data.Name = v.Name
		data.Checked = false
		for _, v1 := range epgList {
			if strings.EqualFold(v1, v.Name) {
				data.Checked = true
			}
		}
		if strings.EqualFold(epg.Name, v.Name) {
			data.Checked = true
		}
		dataList = append(dataList, data)
	}

	return dto.ReturnJsonDto{Code: 1, Msg: "操作成功", Type: "success", Data: dataList}
}

func SaveEpg(params url.Values, editType int) dto.ReturnJsonDto {
	if editType == 1 {
		id := params.Get("id")
		if id == "" {
			return dto.ReturnJsonDto{Code: 0, Msg: "EPG id不能为空", Type: "danger"}
		}
		name := params.Get("name")
		if name == "" {
			return dto.ReturnJsonDto{Code: 0, Msg: "EPG 名称不能为空", Type: "danger"}
		}
		epg := params.Get("epg")
		if epg == "" {
			return dto.ReturnJsonDto{Code: 0, Msg: "EPG 来源不能为空", Type: "danger"}
		}

		remarks := params.Get("remarks")
		namesList := params["names[]"]

		var epgData models.IptvEpg

		if err := dao.DB.Model(&models.IptvEpg{}).Where("id = ?", id).First(&epgData).Error; err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "查询EPG失败", Type: "danger"}
		}
		epgData.Name = epg + "-" + name
		epgData.Content = strings.Join(namesList, ",")
		epgData.Remarks = remarks

		if err := dao.DB.Save(&epgData).Error; err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "保存EPG失败", Type: "danger"}
		}
		return dto.ReturnJsonDto{Code: 1, Msg: "EPG " + epgData.Name + "保存成功", Type: "success"}
	} else {
		name := params.Get("name")
		if name == "" {
			return dto.ReturnJsonDto{Code: 0, Msg: "EPG 名称不能为空", Type: "danger"}
		}
		epg := params.Get("epg")
		if epg == "" {
			return dto.ReturnJsonDto{Code: 0, Msg: "EPG 来源不能为空", Type: "danger"}
		}

		remarks := params.Get("remarks")
		namesList := params["names[]"]

		var epgData models.IptvEpg
		epgData.Name = epg + "-" + name
		epgData.Content = strings.Join(namesList, ",")
		epgData.Remarks = remarks
		epgData.Status = 1

		if err := dao.DB.Save(&epgData).Error; err != nil {
			return dto.ReturnJsonDto{Code: 0, Msg: "保存EPG失败", Type: "danger"}
		}
		return dto.ReturnJsonDto{Code: 1, Msg: "EPG " + epgData.Name + "保存成功", Type: "success"}
	}
}

func ChangeStatus(params url.Values) dto.ReturnJsonDto {
	id := params.Get("change_status")
	if id == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "EPG id不能为空", Type: "danger"}
	}

	var epgData models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Where("id = ?", id).First(&epgData).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "查询EPG失败", Type: "danger"}
	}

	if epgData.Status == 1 {
		dao.DB.Model(&models.IptvEpg{}).Where("id = ?", id).Update("status", 0)
	} else {
		dao.DB.Model(&models.IptvEpg{}).Where("id = ?", id).Update("status", 1)
	}
	return dto.ReturnJsonDto{Code: 1, Msg: "EPG " + epgData.Name + "状态修改成功", Type: "success"}
}

func DeleteEpg(params url.Values) dto.ReturnJsonDto {
	id := params.Get("delepg")
	if id == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "EPG id不能为空", Type: "danger"}
	}
	dao.DB.Where("id = ?", id).Delete(&models.IptvEpg{})
	return dto.ReturnJsonDto{Code: 1, Msg: "EPG删除成功", Type: "success"}
}

func BindChannel() dto.ReturnJsonDto {
	ClearBind() // 清空绑定
	var channelList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Select("distinct name").Order("category,id").Find(&channelList).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "查询频道失败", Type: "danger"}
	}

	var epgList []models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Find(&epgList).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "查询EPG失败", Type: "danger"}
	}

	for _, epgData := range epgList {
		var tmpList []string
		for _, channelData := range channelList {

			if strings.EqualFold(channelData.Name, epgData.Name) {
				tmpList = append(tmpList, channelData.Name)
				break
			}

			nameList := strings.Split(epgData.Remarks, "|")
			for _, name := range nameList {
				if strings.EqualFold(channelData.Name, name) {
					tmpList = append(tmpList, channelData.Name)
					break
				}
			}
		}
		epgData.Content = strings.Join(tmpList, ",")
		if epgData.Content != "" {
			dao.DB.Save(&epgData)
		}
	}
	go until.GetCCTVChannelList(true)
	go until.GetProvinceChannelList(true)
	go until.CleanMealsXmlCacheAll()
	return dto.ReturnJsonDto{Code: 1, Msg: "绑定成功", Type: "success"}
}

func ClearBind() dto.ReturnJsonDto {
	dao.DB.Model(&models.IptvEpg{}).Where("content != ''").Update("content", "")
	return dto.ReturnJsonDto{Code: 1, Msg: "清除绑定成功", Type: "success"}
}

func ClearCache() dto.ReturnJsonDto {
	dao.Cache.Clear()
	return dto.ReturnJsonDto{Code: 1, Msg: "清除缓存成功", Type: "success"}
}

func EpgImportFile(c *gin.Context) dto.ReturnJsonDto {

	epgFromName := c.PostForm("epgfromname")
	if epgFromName == "" || strings.Contains(epgFromName, "-") || !until.IsSafe(epgFromName) {
		return dto.ReturnJsonDto{Code: 0, Msg: "EPG来源名称不合法", Type: "danger"}
	}

	file, err := c.FormFile("epgfile")
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取文件失败:" + err.Error(), Type: "danger"}
	}

	f, err := file.Open()
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "打开文件失败: " + err.Error(), Type: "danger"}
	}
	defer f.Close()

	// 读取内容
	data, err := io.ReadAll(f)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "读取文件失败: " + err.Error(), Type: "danger"}
	}

	var xmlData dto.TV
	errXml := xml.Unmarshal(data, &xmlData)
	if errXml != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "xml格式化失败: " + errXml.Error(), Type: "danger"}
	}

	var eOld models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Where("name like ?", epgFromName+"-%").First(&eOld).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.ReturnJsonDto{Code: 0, Msg: "EPG查询失败: " + err.Error(), Type: "danger"}
	}
	if eOld.ID != 0 {
		dao.DB.Model(&models.IptvEpg{}).Where("name like ?", epgFromName+"-%").Delete(&models.IptvEpg{})
	}

	var epgs []models.IptvEpg
	for _, channel := range xmlData.Channels {
		epgs = append(epgs, models.IptvEpg{
			Name:    epgFromName + "-" + channel.DisplayName.Value,
			Status:  1,
			Remarks: channel.DisplayName.Value,
		})
	}
	if err := dao.DB.Model(&models.IptvEpg{}).Create(&epgs).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "保存失败: " + err.Error(), Type: "danger"}
	}
	go BindChannel()
	return dto.ReturnJsonDto{Code: 1, Msg: "保存成功", Type: "success"}
}

func UploadLogo(c *gin.Context) dto.ReturnJsonDto {

	epgFromName := c.PostForm("epgname")
	log.Println(epgFromName)
	if epgFromName == "" || !strings.Contains(epgFromName, "-") || !until.IsSafe(epgFromName) {
		return dto.ReturnJsonDto{Code: 0, Msg: "EPG名称不合法", Type: "danger"}
	}

	epgName := strings.Split(epgFromName, "-")[1]

	file, err := c.FormFile("uploadlogo")
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "获取文件失败:" + err.Error(), Type: "danger"}
	}

	f, err := file.Open()
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "打开文件失败:" + err.Error(), Type: "danger"}
	}
	defer f.Close()

	// 读取前 512 字节判断 MIME 类型
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	contentType := http.DetectContentType(buf[:n])

	if contentType != "image/png" {
		return dto.ReturnJsonDto{Code: 0, Msg: "只允许上传 PNG 文件", Type: "danger"}
	}

	dst := "/config/logo/" + epgName + ".png"
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "保存文件失败:" + err.Error(), Type: "danger"}
	}
	go until.CleanMealsXmlCacheAll()
	return dto.ReturnJsonDto{Code: 1, Msg: "上传成功", Type: "success"}
}

// func SaveEpgApi(params url.Values) dto.ReturnJsonDto {
// 	err1000 := params.Get("tipepgerror_1000")
// 	err1001 := params.Get("tipepgerror_1001")
// 	err1002 := params.Get("tipepgerror_1002")
// 	err1003 := params.Get("tipepgerror_1003")
// 	err1004 := params.Get("tipepgerror_1004")
// 	err1005 := params.Get("tipepgerror_1005")
// 	epgapiChk := params.Get("epgapi_chk")

// 	if err1000 == "" && err1001 == "" && err1002 == "" && err1003 == "" && err1004 == "" && err1005 == "" {
// 		return dto.ReturnJsonDto{Code: 0, Msg: "错误提示存在空", Type: "error"}
// 	}

// 	cfg := dao.GetConfig()

// 	cfg.EPGErrors.Err1000 = err1000
// 	cfg.EPGErrors.Err1001 = err1001
// 	cfg.EPGErrors.Err1002 = err1002
// 	cfg.EPGErrors.Err1003 = err1003
// 	cfg.EPGErrors.Err1004 = err1004
// 	cfg.EPGErrors.Err1005 = err1005

// 	if epgapiChk == "on" {
// 		cfg.App.EPGApiChk = 1
// 	} else {
// 		cfg.App.EPGApiChk = 0
// 	}

// 	dao.SetConfig(cfg)

// 	return dto.ReturnJsonDto{Code: 1, Msg: "保存EPG成功", Type: "success"}
// }
