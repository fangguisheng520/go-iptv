package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"net/url"
	"strings"
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
	var channeList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Select("distinct name").Order("category,id").Find(&channeList).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "查询频道失败", Type: "danger"}
	}

	var epgList []models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Find(&epgList).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "查询EPG失败", Type: "danger"}
	}

	for _, epgData := range epgList {
		var tmpList []string
		for _, channelData := range channeList {

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
