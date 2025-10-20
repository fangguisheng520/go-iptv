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
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ConvertCntvToXml(cntv dto.CntvJsonChannel, eName string) dto.XmlTV {
	tv := dto.XmlTV{
		GeneratorName: "清和IPTV管理系统",
		GeneratorURL:  "https://www.qingh.xyz",
	}

	// 添加频道
	tv.Channels = append(tv.Channels, dto.XmlChannel{
		ID: eName,
		DisplayName: []dto.DisplayName{
			{Lang: "zh",
				Value: eName,
			},
		},
	})

	// 添加节目表
	for _, p := range cntv.Program {
		start := time.Unix(p.StartTime, 0).UTC().Format("20060102150405 -0700")
		stop := time.Unix(p.EndTime, 0).UTC().Format("20060102150405 -0700")

		tv.Programmes = append(tv.Programmes, dto.Programme{
			Start:   start,
			Stop:    stop,
			Channel: eName,
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

	MakeMealsXmlCacheAll()
}

func MakeMealsXmlCacheAll() {
	var meals []models.IptvMeals
	dao.DB.Model(&models.IptvMeals{}).Where("status = 1").Find(&meals)
	log.Println("重建套餐EPG订阅缓存")
	for _, meal := range meals {
		GetEpg(meal.ID)
	}
	log.Println("重建套餐EPG订阅缓存完成")
}

func CleanMealsXmlCacheOne(id int64) {
	dao.Cache.Delete("rssEpgXml_" + strconv.FormatInt(id, 10))
	MakeMealsXmlCacheAll()
}

func UpdataEpgList() bool {
	var epgLists []models.IptvEpgList
	dao.DB.Model(&models.IptvEpgList{}).Find(&epgLists)
	for _, list := range epgLists {
		cacheKey := "epgXmlFrom_" + list.Name
		dao.Cache.Delete(cacheKey)
		xmlStr := GetUrlData(url.QueryEscape(strings.TrimSpace(list.Url)), list.UA)
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
				remarks := channel.DisplayName[0].Value
				upper := strings.ToUpper(remarks)
				if strings.Contains(upper, "CCTV") {
					switch {
					case reNum.MatchString(upper):
						match := reNum.FindStringSubmatch(upper)
						num := match[1]
						remarks = fmt.Sprintf("CCTV%s|CCTV-%s|CCTV%s 4K|CCTV-%s 4K|CCTV%s HD|CCTV-%s HD", num, num, num, num, num, num)

					case reAlpha.MatchString(upper):
						match := reAlpha.FindStringSubmatch(upper)
						suffix := match[1]
						remarks = fmt.Sprintf("CCTV%s|CCTV-%s", suffix, suffix)
					}
				} else {
					remarks = fmt.Sprintf("%s|%s 4K|%s HD", remarks, remarks, remarks)
				}
				epgs = append(epgs, models.IptvEpg{
					Name:    list.Remarks + "-" + channel.DisplayName[0].Value,
					Status:  1,
					Remarks: remarks,
				})
			}
			if len(epgs) > 0 {
				dao.DB.Model(&models.IptvEpgList{}).Where("id = ?", list.ID).Updates(&models.IptvEpgList{Status: 1, LastTime: time.Now().Unix()})
				// dao.DB.Model(&models.IptvEpg{}).Where("name like ?", list.Remarks+"-%").Delete(&models.IptvEpg{})
				// dao.DB.Model(&models.IptvEpg{}).Create(&epgs)
				SyncEpgs(list.Remarks, epgs) // 同步
				go BindChannel()
				// CleanMealsXmlCacheAll() // 清除缓存
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
	xmlStr := GetUrlData(url.QueryEscape(strings.TrimSpace(list.Url)), list.UA)
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
			remarks := channel.DisplayName[0].Value
			if remarks == "" {
				continue
			}
			upper := strings.ToUpper(remarks)
			if strings.Contains(upper, "CCTV") {
				switch {
				case reNum.MatchString(upper):
					match := reNum.FindStringSubmatch(upper)
					num := match[1]
					remarks = fmt.Sprintf("CCTV%s|CCTV-%s|CCTV%s 4K|CCTV-%s 4K|CCTV%s HD|CCTV-%s HD", num, num, num, num, num, num)

				case reAlpha.MatchString(upper):
					match := reAlpha.FindStringSubmatch(upper)
					suffix := match[1]
					remarks = fmt.Sprintf("CCTV%s|CCTV-%s", suffix, suffix)
				}
			} else {
				remarks = fmt.Sprintf("%s|%s 4K|%s HD", remarks, remarks, remarks)
			}

			epgs = append(epgs, models.IptvEpg{
				Name:    list.Remarks + "-" + channel.DisplayName[0].Value,
				Status:  1,
				Remarks: remarks,
			})
		}
		if len(epgs) > 0 {
			dao.DB.Model(&models.IptvEpgList{}).Where("id = ?", list.ID).Updates(&models.IptvEpgList{Status: 1, LastTime: time.Now().Unix()})
			// dao.DB.Model(&models.IptvEpg{}).Where("name like ?", list.Remarks+"-%").Delete(&models.IptvEpg{})
			// dao.DB.Model(&models.IptvEpg{}).Create(&epgs)
			SyncEpgs(list.Remarks, epgs) // 同步
			// CleanMealsXmlCacheAll()      // 清空缓存
			go BindChannel() // 绑定频道
			return true
		}
		return false
	}
	return false
}

func BindChannel() bool {
	// ClearBind() // 清空绑定
	var channelList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Select("distinct name").Where("status = 1").Order("category,id").Find(&channelList).Error; err != nil {
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
		epgData.Content = strings.Join(MergeAndUnique(strings.Split(epgData.Content, ","), tmpList), ",")
		if epgData.Content != "" {
			dao.DB.Save(&epgData)
		}
	}
	GetCCTVChannelList(true)
	GetProvinceChannelList(true)
	go CleanMealsXmlCacheAll()
	// go MakeMealsXmlCacheAll()
	return true
}

// SyncEpgs 同步 IPTV EPG 数据：
// - 保留数据库中已存在的记录（不更新）
// - 新数据中有但数据库没有的 → 新增
// - 数据库中有但新数据中没有的 → 删除
func SyncEpgs(prefix string, epgs []models.IptvEpg) error {
	// 1. 查询数据库中已有的记录
	var oldEpgs []models.IptvEpg
	if err := dao.DB.Where("name LIKE ?", prefix+"-%").Find(&oldEpgs).Error; err != nil {
		return err
	}

	// 2. 建立 name 映射方便比对
	oldMap := make(map[string]bool)
	for _, o := range oldEpgs {
		oldMap[o.Name] = true
	}

	newMap := make(map[string]bool)
	for _, n := range epgs {
		newMap[n.Name] = true
	}

	// 3. 计算需要新增与删除的数据
	var toAdd []models.IptvEpg
	var toDelete []string

	for _, n := range epgs {
		if !oldMap[n.Name] {
			toAdd = append(toAdd, n)
		}
	}

	for _, o := range oldEpgs {
		if !newMap[o.Name] {
			toDelete = append(toDelete, o.Name)
		}
	}

	// 4. 执行数据库操作
	if len(toDelete) > 0 {
		if err := dao.DB.Where("name IN ?", toDelete).Delete(&models.IptvEpg{}).Error; err != nil {
			return err
		}
		log.Printf("删除 %d 条无效 EPG 记录\n", len(toDelete))
	}

	if len(toAdd) > 0 {
		if err := dao.DB.Create(&toAdd).Error; err != nil {
			return err
		}
		log.Printf("新增 %d 条 EPG 记录\n", len(toAdd))
	}

	log.Printf("同步完成：新增 %d，删除 %d\n", len(toAdd), len(toDelete))
	return nil
}

func GetTxt(id int64) string {
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
			data := GetCCTVChannelList(false)
			if data != "" {
				res += category.Name + ",#genre#\n"
				res += data

			}
		case -1:
			data := GetProvinceChannelList(false)
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

func GetEpg(id int64) dto.XmlTV {

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
			// output, _ := xml.MarshalIndent(tvData, "", "  ")
			// log.Println(string(output))

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
		if !seen[ch.ID] {
			seen[ch.ID] = true
			ids[ch.ID] = i
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
			t, err := time.Parse("20060102150405 -0700", p.Start)
			if err != nil {
				log.Println("解析时间错误:", err)
				continue
			}
			key := p.Channel + "_" + fmt.Sprintf("%d", t.Unix()) + "_" + p.Title.Value // 唯一键

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
		eName := strings.SplitN(epg.Name, "-", 2)[1]
		nameList := strings.Split(epg.Content, ",")
		var channelList []models.IptvChannel
		if err := dao.DB.Model(&models.IptvChannel{}).Where("name in (?)", nameList).Order("sort asc").Find(&channelList).Error; err != nil {
			continue
		}
		dName := []dto.DisplayName{}
		exists := false
		for _, channel := range channelList {
			tmpData, err := GetEpgCntv(eName)
			if err == nil {
				tmpXml := ConvertCntvToXml(tmpData, eName)
				for k, c := range cntvXml.Channels {
					if c.ID == eName {
						exists = true
						dName = append(c.DisplayName, dto.DisplayName{
							Lang:  "zh",
							Value: channel.Name,
						})
						cntvXml.Channels[k].DisplayName = dName
					}
				}

				if !exists {
					dName = append(dName, dto.DisplayName{
						Lang:  "zh",
						Value: channel.Name,
					})
					cntvXml.Channels = append(cntvXml.Channels, dto.XmlChannel{
						ID:          eName,
						DisplayName: dName,
					})
				}

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
	if err := dao.DB.Model(&models.IptvEpg{}).Where("name like ? and status = 1", "%"+"-%卫视%").Find(&epgs).Error; err != nil {
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
		eFrom := strings.SplitN(epg.Name, "-", 2)[0]
		eName := strings.SplitN(epg.Name, "-", 2)[1]

		var epgList models.IptvEpgList
		if err := dao.DB.Model(&models.IptvEpgList{}).Where("remarks = ? and status = 1", eFrom).First(&epgList).Error; err != nil {
			continue
		}
		dName := []dto.DisplayName{}
		exists := false
		for _, channel := range channelList {
			tmpXml := GetEpgListXml(epgList.Name, epgList.Url)

			for k, c := range epgXml.Channels {
				if c.ID == eName {
					exists = true
					dName = append(c.DisplayName, dto.DisplayName{
						Lang:  "zh",
						Value: channel.Name,
					})
					epgXml.Channels[k].DisplayName = dName
				}
			}

			if !exists {
				dName = append(dName, dto.DisplayName{
					Lang:  "zh",
					Value: channel.Name,
				})
				epgXml.Channels = append(epgXml.Channels, dto.XmlChannel{
					ID:          eName,
					DisplayName: dName,
				})
			}

			var cId string
			for _, c := range tmpXml.Channels {
				if c.DisplayName[0].Value == eName {
					cId = c.ID
					break
				}
			}

			for _, p := range tmpXml.Programmes {
				if p.Channel == cId {
					p.Channel = eName
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
			dName := []dto.DisplayName{}
			exists := false
			if eType == "cntv" {
				tmpData, err := GetEpgCntv(eName)
				if err == nil {
					tmpXml := ConvertCntvToXml(tmpData, eName)
					for k, c := range epgXml.Channels {
						if c.ID == eName {
							exists = true
							dName = append(c.DisplayName, dto.DisplayName{
								Lang:  "zh",
								Value: channel.Name,
							})
							epgXml.Channels[k].DisplayName = dName
						}
					}

					if !exists {
						dName = append(dName, dto.DisplayName{
							Lang:  "zh",
							Value: channel.Name,
						})
						epgXml.Channels = append(epgXml.Channels, dto.XmlChannel{
							ID:          eName,
							DisplayName: dName,
						})
					}

					for _, p := range tmpXml.Programmes {
						p.Channel = eName
						epgXml.Programmes = append(epgXml.Programmes, p)
					}
					if len(epgXml.Channels) > 0 && len(epgXml.Programmes) > 0 {
						break
					}
					continue
				}
			}

			var epgList models.IptvEpgList
			if err := dao.DB.Model(&models.IptvEpgList{}).Where("remarks = ? and status = 1", eType).First(&epgList).Error; err != nil {
				continue
			}
			tmpXml := GetEpgListXml(epgList.Name, epgList.Url)
			for k, c := range epgXml.Channels {
				if c.ID == eName {
					exists = true
					dName = append(c.DisplayName, dto.DisplayName{
						Lang:  "zh",
						Value: channel.Name,
					})
					epgXml.Channels[k].DisplayName = dName
				}
			}

			if !exists {
				dName = append(dName, dto.DisplayName{
					Lang:  "zh",
					Value: channel.Name,
				})
				epgXml.Channels = append(epgXml.Channels, dto.XmlChannel{
					ID:          eName,
					DisplayName: dName,
				})
			}

			var cId string
			for _, c := range tmpXml.Channels {
				if c.DisplayName[0].Value == eName {
					cId = c.ID
					break
				}
			}

			for _, p := range tmpXml.Programmes {
				if p.Channel == cId {
					p.Channel = eName
					epgXml.Programmes = append(epgXml.Programmes, p)
				}
			}
		}
	}

	return epgXml
}
