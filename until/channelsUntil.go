package until

import (
	"fmt"
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
			channelName := match[1]
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
			channelName := match[1]
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

// AddOrUpdateGenreChannels 根据 map 更新或新增分类及频道
