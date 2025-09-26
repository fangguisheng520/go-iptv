package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"net/url"
)

func ChangeMovieStatus(params url.Values) dto.ReturnJsonDto {
	id := params.Get("change_status")
	if id == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "点播 id不能为空", Type: "danger"}
	}

	if !until.IsSafe(id) {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数不合法，请勿输入特殊字符", Type: "danger"}
	}

	var epgData models.IptvMovie
	if err := dao.DB.Model(&models.IptvMovie{}).Where("id = ?", id).First(&epgData).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "查询点播失败", Type: "danger"}
	}

	if epgData.State == 1 {
		dao.DB.Model(&models.IptvMovie{}).Where("id = ?", id).Update("state", 0)
	} else {
		dao.DB.Model(&models.IptvMovie{}).Where("id = ?", id).Update("state", 1)
	}
	return dto.ReturnJsonDto{Code: 1, Msg: "点播 " + epgData.Name + "状态修改成功", Type: "success"}
}

func DeleteMovie(params url.Values) dto.ReturnJsonDto {
	movieId := params.Get("delmovie")
	if movieId == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "没有获取到点播ID", Type: "danger"}
	}
	if !until.IsSafe(movieId) {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数不合法，请勿输入特殊字符", Type: "danger"}
	}
	if err := dao.DB.Where("id = ?", movieId).Delete(&models.IptvMovie{}).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "删除失败", Type: "danger"}
	}
	return dto.ReturnJsonDto{Code: 1, Msg: "删除成功", Type: "success"}
}

func EditMovie(params url.Values) dto.ReturnJsonDto {
	movieId := params.Get("movieId")
	movieName := params.Get("movieName")
	movieApi := params.Get("movieApi")

	if movieName == "" || movieApi == "" {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数错误", Type: "danger"}
	}
	if !until.IsSafe(movieId) || !until.IsSafe(movieName) || !until.IsSafe(movieApi) {
		return dto.ReturnJsonDto{Code: 0, Msg: "参数不合法，请勿输入特殊字符", Type: "danger"}
	}

	if movieId == "" {
		dao.DB.Model(&models.IptvMovie{}).Create(&models.IptvMovie{Name: movieName, Api: movieApi})
		return dto.ReturnJsonDto{Code: 1, Msg: "点播 " + movieName + "添加成功", Type: "success"}
	} else {
		dao.DB.Model(&models.IptvMovie{}).Where("id = ?", movieId).Update("name", movieName).Update("api", movieApi)
		return dto.ReturnJsonDto{Code: 1, Msg: "点播 " + movieName + "修改成功", Type: "success"}
	}
}
