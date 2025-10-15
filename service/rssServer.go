package service

import (
	"encoding/json"
	"encoding/xml"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"log"
	"strconv"
	"strings"
	"time"
)

type RssUrl struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type AesData struct {
	I int64 `json:"i"`
}

func getAesdata(aesData AesData) (string, error) {
	jsonBytes, err := json.Marshal(aesData)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func getAesType(jsonStr string) (AesData, error) {
	var data AesData

	// 反序列化（字符串 -> 结构体）
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetRssUrl(id, host string, getnewkey bool) dto.ReturnJsonDto {
	var res []RssUrl

	var meal models.IptvMeals
	if err := dao.DB.Model(&models.IptvMeals{}).Where("id = ? and status = 1", id).First(&meal).Error; err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "未找到上线套餐", Type: "danger"}
	}

	aesData := AesData{
		I: meal.ID,
	}
	aesDataStr, err := getAesdata(aesData)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "生成key失败", Type: "danger"}
	}

	if getnewkey {
		cfg := dao.GetConfig()
		cfg.Rss.Key = until.Md5(time.Now().Format("2006-01-02 15:04:05"))
		until.RssKey = []byte(cfg.Rss.Key)
		dao.SetConfig(cfg)
	}

	aes := until.NewChaCha20(string(until.RssKey))
	token, err := aes.Encrypt(aesDataStr)
	if err != nil {
		return dto.ReturnJsonDto{Code: 0, Msg: "生成链接失败" + err.Error(), Type: "danger"}
	}

	res = append(res, RssUrl{Type: "m3u8", Url: host + "/getRss/" + token + "/paylist.m3u"})
	res = append(res, RssUrl{Type: "txt", Url: host + "/getRss/" + token + "/paylist.txt"})
	res = append(res, RssUrl{Type: "epg", Url: host + "/epg/" + token + "/e.xml"})

	return dto.ReturnJsonDto{Code: 1, Msg: "订阅生成成功", Type: "success", Data: res}
}

func GetRss(token, host, t string) string {

	aes := until.NewChaCha20(string(until.RssKey))
	jsonStr, err := aes.Decrypt(token)
	if err != nil {
		return "订阅失败,token解密错误"
	}
	aesData, err := getAesType(jsonStr)
	if err != nil {
		return "订阅失败，token读取错误"
	}
	if t == "t" {
		return getTxt(aesData.I)
	} else {
		return until.Txt2M3u8(getTxt(aesData.I), host, token)
	}
}

func getTxt(id int64) string {
	var res string
	var meal models.IptvMeals
	if err := dao.DB.Model(&models.IptvMeals{}).Where("id = ? and status = 1", id).First(&meal).Error; err != nil {
		return res
	}
	categoryNameList := strings.Split(meal.Content, "_")
	var categoryList []models.IptvCategory
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name in (?) and enable = 1", categoryNameList).Order("sort asc").Find(&categoryList).Error; err != nil {
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

func GetRssEpg(token, host string) dto.XmlTV {

	res := dto.XmlTV{
		GeneratorName: "清和IPTV管理系统",
		GeneratorURL:  "https://www.qingh.xyz",
	}
	aes := until.NewChaCha20(string(until.RssKey))
	jsonStr, err := aes.Decrypt(token)
	if err != nil {
		return res
	}
	aesData, err := getAesType(jsonStr)
	if err != nil {
		return res
	}
	return getEpg(aesData.I)
}

func getEpg(id int64) dto.XmlTV {

	res := dto.XmlTV{
		GeneratorName: "清和IPTV管理系统",
		GeneratorURL:  "https://www.qingh.xyz",
	}

	epgCaCheKey := "rssEpgXml_" + strconv.FormatInt(id, 10)
	if dao.Cache.Exists(epgCaCheKey) {
		cacheData, err := dao.Cache.Get(epgCaCheKey)
		if err == nil {
			err := xml.Unmarshal(cacheData, &res)
			if err == nil {
				return res
			}
		}
	}

	var meal models.IptvMeals
	if err := dao.DB.Model(&models.IptvMeals{}).Where("id = ? and status = 1", id).First(&meal).Error; err != nil {
		return res
	}
	categoryNameList := strings.Split(meal.Content, "_")
	var categoryList []models.IptvCategory
	if err := dao.DB.Model(&models.IptvCategory{}).Where("name in (?) and enable = 1", categoryNameList).Order("sort asc").Find(&categoryList).Error; err != nil {
		return res
	}
	for _, category := range categoryList {
		switch category.Sort {
		case -2:
			tvData := GetCntvEpgXml()
			res.Channels = append(res.Channels, tvData.Channels...)
			res.Programmes = append(res.Programmes, tvData.Programmes...)
		case -1:
			tvData := GetProvinceEpgXml()
			res.Channels = append(res.Channels, tvData.Channels...)
			res.Programmes = append(res.Programmes, tvData.Programmes...)
		default:
			tvData := GetEpgXml(category.Name)
			res.Channels = append(res.Channels, tvData.Channels...)
			res.Programmes = append(res.Programmes, tvData.Programmes...)
		}
	}
	res = CleanTV(res)

	data, err := xml.Marshal(res)
	if err == nil {
		err := dao.Cache.Set(epgCaCheKey, data)
		if err != nil {
			log.Println("epg缓存设置失败:", err)
			dao.Cache.Delete(epgCaCheKey)
		}
	} else {
		log.Println("epg缓存序列化失败:", err)
		dao.Cache.Delete(epgCaCheKey)
	}
	return res
}

func CleanTV(tv dto.XmlTV) dto.XmlTV {
	// 1️⃣ 去重 Channel（按 ID 保留第一个）
	uniqueChannels := make([]dto.XmlChannel, 0, len(tv.Channels))
	seen := make(map[string]bool)
	ids := make(map[string]int)
	i := 1
	for _, ch := range tv.Channels {
		if !seen[ch.DisplayName.Value] {
			seen[ch.DisplayName.Value] = true
			ids[ch.DisplayName.Value] = i
			ch.ID = strconv.Itoa(i)
			uniqueChannels = append(uniqueChannels, ch)
			i++
		}
	}
	tv.Channels = uniqueChannels

	// 2️⃣ 删除无效的 Programme（仅保留 channel 存在的）
	validProgrammes := make([]dto.Programme, 0, len(tv.Programmes))
	progSet := make(map[string]bool) // 记录唯一键

	for _, p := range tv.Programmes {
		if seen[p.Channel] {
			p.Channel = strconv.Itoa(ids[p.Channel])
			key := p.Channel + "_" + p.Start + "_" + p.Title.Value // 唯一键

			if !progSet[key] {
				validProgrammes = append(validProgrammes, p)
				progSet[key] = true
			}
		}
	}
	tv.Programmes = validProgrammes

	return tv
}

func GetCntvEpgXml() dto.XmlTV {
	cntvXml := dto.XmlTV{
		GeneratorName: "清和IPTV管理系统",
		GeneratorURL:  "https://www.qingh.xyz",
	}

	var epgs []models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Where("name like ? and status = 1", "cntv-%").Find(&epgs).Error; err != nil {
		return cntvXml
	}

	for _, epg := range epgs {
		if epg.Content == "" {
			continue
		}
		// eFrom := strings.Split(epg.Name, "-")[0]
		eName := strings.Split(epg.Name, "-")[1]
		nameList := strings.Split(epg.Content, ",")
		var channelList []models.IptvChannel
		if err := dao.DB.Model(&models.IptvChannel{}).Where("name in (?)", nameList).Order("sort asc").Find(&channelList).Error; err != nil {
			continue
		}
		for _, channel := range channelList {
			tmpData, err := until.GetEpgCntv(eName)
			if err == nil {
				tmpXml := until.ConvertCntvToXml(tmpData, eName)
				cntvXml.Channels = append(cntvXml.Channels, dto.XmlChannel{
					ID: eName,
					DisplayName: dto.DisplayName{
						Lang:  "zh",
						Value: channel.Name,
					},
				})

				for _, p := range tmpXml.Programmes {
					p.Channel = channel.Name
					cntvXml.Programmes = append(cntvXml.Programmes, p)
				}
			}
		}
	}
	return cntvXml
}

func GetProvinceEpgXml() dto.XmlTV {
	epgXml := dto.XmlTV{
		GeneratorName: "清和IPTV管理系统",
		GeneratorURL:  "https://www.qingh.xyz",
	}

	var epgs []models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Where("name like ? and status = 1", "%"+"-%卫视").Find(&epgs).Error; err != nil {
		return epgXml
	}

	for _, epg := range epgs {
		if epg.Content == "" {
			continue
		}
		nameList := strings.Split(epg.Content, ",")
		var channelList []models.IptvChannel
		if err := dao.DB.Model(&models.IptvChannel{}).Where("name in (?)", nameList).Order("sort asc").Find(&channelList).Error; err != nil {
			continue
		}
		if len(channelList) == 0 {
			continue
		}
		eFrom := strings.Split(epg.Name, "-")[0]
		eName := strings.Split(epg.Name, "-")[1]

		var epgList models.IptvEpgList
		if err := dao.DB.Model(&models.IptvEpgList{}).Where("remarks = ? and status = 1", eFrom).First(&epgList).Error; err != nil {
			continue
		}
		for _, channel := range channelList {
			tmpXml := until.GetEpgListXml(epgList.Name, epgList.Url)
			epgXml.Channels = append(epgXml.Channels, dto.XmlChannel{
				ID: eName,
				DisplayName: dto.DisplayName{
					Lang:  "zh",
					Value: channel.Name,
				},
			})
			var cId string
			for _, c := range tmpXml.Channels {
				if c.DisplayName.Value == channel.Name {
					cId = c.ID
					break
				}
			}

			for _, p := range tmpXml.Programmes {
				if p.Channel == cId {
					p.Channel = channel.Name
					epgXml.Programmes = append(epgXml.Programmes, p)
				}
			}
		}

	}
	return epgXml
}

func GetEpgXml(cname string) dto.XmlTV {
	epgXml := dto.XmlTV{
		GeneratorName: "清和IPTV管理系统",
		GeneratorURL:  "https://www.qingh.xyz",
	}

	var channelList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).
		Where("category = ?", cname).
		Order("sort asc").
		Find(&channelList).Error; err != nil {
		return epgXml
	}
	if len(channelList) == 0 {
		return epgXml
	}

	for _, channel := range channelList {
		var epgs []models.IptvEpg
		if err := dao.DB.Model(&models.IptvEpg{}).Where("content like ? and status = 1", "%"+channel.Name+"%").Find(&epgs).Error; err != nil {
			continue
		}
		for _, epg := range epgs {
			eType := strings.SplitN(epg.Name, "-", 2)[0]
			eName := strings.SplitN(epg.Name, "-", 2)[1]
			if eType == "cntv" {
				tmpData, err := until.GetEpgCntv(eName)
				if err == nil {
					tmpXml := until.ConvertCntvToXml(tmpData, eName)
					epgXml.Channels = append(epgXml.Channels, dto.XmlChannel{
						ID: epg.Name,
						DisplayName: dto.DisplayName{
							Lang:  "zh",
							Value: channel.Name,
						},
					})
					for _, p := range tmpXml.Programmes {
						p.Channel = channel.Name
						epgXml.Programmes = append(epgXml.Programmes, p)
					}
				}
				if len(epgXml.Channels) > 0 && len(epgXml.Programmes) > 0 {
					break
				}
				continue
			}

			var epgList models.IptvEpgList
			if err := dao.DB.Model(&models.IptvEpgList{}).Where("remarks = ? and status = 1", eType).First(&epgList).Error; err != nil {
				continue
			}
			tmpXml := until.GetEpgListXml(epgList.Name, epgList.Url)
			epgXml.Channels = append(epgXml.Channels, dto.XmlChannel{
				ID: epg.Name,
				DisplayName: dto.DisplayName{
					Lang:  "zh",
					Value: channel.Name,
				},
			})
			var cId string
			for _, c := range tmpXml.Channels {
				if c.DisplayName.Value == channel.Name {
					cId = c.ID
					break
				}
			}

			for _, p := range tmpXml.Programmes {
				if p.Channel == cId {
					p.Channel = epg.Name
					epgXml.Programmes = append(epgXml.Programmes, p)
				}
			}
		}
	}
	return epgXml
}
