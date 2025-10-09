package until

import (
	"fmt"
	"go-iptv/dao"
	"go-iptv/models"
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
