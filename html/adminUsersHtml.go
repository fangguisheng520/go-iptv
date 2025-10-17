package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Users(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	var pageData = dto.AdminUserDto{
		LoginUser: username,
		Title:     "用户管理",
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

	pageData.Order = c.DefaultQuery("order", "id")
	pageData.Keywords = c.DefaultQuery("keywords", "")

	if !until.IsSafe(pageData.Order) || !until.IsSafe(pageData.Keywords) {
		pageData.Order = "id"
		pageData.Keywords = ""
	}

	today := time.Now().Truncate(24 * time.Hour).Unix()

	// dao.DB.Model(&models.IptvUserShow{}).Count(&pageData.UserTotal)
	dao.DB.Model(&models.IptvUserShow{}).Where("lasttime > ?", today).Count(&pageData.UserToday)

	recStart := recCounts * (pageData.Page - 1)
	keywords := "%" + pageData.Keywords + "%"

	// 基础查询
	dbQuery := dao.DB.Table(models.IptvUserShow{}.TableName()+" u").Select(`u.*, m.name AS mealname`).
		Joins("LEFT JOIN iptv_meals m ON u.meal = m.id").Where("u.status > ?", 0)

	// 如果不是 admin，增加 author 条件
	if username != "admin" {
		dbQuery = dbQuery.Where("author = ?", username)
	}

	// 增加搜索条件
	dbQuery = dbQuery.Where(
		"u.name LIKE ? OR u.deviceid LIKE ? OR u.mac LIKE ? OR u.model LIKE ? OR u.ip LIKE ? OR u.region LIKE ? OR u.author LIKE ? OR u.marks LIKE ? OR CAST(u.status AS CHAR) LIKE ?",
		keywords, keywords, keywords, keywords, keywords, keywords, keywords, keywords, keywords,
	)

	err = dbQuery.Count(&pageData.UserTotal).Error
	if err != nil {
		pageData.PageCount = 1
	} else {
		if pageData.UserTotal == 0 {
			pageData.PageCount = 1
		} else {
			pageData.PageCount = int64(math.Ceil(float64(pageData.UserTotal) / float64(recCounts)))
		}
	}

	err = dbQuery.Offset(int(recStart)).Limit(int(recCounts)).Order("u." + pageData.Order).Find(&pageData.Users).Error
	if err != nil {
		log.Println("查询用户失败:", err)
	}
	pageData.Users = until.CheckUserDay(pageData.Users)

	dao.DB.Model(&models.IptvMeals{}).Find(&pageData.Meals)

	c.HTML(200, "admin_user.html", pageData)
}
