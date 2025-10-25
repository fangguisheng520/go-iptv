package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"net/url"
	"strconv"
)

func Notice(params url.Values) dto.ReturnJsonDto {
	adtext := params.Get("adtext")
	showtime := params.Get("showtime")
	showinterval := params.Get("showinterval")

	showtimeNum, err := strconv.ParseInt(showtime, 10, 64)
	if err != nil {
		showtimeNum = 0
	}

	showintervalNum, err := strconv.ParseInt(showinterval, 10, 64)
	if err != nil {
		showintervalNum = 0
	}

	cfg := dao.GetConfig()

	cfg.Ad.AdText = adtext
	cfg.Ad.ShowInterval = showintervalNum
	cfg.Ad.ShowTime = showtimeNum

	dao.SetConfig(cfg)

	// TODO: 调用接口
	return dto.ReturnJsonDto{Code: 1, Msg: "修改成功", Type: "success"}
}
