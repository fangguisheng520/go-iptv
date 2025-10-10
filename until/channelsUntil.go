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
	if err := dao.DB.Model(&models.IptvEpg{}).Where("name like ?", "cntv-%").Find(&epgs).Error; err != nil {
		return res
	}
	var channelList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Find(&channelList).Error; err != nil {
		return res
	}

	for _, epg := range epgs {
		nameList := strings.Split(epg.Remarks, "|")
		for _, channel := range channelList {
			for _, name := range nameList {
				if strings.EqualFold(channel.Name, name) {
					res += fmt.Sprintf("%s,%s\n", channel.Name, channel.Url)
				}
			}
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
	if err := dao.DB.Model(&models.IptvEpg{}).Where("name like ?", "51zmt-%卫视").Find(&epgs).Error; err != nil {
		return res
	}
	var channelList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Find(&channelList).Error; err != nil {
		return res
	}

	for _, epg := range epgs {
		nameList := strings.Split(epg.Remarks, "|")
		for _, channel := range channelList {
			for _, name := range nameList {
				if strings.EqualFold(channel.Name, name) {
					res += fmt.Sprintf("%s,%s\n", channel.Name, channel.Url)
				}
			}
		}
	}
	go dao.Cache.Set(channelCache, []byte(res))
	return res
}

func Txt2M3u8(txtData string) string {

	cfg := dao.GetConfig()

	epgURL := "https://epg.51zmt.top:8000/e.xml" // ✅ 可自行修改 EPG 地址
	logoBase := cfg.ServerUrl + "/logo/"         // ✅ 可自行修改 logo 前缀

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

func GetEpgName(name string) string {
	var epgs []models.IptvEpg
	dao.DB.Model(&models.IptvEpg{}).Where("content like ?", "%"+name+"%").Find(&epgs)

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
