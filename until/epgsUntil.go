package until

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"strconv"
	"strings"
	"time"
)

// EpgName 根据 epg 前缀返回中文名
func EpgName(name string) string {

	var epgMap = map[string]string{
		"cntv": "CCTV官网",
		// "jisu":   "极速数据",
		// "tvming": "天脉聚源",
		// "tvmao":  "电视猫",
		// "tvsou":  "搜视网",
		"51zmt": "51zmt",
		// "112114": "112114",
	}
	parts := strings.SplitN(name, "-", 2)
	key := parts[0]
	if val, ok := epgMap[key]; ok {
		return val
	}
	return key // 没找到就返回原始 key
}

func ConvertCntvToXml(cntv dto.CntvJsonChannel, cName string) dto.TV {
	tv := dto.TV{
		GeneratorName: "清和IPTV管理系统",
		GeneratorURL:  "https://www.qingh.xyz",
	}

	// 添加频道
	tv.Channels = append(tv.Channels, dto.XmlChannel{
		ID: "1",
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
			Channel: "1",
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

func Get51zmtXml() dto.TV {
	epgUrl := "http://epg.51zmt.top:8000/e.xml"
	cacheKey := "epg51ZmtXml"
	var zmtTV dto.TV
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
	xml.Unmarshal(xmlByte, &zmtTV)
	return zmtTV
}

func GetEpgCntv(name string) (dto.CntvJsonChannel, error) {

	var cacheKey = "cntv_" + name

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
