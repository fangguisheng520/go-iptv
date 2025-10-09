package crontab

import (
	"fmt"
	"go-iptv/dao"
	"go-iptv/models"
	"go-iptv/until"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var CrontabStatus bool = false
var UpdateStatus bool = false
var StopChan = make(chan struct{})

func Crontab() {
	defer func() { CrontabStatus = false }()
	if CrontabStatus {
		log.Println("定时任务服务已启动，请勿重复启动")
		return
	}
	cfg := dao.GetConfig()
	autoUpdate := cfg.Channel.Auto
	upInterval := cfg.Channel.Interval
	if autoUpdate == 1 && upInterval > 0 {
		log.Println("定时任务服务启动...")
		CrontabStatus = true
		interval := time.Duration(upInterval) * time.Minute
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				if UpdateStatus {
					log.Println("正在更新频道，请稍后...")
					continue
				}
				log.Println("开始执行更新频道任务：", t.Format("2006-01-02 15:04:05"))
				UpdateList()
			case <-StopChan:
				log.Println("收到停止信号，停止更新频道任务")
				ticker.Stop()
				return
			}
		}
	} else {
		log.Println("定时任务服务未开启...")
	}
}

func UpdateList() {
	UpdateStatus = true
	defer func() { UpdateStatus = false }()
	// TODO: 定时任务
	var lists []models.IptvCategory
	res := dao.DB.Model(&models.IptvCategory{}).Where("url != ?", "").Find(&lists)

	if res.RowsAffected == 0 {
		log.Println("没有可更新的频道列表")
		return
	}
	for _, v := range lists {
		resp, err := http.Get(v.Url)
		if err != nil {
			log.Println("更新频道列表失败--->链接响应失败1：", v.Url)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Println("更新频道列表失败--->链接响应失败2：", v.Url)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("更新频道列表失败--->读取响应失败：", v.Url)
			return
		}

		urlData := until.FilterEmoji(string(body)) // 过滤emoji表情

		if v.AutoCategory == 1 {
			if !strings.Contains(urlData, "#genre#") {
				dao.DB.Model(&models.IptvCategory{}).Where("name = ?", v.Name).Updates(map[string]interface{}{
					"latesttime":   time.Now().Format("2006-01-02 15:04:05"),
					"autocategory": 0,
				})
				AddChannelList(v.Name, urlData)
			}
			GenreChannels(v.Name, urlData)
		} else {
			dao.DB.Model(&models.IptvCategory{}).Where("name = ?", v.Name).Updates(map[string]interface{}{
				"latesttime": time.Now().Format("2006-01-02 15:04:05"),
			})
			AddChannelList(v.Name, urlData)
		}
	}
	log.Println("定时执行更新频道任务结束")
}

func GenreChannels(listName, srclist string) {

	data := until.ConvertDataToMap(srclist)

	for genreName, genreList := range data {
		genreName = strings.TrimSpace(genreName)
		if genreName == "" {
			continue
		}

		categoryName := strings.ReplaceAll(fmt.Sprintf("%s(%s)", genreName, listName), " ", "")

		var count int64
		if err := dao.DB.Model(&models.IptvCategory{}).
			Where("name = ?", categoryName).
			Count(&count).Error; err != nil {
			continue
		}

		if count > 0 {
			dao.DB.Model(&models.IptvCategory{}).Where("name = ?", categoryName).Updates(map[string]interface{}{
				"latesttime": time.Now().Format("2006-01-02 15:04:05"),
			})
			AddChannelList(categoryName, genreList)
		} else {
			var maxSort int
			dao.DB.Model(&models.IptvCategory{}).Select("IFNULL(MAX(sort),0)").Scan(&maxSort)
			newCategory := models.IptvCategory{
				LatestTime: time.Now().Format("2006-01-02 15:04:05"),
				Name:       categoryName,
				Sort:       maxSort + 1,
				Type:       "add",
			}

			if err := dao.DB.Create(&newCategory).Error; err != nil {
				log.Println("创建分类"+categoryName+"失败", err)
				return
			}

			AddChannelList(categoryName, genreList)
		}
	}
	log.Println("更新" + listName + "分类结束")
}

func AddChannelList(cname, srclist string) {
	if cname == "" || srclist == "" {
		log.Println("分类名或频道列表不能为空")
	}

	// 转换为 "频道,URL" 格式
	srclist = until.ConvertListFormat(srclist)

	// 获取 cname 分类下已有的频道
	var oldChannels []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Where("category = ?", cname).Find(&oldChannels).Error; err != nil {
		log.Println("分类"+cname+"获取已有频道失败", err)
	}

	// 获取数据库中所有 URL，用于跨分类去重
	var allChannels []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Find(&allChannels).Error; err != nil {
		log.Println("获取已有频道失败", err)
	}

	allUrls := make(map[string]string)               // url -> category
	existUrls := make(map[string]models.IptvChannel) // 当前分类已有 URL
	for _, ch := range allChannels {
		if ch.Url != "" {
			allUrls[ch.Url] = ch.Category
		}
	}
	for _, ch := range oldChannels {
		if ch.Url != "" {
			existUrls[ch.Url] = ch
		}
	}

	// 预处理正则清洗
	reSpaces := regexp.MustCompile(`\s+`)
	reGenre := regexp.MustCompile(`#genre#`)
	reVer := regexp.MustCompile(`ver\..*?\.m3u8`)
	reTme := regexp.MustCompile(`t\.me.*?\.m3u8`)
	reBbsok := regexp.MustCompile(`https(.*)www\.bbsok\.cf[^>]*`)

	lines := strings.Split(srclist, "\n")
	srclistUrls := make(map[string]struct{}) // 记录新列表 URL
	var newChannels []models.IptvChannel     // 待批量新增
	repetNum := 0

	for _, line := range lines {
		line = strings.ReplaceAll(line, " ,", ",")
		line = strings.ReplaceAll(line, "\r", "")
		line = reSpaces.ReplaceAllString(line, "")
		line = reGenre.ReplaceAllString(line, "")
		line = reVer.ReplaceAllString(line, "")
		line = reTme.ReplaceAllString(line, "")
		line = reBbsok.ReplaceAllString(line, "")

		if !strings.Contains(line, ",") {
			continue
		}

		parts := strings.SplitN(line, ",", 2)
		channelName := parts[0]
		source := parts[1]

		srcList := strings.Split(source, "#")

		for _, src := range srcList {
			src2 := strings.NewReplacer(`"`, "", "'", "", "}", "", "{", "").Replace(src)
			if src2 == "" || channelName == "" {
				continue
			}

			srclistUrls[src2] = struct{}{}

			// 如果 URL 已在其它分类中存在，则跳过
			if cat, ok := allUrls[src2]; ok && cat != cname {
				repetNum++
				continue
			}

			// 如果 URL 已在当前分类中，则保留
			if _, exists := existUrls[src2]; exists {
				continue
			}

			// 新增数据
			newChannels = append(newChannels, models.IptvChannel{
				Name:     channelName,
				Url:      src2,
				Category: cname,
			})
			existUrls[src2] = models.IptvChannel{Name: channelName, Url: src2, Category: cname}
		}
	}

	// 批量新增
	if len(newChannels) > 0 {
		if err := dao.DB.Model(&models.IptvChannel{}).Create(&newChannels).Error; err != nil {
			log.Println("新增数据失败：", err)
		}
	}

	// 批量删除：删除 cname 分类下不在 srclistUrls 中的旧数据
	var delIDs []int64
	for url, ch := range existUrls {
		if _, ok := srclistUrls[url]; !ok {
			delIDs = append(delIDs, ch.ID)
		}
	}
	if len(delIDs) > 0 {
		dao.DB.Delete(&models.IptvChannel{}, delIDs)
	}

	go BindChannel()
	go until.UpdateChannelsId()
	log.Println("重复:", repetNum)
}

func BindChannel() {
	dao.DB.Model(&models.IptvEpg{}).Where("content != ''").Update("content", "")
	var channeList []models.IptvChannel
	if err := dao.DB.Model(&models.IptvChannel{}).Select("distinct name").Order("category,id").Find(&channeList).Error; err != nil {
		return
	}

	var epgList []models.IptvEpg
	if err := dao.DB.Model(&models.IptvEpg{}).Find(&epgList).Error; err != nil {
		return
	}

	for _, epgData := range epgList {
		var tmpList []string
		for _, channelData := range channeList {

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
	go until.GetCCTVChannelList(true)
	go until.GetProvinceChannelList(true)
}
