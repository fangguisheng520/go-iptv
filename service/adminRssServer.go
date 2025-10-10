package service

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"strconv"
	"strings"
)

type RssUrl struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

func GetRssUrl(id string) dto.ReturnJsonDto {
	var res []RssUrl

	var meal models.IptvMeals
	if err := dao.DB.Model(&models.IptvMeals{}).Where("id = ?", id).First(&meal).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到套餐", Type: "danger"}
	}
	tokenM3u8, err1 := until.GenerateJWTRss("m3u8", strconv.FormatInt(meal.ID, 10))
	if err1 != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "生成链接失败", Type: "danger"}
	}
	tokenTxt, err2 := until.GenerateJWTRss("txt", strconv.FormatInt(meal.ID, 10))
	if err2 != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "生成链接失败", Type: "danger"}
	}
	cfg := dao.GetConfig()
	res = append(res, RssUrl{Type: "m3u8", Url: cfg.ServerUrl + "/getRss?token=" + tokenM3u8})
	res = append(res, RssUrl{Type: "txt", Url: cfg.ServerUrl + "/getRss?token=" + tokenTxt})

	return dto.ReturnJsonDto{Code: 1, Msg: "订阅生成成功", Type: "success", Data: res}
}

func GetRss(token string) string {

	rssType, id, err := until.VerifyJWTRss(token)
	if err != nil {
		return "订阅失败"
	}
	if rssType == "txt" {
		return getTxt(id)
	} else {
		return until.Txt2M3u8(getTxt(id))
	}
}

func getTxt(id int64) string {
	var res string
	var meal models.IptvMeals
	if err := dao.DB.Model(&models.IptvMeals{}).Where("id = ?", id).First(&meal).Error; err != nil {
		return res
	}
	categoryNameList := strings.Split(meal.Content, "_")
	var categoryList []models.IptvCategory
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name in (?)", categoryNameList).Order("sort asc").Find(&categoryList).Error; err != nil {
		return res
	}

	for _, category := range categoryList {
		switch category.Sort {
		case -2:
			data := until.GetCCTVChannelList(false)
			if data != "" {
				res += category.Name + ",#genre#\n"
				res += data

			}
		case -1:
			data := until.GetProvinceChannelList(false)
			if data != "" {
				res += category.Name + ",#genre#\n"
				res += data

			}
		default:
			var channels []models.IptvChannel
			if err := dao.DB.Model(&models.IptvChannel{}).
				Where("category = ?", category.Name).
				Order("sort asc").Find(&channels).Error; err != nil {
				continue
			}
			if len(channels) == 0 {
				continue
			}
			res += category.Name + ",#genre#\n"
			for _, channel := range channels {
				res += channel.Name + "," + channel.Url + "\n"
			}
		}
	}
	return res
}
