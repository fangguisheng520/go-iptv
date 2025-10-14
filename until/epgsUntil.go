package until

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ConvertCntvToXml(cntv dto.CntvJsonChannel, cName string) dto.XmlTV {
	tv := dto.XmlTV{
		GeneratorName: "清和IPTV管理系统",
		GeneratorURL:  "https://www.qingh.xyz",
	}

	log.Println("开始转换", cName)

	// 添加频道
	tv.Channels = append(tv.Channels, dto.XmlChannel{
		ID: cName,
		DisplayName: dto.DisplayName{
			Lang:  "zh",
			Value: cName,
		},
	})

	// 添加节目表
	for _, p := range cntv.Program {
		start := time.Unix(p.StartTime, 0).UTC().Format("20060102150405 -0700")
		stop := time.Unix(p.EndTime, 0).UTC().Format("20060102150405 -0700")

		tv.Programmes = append(tv.Programmes, dto.Programme{
			Start:   start,
			Stop:    stop,
			Channel: cName,
			Title: dto.Title{
				Lang:  "zh",
				Value: p.Title,
			},
			Desc: dto.Desc{
				Lang:  "zh",
				Value: p.Title,
			},
		})
	}

	return tv
}

func GetEpgListXml(name, url string) dto.XmlTV {
	epgUrl := url
	cacheKey := "epgXmlFrom_" + name
	var xmlTV dto.XmlTV
	var xmlByte []byte
	readCacheOk := false
	if dao.Cache.Exists(cacheKey) {
		tmpByte, err := dao.Cache.Get(cacheKey)
		if err == nil {
			xmlByte = tmpByte
			readCacheOk = true
		}
	}

	if !readCacheOk {
		xmlByte = []byte(GetUrlData(epgUrl))
		if dao.Cache.Set(cacheKey, xmlByte) != nil {
			dao.Cache.Delete(cacheKey)
		}
	}
	xml.Unmarshal(xmlByte, &xmlTV)
	return xmlTV
}

func GetEpgCntv(name string) (dto.CntvJsonChannel, error) {

	var cacheKey = "cntv_" + strings.ToUpper(name)

	var cntvJson dto.CntvData

	if name == "" {
		return dto.CntvJsonChannel{}, errors.New("id is empty")
	}
	name = strings.ToLower(name)

	epgUrl := "https://api.cntv.cn/epg/epginfo?c=" + name + "&serviceId=channel&d="

	readCacheOk := false
	if dao.Cache.Exists(cacheKey) {
		err := dao.Cache.GetJSON(cacheKey, cntvJson)
		if err == nil {
			readCacheOk = true
		}
	}

	if !readCacheOk {
		jsonStr := GetUrlData(epgUrl)
		err := json.Unmarshal([]byte(jsonStr), &cntvJson)
		if err != nil {
			return dto.CntvJsonChannel{}, err
		}
		if dao.Cache.SetJSON(cacheKey, cntvJson) != nil {
			dao.Cache.Delete(cacheKey)
		}
	}
	return cntvJson[name], nil
}

func CleanMealsXmlCacheAll() {
	var meals []models.IptvMeals
	dao.DB.Model(&models.IptvMeals{}).Find(&meals)
	for _, meal := range meals {
		dao.Cache.Delete("rssEpgXml_" + strconv.FormatInt(meal.ID, 10))
	}
}

func CleanMealsXmlCacheOne(id int64) {
	dao.Cache.Delete("rssEpgXml_" + strconv.FormatInt(id, 10))
}

func UpdataEpgList() bool {
	var epgLists []models.IptvEpgList
	dao.DB.Model(&models.IptvEpgList{}).Find(&epgLists)
	for _, list := range epgLists {
		cacheKey := "epgXmlFrom_" + list.Name
		dao.Cache.Delete(cacheKey)
		xmlStr := GetUrlData(list.Url)
		if xmlStr != "" {
			xmlByte := []byte(xmlStr)
			if dao.Cache.Set(cacheKey, xmlByte) != nil {
				dao.Cache.Delete(cacheKey)
			}
			var xmlTV dto.XmlTV
			if xml.Unmarshal(xmlByte, &xmlTV) != nil {
				continue
			}
			var epgs []models.IptvEpg
			// 1️⃣ 匹配数字台，如 CCTV1、CCTV-5+、CCTV13 等
			reNum := regexp.MustCompile(`(?i)CCTV-?(\d+\+?)`)

			// 2️⃣ 匹配字母台，如 CCTV4EUO、CCTV4AME、CCTVF、CCTVE 等
			reAlpha := regexp.MustCompile(`(?i)CCTV(\d*[A-Z]+)`)
			for _, channel := range xmlTV.Channels {
				remarks := channel.DisplayName.Value
				upper := strings.ToUpper(remarks)
				if strings.Contains(upper, "CCTV") {
					switch {
					case reNum.MatchString(upper):
						match := reNum.FindStringSubmatch(upper)
						num := match[1]
						remarks = fmt.Sprintf("CCTV%s|CCTV-%s", num, num)

					case reAlpha.MatchString(upper):
						match := reAlpha.FindStringSubmatch(upper)
						suffix := match[1]
						remarks = fmt.Sprintf("CCTV%s|CCTV-%s", suffix, suffix)
					}
				}
				epgs = append(epgs, models.IptvEpg{
					Name:    list.Remarks + "-" + channel.DisplayName.Value,
					Status:  1,
					Remarks: remarks,
				})
			}
			if len(epgs) > 0 {
				dao.DB.Model(&models.IptvEpgList{}).Where("id = ?", list.ID).Updates(&models.IptvEpgList{Status: 1, LastTime: time.Now().Unix()})
				dao.DB.Model(&models.IptvEpg{}).Where("name like ?", list.Remarks+"-%").Delete(&models.IptvEpg{})
				dao.DB.Model(&models.IptvEpg{}).Create(&epgs)
				CleanMealsXmlCacheAll() // 清除缓存
				go BindChannel()        // 绑定频道
			}
		}
	}
	return true
}

func UpdataEpgListOne(id int64) bool {
	var list models.IptvEpgList
	if err := dao.DB.Model(&models.IptvEpgList{}).Where("id = ?", id).First(&list).Error; err != nil {
		return false
	}
	cacheKey := "epgXmlFrom_" + list.Name
	dao.Cache.Delete(cacheKey)
	xmlStr := GetUrlData(list.Url)
	if xmlStr != "" {
		xmlByte := []byte(xmlStr)
		if dao.Cache.Set(cacheKey, xmlByte) != nil {
			dao.Cache.Delete(cacheKey)
		}
		var xmlTV dto.XmlTV
		if xml.Unmarshal(xmlByte, &xmlTV) != nil {
			return false
		}
		var epgs []models.IptvEpg
		// 1️⃣ 匹配数字台，如 CCTV1、CCTV-5+、CCTV13 等
		reNum := regexp.MustCompile(`(?i)CCTV-?(\d+\+?)$`)

		// 2️⃣ 匹配字母台，如 CCTV4EUO、CCTV4AME、CCTVF、CCTVE 等
		reAlpha := regexp.MustCompile(`(?i)CCTV(\d*[A-Z]+)`)
		for _, channel := range xmlTV.Channels {
			remarks := channel.DisplayName.Value
			upper := strings.ToUpper(remarks)
			if strings.Contains(upper, "CCTV") {
				switch {
				case reNum.MatchString(upper):
					match := reNum.FindStringSubmatch(upper)
					num := match[1]
					remarks = fmt.Sprintf("CCTV%s|CCTV-%s", num, num)

				case reAlpha.MatchString(upper):
					match := reAlpha.FindStringSubmatch(upper)
					suffix := match[1]
					remarks = fmt.Sprintf("CCTV%s|CCTV-%s", suffix, suffix)
				}
			}
			epgs = append(epgs, models.IptvEpg{
				Name:    list.Remarks + "-" + channel.DisplayName.Value,
				Status:  1,
				Remarks: remarks,
			})
		}
		if len(epgs) > 0 {
			dao.DB.Model(&models.IptvEpgList{}).Where("id = ?", list.ID).Updates(&models.IptvEpgList{Status: 1, LastTime: time.Now().Unix()})
			dao.DB.Model(&models.IptvEpg{}).Where("name like ?", list.Remarks+"-%").Delete(&models.IptvEpg{})
			dao.DB.Model(&models.IptvEpg{}).Create(&epgs)
			CleanMealsXmlCacheAll() // 清空缓存
			go BindChannel()        // 绑定频道
			return true
		}
		return false
	}
	return false
}

func BindChannel() bool {
	ClearBind() // 清空绑定
	var channelList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Select("distinct name").Order("category,id").Find(&channelList).Error; err != nil {
		return false
	}

	var epgList []models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Find(&epgList).Error; err != nil {
		return false
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
	go GetCCTVChannelList(true)
	go GetProvinceChannelList(true)
	go CleanMealsXmlCacheAll()
	return true
}

func ClearBind() {
	dao.DB.Model(&models.IptvEpg{}).Where("content != ''").Update("content", "")
}
