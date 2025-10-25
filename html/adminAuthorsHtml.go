package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Authors(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	var pageData = dto.AdminAuthorsDto{
		LoginUser: username,
		Title:     "用户授权",
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

	dao.DB.Model(&models.IptvUser{}).Where("status <= 0 and lasttime > ?", today).Count(&pageData.NewUserToday)
	dao.DB.Model(&models.IptvUser{}).Where("status>0 and authortime > ?", today).Count(&pageData.UserTodayAuthor)

	recStart := recCounts * (pageData.Page - 1)
	keywords := "%" + pageData.Keywords + "%"

	// 基础查询
	dbQuery := dao.DB.Model(&models.IptvUser{}).
		Select(`name,deviceid,model,ip,region,lasttime,exp,status`).
		Where("status <= ?", 0)

	dbQuery = dbQuery.Where(
		"name LIKE ? OR deviceid LIKE ? OR model LIKE ? OR ip LIKE ? OR region LIKE ? OR status LIKE ?",
		keywords, keywords, keywords, keywords, keywords, keywords,
	)

	err = dbQuery.Count(&pageData.UnAuthorUserTotal).Error
	if err != nil {
		pageData.PageCount = 1
	} else {
		if pageData.UnAuthorUserTotal == 0 {
			pageData.PageCount = 1
		} else {
			pageData.PageCount = int64(math.Ceil(float64(pageData.UnAuthorUserTotal) / float64(recCounts)))
		}
	}

	dbQuery.Offset(int(recStart)).Limit(int(recCounts)).Order(pageData.Order).Find(&pageData.Users)
	pageData.Users = until.CheckUserDay(pageData.Users)

	dao.DB.Model(&models.IptvMeals{}).Where("status = 1").Find(&pageData.Meals)

	c.HTML(http.StatusOK, "admin_authors.html", pageData)
}
