package api

import (
	"go-iptv/dao"
	"go-iptv/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func EditUsers(c *gin.Context) {
	c.Request.ParseForm()
	params := c.Request.PostForm
	for k := range params {
		switch k {
		case "submitdelall":
			nowTime := time.Now().Unix()
			dao.DB.Model(&models.IptvUser{}).Where("status=1 and exp < ?", nowTime).Delete(&models.IptvUser{})
			c.JSON(200, gin.H{"code": 1, "msg": "已清空所有过期用户", "type": "success"})
			return
		case "submitdel":
			ids := c.PostFormArray("ids[]")
			if len(ids) == 0 {
				c.JSON(200, gin.H{"code": 0, "msg": "请选择要删除的用户账号", "type": "danger"})
				return
			}
			dao.DB.Where("name in (?)", ids).Delete(&models.IptvUser{})
			c.JSON(200, gin.H{"code": 1, "msg": "已删除选中的用户账号", "type": "success"})
			return
		case "submitmodify":
			ids := c.PostFormArray("ids[]")
			if len(ids) == 0 {
				c.JSON(200, gin.H{"code": 0, "msg": "请选择要修改授权天数的用户账号", "type": "danger"})
				return
			}
			expStr := c.PostForm("exp")
			expDays, err := strconv.Atoi(expStr)
			if err != nil {
				c.JSON(200, gin.H{"code": 0, "msg": "授权天数格式不正确", "type": "danger"})
				return
			}
			now := time.Now()
			todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			targetTime := todayStart.Add(time.Duration(expDays) * 24 * time.Hour)
			dao.DB.Model(&models.IptvUser{}).Where("name in (?)", ids).Updates(map[string]interface{}{
				"exp":    targetTime.Unix(),
				"status": 1,
			})
			c.JSON(200, gin.H{"code": 1, "msg": "已修改选中的用户账号的授权天数", "type": "success"})
			return
		case "submitmodifymarks":
			ids := c.PostFormArray("ids[]")
			if len(ids) == 0 {
				c.JSON(200, gin.H{"code": 0, "msg": "请选择要修改备注的用户账号", "type": "danger"})
				return
			}
			marks := c.PostForm("marks")
			dao.DB.Model(&models.IptvUser{}).Where("name in (?)", ids).Update("marks", marks)
			c.JSON(200, gin.H{"code": 1, "msg": "已修改选中的用户账号的备注", "type": "success"})
			return
		case "submitNotExpired":
			ids := c.PostFormArray("ids[]")
			if len(ids) == 0 {
				c.JSON(200, gin.H{"code": 0, "msg": "请选择要设为永不到期的用户账号", "type": "danger"})
				return
			}
			dao.DB.Model(&models.IptvUser{}).Where("name in (?)", ids).Updates(map[string]interface{}{
				"exp":    0,
				"status": 999,
			})
			c.JSON(200, gin.H{"code": 1, "msg": "已设为选中的用户账号为永不到期", "type": "success"})
			return
		case "submitCancelNotExpired":
			ids := c.PostFormArray("ids[]")
			if len(ids) == 0 {
				c.JSON(200, gin.H{"code": 0, "msg": "请选择要取消永不到期的用户账号", "type": "danger"})
				return
			}
			now := time.Now()
			todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			targetTime := todayStart.Add(24 * time.Hour)
			dao.DB.Model(&models.IptvUser{}).Where("name in (?)", ids).Updates(map[string]interface{}{
				"exp":    targetTime.Unix(),
				"status": 1,
			})
			c.JSON(200, gin.H{"code": 1, "msg": "已取消选中的用户账号的永不到期", "type": "success"})
			return
		case "submitforbidden":
			ids := c.PostFormArray("ids[]")
			if len(ids) == 0 {
				c.JSON(200, gin.H{"code": 0, "msg": "请选择要取消授权的用户账号", "type": "danger"})
				return
			}
			dao.DB.Model(&models.IptvUser{}).Where("name in (?)", ids).Updates(map[string]interface{}{
				"status": 0,
			})
			c.JSON(200, gin.H{"code": 1, "msg": "已取消选中的用户账号的授权", "type": "success"})
			return
		case "e_meals":
			ids := c.PostFormArray("ids[]")
			if len(ids) == 0 {
				c.JSON(200, gin.H{"code": 0, "msg": "请选择要修改套餐的用户账号", "type": "danger"})
				return
			}
			mealStr := c.PostForm("s_meals")
			mealID, err := strconv.Atoi(mealStr)
			if err != nil {
				c.JSON(200, gin.H{"code": 0, "msg": "套餐格式不正确", "type": "danger"})
				return
			}
			var meal models.IptvMeals
			err = dao.DB.Where("id = ?", mealID).First(&meal).Error
			if err != nil {
				c.JSON(200, gin.H{"code": 0, "msg": "选择的套餐不存在", "type": "danger"})
				return
			}

			dao.DB.Model(&models.IptvUser{}).Where("name in (?)", ids).Updates(map[string]interface{}{
				"meal":   meal.ID,
				"status": 1,
			})
			c.JSON(200, gin.H{"code": 1, "msg": "已修改选中的用户账号的套餐", "type": "success"})
			return
		}
	}
}
