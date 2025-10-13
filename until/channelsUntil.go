package until

import (
	"bufio"
	"fmt"
	"go-iptv/dao"
	"go-iptv/models"
	"log"
	"regexp"
	"strings"
)

// convertListFormat 将 m3u 或 "频道,URL" 格式统一转换为 "频道,URL\n"
func ConvertListFormat(srclist string) string {
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
			if idx := strings.Index(channelName, " "); idx != -1 {
				channelName = channelName[:idx]
			}
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
			if idx := strings.Index(channelName, " "); idx != -1 {
				channelName = channelName[:idx]
			}
			channelURL := match[2]
			convertedList.WriteString(fmt.Sprintf("%s,%s\n", channelName, channelURL))
		}
		return convertedList.String()
	}

	return srclist
}

// addChannelList 添加频道到数据库

func ConvertDataToMap(data string) map[string]string {
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

func GetCCTVChannelList(rebuild bool) string {
	var res string
	var channelCache = "cctv_channel_list"
	if dao.Cache.ChannelExists(channelCache) && !rebuild {
		data, err := dao.Cache.Get(channelCache)
		if err == nil {
			return string(data)
		}
	}

	var epgs []models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Where("name like ? and status = 1", "cntv-%").Find(&epgs).Error; err != nil {
		return res
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
		for _, channel := range channelList {
			res += fmt.Sprintf("%s,%s\n", channel.Name, channel.Url)
		}
	}
	go dao.Cache.Set(channelCache, []byte(res))
	return res
}

func GetProvinceChannelList(rebuild bool) string {
	var res string
	var channelCache = "province_channel_list"
	if dao.Cache.ChannelExists(channelCache) && !rebuild {
		data, err := dao.Cache.Get(channelCache)
		if err == nil {
			return string(data)
		}
	}

	var epgs []models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Where("name like ? and status = 1", "51zmt-%卫视").Find(&epgs).Error; err != nil {
		return res
	}
	var channelList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Find(&channelList).Error; err != nil {
		return res
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
		for _, channel := range channelList {
			res += fmt.Sprintf("%s,%s\n", channel.Name, channel.Url)
		}
	}
	go dao.Cache.Set(channelCache, []byte(res))
	return res
}

func Txt2M3u8(txtData, host, token string) string {

	epgURL := host + "/epg/" + token + "/e.xml" // ✅ 可自行修改 EPG 地址
	logoBase := host + "/logo/"                 // ✅ 可自行修改 logo 前缀

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("#EXTM3U url-tvg=\"%s\"\n\n", epgURL))

	scanner := bufio.NewScanner(strings.NewReader(txtData))
	currentGroup := "未分组"
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 检查是否为分组行（如 “中央台,#genre#”）
		if strings.HasSuffix(line, "#genre#") {
			group := strings.TrimSuffix(line, ",#genre#")
			currentGroup = strings.TrimSpace(group)
			continue
		}

		// 普通频道行
		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			fmt.Printf("Txt2M3u8: 第 %d 行格式错误: %s\n", lineNum, line)
			continue
		}

		name := strings.TrimSpace(parts[0])
		url := strings.TrimSpace(parts[1])
		epgName := GetEpgName(name)
		var logo string
		if epgName != "" {
			logo = fmt.Sprintf("%s%s.png", strings.TrimRight(logoBase, "/")+"/", epgName)
		}

		// ✅ 生成 #EXTINF 信息
		extinf := fmt.Sprintf(`#EXTINF:-1 tvg-id="%s" tvg-name="%s" tvg-logo="%s" group-title="%s",%s`,
			name, name, logo, currentGroup, name)
		builder.WriteString(extinf + "\n")
		builder.WriteString(url + "\n\n")
	}

	if err := scanner.Err(); err != nil {
		log.Println("Txt2M3u8: m3u8解析出错:", err)
	}

	return builder.String()
}

func M3UToGenreTXT(m3u string) string {
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
			log.Println(matches)
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

func GetEpgName(name string) string {
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

	return strings.Split(epgName, "-")[1]
}

func IsM3UContent(data string) bool {
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
