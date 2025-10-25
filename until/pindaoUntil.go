package until

import (
	"fmt"
	"go-iptv/dao"
	"go-iptv/models"
	"log"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// convertListFormat 将 m3u 或 "频道,URL" 格式统一转换为 "频道,URL\n"
func ConvertListFormat_pindao(srclist string) string {
	if !strings.HasSuffix(srclist, "\n") {
		srclist += "\n"
	}

	var convertedList strings.Builder

	// 匹配 #EXTINF
	reExtInf := regexp.MustCompile(`#EXTINF:-?\d+.*?,(.*?)\n(.*?)\n`)
	matches := reExtInf.FindAllStringSubmatch(srclist, -1)

	if len(matches) > 0 {
		for _, match := range matches {
			channelName := strings.TrimSpace(match[1])
			// if idx := strings.Index(channelName, " "); idx != -1 {
			// 	channelName = channelName[:idx]
			// }
			channelURL := match[2]
			convertedList.WriteString(fmt.Sprintf("%s,%s\n", channelName, channelURL))
		}
		return convertedList.String()
	}

	// 匹配 "频道,URL"
	reLine := regexp.MustCompile(`(.*?),(.*)\n`)
	matches = reLine.FindAllStringSubmatch(srclist, -1)

	if len(matches) > 0 {
		for _, match := range matches {
			channelName := strings.TrimSpace(match[1])
			// if idx := strings.Index(channelName, " "); idx != -1 {
			// 	channelName = channelName[:idx]
			// }
			channelURL := match[2]
			convertedList.WriteString(fmt.Sprintf("%s,%s\n", channelName, channelURL))
		}
		return convertedList.String()
	}

	return srclist
}

// addChannelList 添加频道到数据库

func ConvertDataToMap_pindao(data string) map[string]string {
	lines := strings.Split(data, "\n")
	result := make(map[string]string)
	currentGenre := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "#genre#") {
			currentGenre = strings.ReplaceAll(line, ",#genre#", "")
			result[currentGenre] = ""
		} else if currentGenre != "" {
			result[currentGenre] += line + "\n"
		}
	}

	for k, v := range result {
		result[k] = strings.TrimSpace(v)
	}

	return result
}

func M3UToGenreTXT_pindao(m3u string) string {
	lines := strings.Split(m3u, "\n")

	genreMap := make(map[string][]string)
	var groupsOrder []string // 记录首次出现的分组顺序

	// 更稳健的正则：在任意位置捕获 group-title="xx"，最后一个逗号后是频道名
	reExtinf := regexp.MustCompile(`(?i)#EXTINF:[^,]*?(?:.*?group-title=["']([^"']+)["'])?.*?,\s*(.*)$`)

	var lastGroup, lastName string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#EXTM3U") {
			continue
		}

		if strings.HasPrefix(line, "#EXTINF:") {
			matches := reExtinf.FindStringSubmatch(line)
			if len(matches) >= 3 {
				group := strings.TrimSpace(matches[1])
				name := strings.TrimSpace(matches[2])

				if group == "" {
					group = "未分组"
				}

				lastGroup = group
				lastName = name

				// 若首次见到该分组，记录顺序
				if _, ok := genreMap[group]; !ok {
					groupsOrder = append(groupsOrder, group)
					genreMap[group] = []string{}
				}
			}
		} else if strings.HasPrefix(line, "http") || strings.HasPrefix(line, "rtsp") || strings.HasPrefix(line, "rtmp") {
			if lastName != "" && lastGroup != "" {
				genreMap[lastGroup] = append(genreMap[lastGroup], fmt.Sprintf("%s,%s", lastName, line))
				// 清空以避免错误关联
				lastName, lastGroup = "", ""
			}
		}
	}

	// 按首次出现顺序输出（避免 sort 后改变顺序）
	var builder strings.Builder
	for _, group := range groupsOrder {
		builder.WriteString(fmt.Sprintf("%s,#genre#\n", group))
		for _, item := range genreMap[group] {
			builder.WriteString(item + "\n")
		}
		builder.WriteString("\n")
	}

	return strings.TrimSpace(builder.String())
}

func GetEpgName_pindao(name string) string {
	var epgs []models.IptvEpg
	dao.DB.Model(&models.IptvEpg{}).Where("content like ? and status = 1", "%"+name+"%").Find(&epgs)

	var epgName string
	for _, epg := range epgs {
		for _, v := range strings.Split(epg.Content, ",") {
			if strings.EqualFold(name, v) {
				epgName = epg.Name
				break
			}
		}
		if epgName != "" {
			break
		}
	}

	if epgName == "" {
		return epgName
	}

	return strings.SplitN(epgName, "-", 2)[1]
}

func IsM3UContent_pindao(data string) bool {
	// 去除前后空白
	trimmed := strings.TrimSpace(data)

	// 必须以 #EXTM3U 开头
	if !strings.HasPrefix(trimmed, "#EXTM3U") {
		return false
	}

	// 检查是否包含至少一个 #EXTINF
	if !strings.Contains(data, "#EXTINF:") {
		return false
	}

	return true
}

func GetAutoChannelList_pindao(category models.IptvFenlei) []models.IptvChannelShow {

	var result []models.IptvChannelShow

	autoCaCheKey := "autoCategory_" + strconv.FormatInt(category.ID, 10)
	if dao.Cache.Exists(autoCaCheKey) {
		err := dao.Cache.GetStruct(autoCaCheKey, result)
		if err == nil {
			return result
		}
	}

	var channelList []models.IptvChannelShow
	if err := dao.DB.Table(models.IptvChannelShow{}.TableName() + " AS c").
		Select("c.*, e.name AS epg_name").
		Joins("LEFT JOIN " + models.IptvEpg{}.TableName() + " AS e ON c.e_id = e.id AND e.status = 1").
		Where("c.e_id != 0 and c.status = 1").
		Order("c_id,sort asc").
		Find(&channelList).Error; err != nil {
		log.Println("获取频道列表失败:", err)
		return result
	}

	re := regexp.MustCompile(category.Rules)

	for _, ch := range channelList {
		if strings.Contains(ch.Name, category.Rules) {
			if ch.EpgName != "" {
				ch.Logo = EpgNameGetLogo(ch.EpgName)
			}
			result = append(result, ch)
			continue
		}
		if re.MatchString(ch.Name) {
			if ch.EpgName != "" {
				ch.Logo = EpgNameGetLogo(ch.EpgName)
			}
			result = append(result, ch)
		}
	}

	if err := dao.Cache.SetStruct(autoCaCheKey, result); err != nil {
		log.Println("epg缓存设置失败:", err)
		dao.Cache.Delete(autoCaCheKey)
	}

	return result
}

func CaGetChannels_pindao(categoryDb models.IptvFenlei) []models.IptvPindaoShow {
	var channels []models.IptvPindaoShow
	dao.DB.Table(models.IptvPindaoShow{}.TableName()+" AS c").
		Select("c.*, e.name AS epg_name").
		Joins("LEFT JOIN "+models.IptvEpg{}.TableName()+" AS e ON c.e_id = e.id AND e.status = 1").
		Where("c.c_id = ?", categoryDb.ID).
		Order("sort asc").
		Find(&channels)
	println("打印频道列表：")
	for i, ch := range channels {
		println(channels[i].Name, channels[i].Rss_Rules, channels[i].Epg_Rules)
		//println(ch.Name)
		channels[i].Logo = EpgNameGetLogo(ch.Name)
	}
	println("------------")
	return channels

	//if categoryDb.Type == "auto" {
	//	return GetAutoChannelList_pindao(categoryDb)
	//} else {
	//	var channels []models.IptvChannelShow
	//	dao.DB.Table(models.IptvChannelShow{}.TableName()+" AS c").
	//		Select("c.*, e.name AS epg_name").
	//		Joins("LEFT JOIN "+models.IptvEpg{}.TableName()+" AS e ON c.e_id = e.id AND e.status = 1").
	//		Where("c.c_id = ?", categoryDb.ID).
	//		Order("sort asc").
	//		Find(&channels)
	//	for i, ch := range channels {
	//		println(channels[i].Name)
	//		if ch.EpgName != "" {
	//			channels[i].Logo = EpgNameGetLogo(ch.EpgName)
	//		}
	//	}
	//	return channels
	//}

}

func AddChannelList_pindao(srclist string, cId, listId int64, doRepeat bool) (int, error) {
	if srclist == "" {
		// 如果 srclist 为空，删除当前分类下所有数据
		if err := dao.DB.Transaction(func(tx *gorm.DB) error {
			return tx.Delete(&models.IptvChannel{}, "c_id = ?", cId).Error
		}); err != nil {
			return 0, err
		}
		BindChannel()
		return 0, nil
	}

	// 转换为 "频道,URL" 格式
	srclist = ConvertListFormat(srclist)

	// 获取 cname 分类下已有的频道
	var oldChannels []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Where("c_id = ?", cId).Find(&oldChannels).Error; err != nil {
		return 0, err
	}

	// 当前分类已有 URL -> channelName（大小写敏感）
	existMap := make(map[string]string)
	for _, ch := range oldChannels {
		if ch.Url != "" && ch.Name != "" {
			existMap[ch.Url] = ch.Name
		}
	}

	existHandMap := make(map[string]string)
	if doRepeat {
		var handChannels []models.IptvChannel
		dao.DB.Table(models.IptvChannel{}.TableName()+" AS c").
			Select("c.name, c.url").
			Joins("LEFT JOIN "+models.IptvCategory{}.TableName()+" AS cat ON c.c_id = cat.id and cat.enable = 1").
			Where("cat.type = ?", "user").
			Scan(&handChannels)

		for _, ch := range handChannels {
			if ch.Url != "" && ch.Name != "" {
				existHandMap[ch.Url] = ch.Name
			}
		}
	}

	// 正则清洗
	reSpaces := regexp.MustCompile(`\s+`)
	reGenre := regexp.MustCompile(`#genre#`)
	reVer := regexp.MustCompile(`ver\..*?\.m3u8`)
	reTme := regexp.MustCompile(`t\.me.*?\.m3u8`)
	reBbsok := regexp.MustCompile(`https(.*)www\.bbsok\.cf[^>]*`)

	lines := strings.Split(srclist, "\n")
	newChannels := make([]models.IptvChannel, 0)
	srclistUrls := make(map[string]struct{})
	repetNum := 0
	delIDs := make([]int64, 0)
	var sortIndex int64 = 1
	var rawCount int64 = 0

	// 先处理循环，准备新增和标记要删除的旧数据
	for _, line := range lines {
		line = strings.ReplaceAll(line, " ,", ",")
		line = strings.ReplaceAll(line, "\r", "")
		line = reSpaces.ReplaceAllString(line, "")
		line = reGenre.ReplaceAllString(line, "")
		line = reVer.ReplaceAllString(line, "")
		line = reTme.ReplaceAllString(line, "")
		line = reBbsok.ReplaceAllString(line, "")

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "http") {
			if _, ok := srclistUrls[line]; ok {
				repetNum++
				continue
			}
			srclistUrls[line] = struct{}{}
		}

		parts := strings.SplitN(line, ",", 2)
		channelName := parts[0]

		var chStatus int64 = 1
		if strings.Contains(channelName, "|") {
			tmp := strings.SplitN(channelName, "|", 2)
			if tmp[0] == "0" {
				chStatus = 0
			}
			channelName = tmp[1]
		}

		source := parts[1]

		srcList := strings.Split(source, "#")

		for _, src := range srcList {
			src2 := strings.Trim(strings.NewReplacer(`"`, "", "'", "", "}", "", "{", "").Replace(src), " \r\n\t")
			if src2 == "" || channelName == "" {
				continue
			}
			rawCount++

			srclistUrls[src2] = struct{}{}

			if doRepeat {
				if _, exists := existHandMap[src2]; exists {
					for _, ch := range oldChannels {
						if ch.Url == src2 {
							delIDs = append(delIDs, ch.ID)
						}
					}
					repetNum++
					continue
				}
			}

			if oldName, exists := existMap[src2]; exists {
				if oldName != channelName {
					// URL 相同但 channelName 不同 → 删除旧数据
					for _, ch := range oldChannels {
						if ch.Url == src2 {
							delIDs = append(delIDs, ch.ID)
						}
					}
				} else {
					// URL + channelName 相同 → 检查顺序
					for _, ch := range oldChannels {
						if ch.Url == src2 && ch.Name == channelName && ch.Sort != sortIndex || ch.Status != chStatus {
							ch.Sort = sortIndex
							if err := dao.DB.Model(&models.IptvChannel{}).
								Where("id = ?", ch.ID).
								Updates(map[string]interface{}{
									"sort":   sortIndex,
									"status": chStatus,
								}).Error; err != nil {
								log.Println("更新顺序失败:", err)
							}
							break
						}
					}
					sortIndex++
					continue
				}
			}

			// 新增数据
			newChannels = append(newChannels, models.IptvChannel{
				Name:   channelName,
				Url:    src2,
				CId:    cId,
				ListId: listId,
				Sort:   sortIndex,
				Status: chStatus,
			})
			existMap[src2] = channelName
			sortIndex++
		}
	}

	// 批量删除数据库中当前分类但新列表中没有的 URL
	for _, ch := range oldChannels {
		if _, ok := srclistUrls[ch.Url]; !ok {
			delIDs = append(delIDs, ch.ID)
		}
	}

	// 在事务中执行删除和新增
	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if len(delIDs) > 0 {
			if err := tx.Delete(&models.IptvChannel{}, delIDs).Error; err != nil {
				return err
			}
		}
		if len(newChannels) > 0 {
			if err := tx.Create(&newChannels).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return repetNum, err
	}

	// 只有当有新增或删除时才执行异步更新
	if len(newChannels) > 0 || len(delIDs) > 0 {
		BindChannel()
	}
	log.Printf("订阅频道数量: %d", rawCount) // 新增日志输出
	dao.DB.Model(&models.IptvCategory{}).Where("id = ?", cId).Update("rawcount", rawCount)
	return repetNum, nil
}
