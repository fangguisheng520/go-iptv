package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Epgs(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	var pageData = dto.AdminEpgsDto{
		LoginUser: username,
		Title:     "EPG管理",
	}

	recCountsStr := c.DefaultQuery("recCounts", "20")
	jumptoStr := c.DefaultQuery("jumpto", "")
	pageStr := c.DefaultQuery("page", "")

	if !until.IsSafe(recCountsStr) || !until.IsSafe(jumptoStr) || !until.IsSafe(pageStr) {
		recCountsStr = "20"
		jumptoStr = ""
		pageStr = ""
	}

	recCounts, err := strconv.ParseInt(recCountsStr, 10, 64)
	if err != nil {
		// 转换失败时设置默认值，比如 20
		recCounts = 20
	}
	pageData.RecCounts = recCounts

	if jumptoStr != "" {
		pageData.Page, err = strconv.ParseInt(jumptoStr, 10, 64)
		if err != nil {
			pageData.Page = 1
		}
	} else if pageStr != "" {
		pageData.Page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			pageData.Page = 1
		}
	} else {
		pageData.Page = 1
	}

	pageData.Keywords = c.DefaultQuery("keywords", "") // 关键词

	if !until.IsSafe(pageData.Keywords) {
		pageData.Keywords = ""
	}

	recStart := recCounts * (pageData.Page - 1)
	keywords := "%" + pageData.Keywords + "%" // 模糊查询

	dbQuery := dao.DB.Model(&models.IptvEpg{}).Where("name like ? or remarks like ? or content like ?", keywords, keywords, keywords)

	var count int64
	err = dbQuery.Count(&count).Error
	if err != nil {
		pageData.PageCount = 1
	} else {
		if count == 0 {
			pageData.PageCount = 1
		} else {
			pageData.PageCount = int64(math.Ceil(float64(count) / float64(recCounts)))
		}
	}

	err = dbQuery.Offset(int(recStart)).Limit(int(recCounts)).Find(&pageData.Epgs).Error
	if err != nil {
		log.Println("查询epg失败:", err)
	}

	logoList := until.GetLogs() // 获取logo列表

	cfg := dao.GetConfig()

	for k, v := range pageData.Epgs {
		epgName := strings.Split(v.Name, "-")[1]
		log.Println("epgName", epgName)
		for _, logo := range logoList {
			logoName := strings.Split(logo, ".")[0]
			log.Println("logoName", logoName)
			if strings.EqualFold(epgName, logoName) {
				pageData.Epgs[k].Logo = cfg.ServerUrl + "/logo/" + logo
			}
		}
	}

	// cfg := dao.GetConfig()

	// pageData.EpgErr = cfg.EPGErrors
	// pageData.EPGApiChk = cfg.App.EPGApiChk

	c.HTML(200, "admin_epgs.html", pageData)
}
