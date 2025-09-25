package api

import (
	"encoding/json"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ApkLogin(c *gin.Context) {
	var user dto.ApkUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	ip := c.ClientIP() //获取ip

	if !service.CheckIpMax(ip) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if strings.Contains(user.Mac, "获取地址失败") {
		user.Mac = user.DeviceID
	}

	dbUser := service.CheckUserInfo(user, ip)

	result := service.ApkLogin(dbUser)

	resObj, _ := json.Marshal(result)

	aes := until.NewAes(until.GetAesKey()[5:21], "AES-128-ECB", "")
	reAes, _ := aes.Encrypt(string(resObj))

	c.String(http.StatusOK, reAes)
}

func Getver(c *gin.Context) {
	result := service.Getver()
	c.JSON(http.StatusOK, result)
}

func GetBg(c *gin.Context) {
	imgName := service.GetBg()
	if imgName == "" {
		c.String(http.StatusOK, "")
		return
	}
	c.String(http.StatusOK, dao.GetConfig().ServerUrl+"/images/bj/"+imgName)
}

func GetChannels(c *gin.Context) {

	var channel dto.DataReqDto
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if strings.Contains(channel.Mac, "获取地址失败") {
		channel.Mac = channel.DeviceID
	}

	result := service.GetChannels(channel)

	c.String(http.StatusOK, result)
}

func GetWeather(c *gin.Context) {
	result := service.GetWeather()
	c.JSON(http.StatusOK, result)
}

func GetEpg(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	simple := c.DefaultQuery("simple", "")

	if simple != "1" {
		result := service.GetEpg(id)
		c.JSON(http.StatusOK, result)
	} else {
		result := service.GetSimpleEpg(id)
		c.JSON(http.StatusOK, result)
	}
}
