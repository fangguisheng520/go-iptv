package until

import "strings"

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
