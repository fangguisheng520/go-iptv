package bootstrap

import (
	"go-iptv/dao"
	"go-iptv/models"
	"go-iptv/until"
	"log"
	"os"
	"os/exec"

	"gorm.io/gorm"
)

func InitDB() bool {
	dao.DB.AutoMigrate(&models.IptvAdmin{})
	dao.DB.AutoMigrate(&models.IptvUser{})

	has := dao.DB.Migrator().HasColumn(&models.IptvChannel{}, "sort")
	if !has {
		dao.DB.AutoMigrate(&models.IptvChannel{})
		if err := dao.DB.Transaction(func(tx *gorm.DB) error {
			var channels []models.IptvChannel
			if err := tx.Model(&models.IptvChannel{}).Order("id").Find(&channels).Error; err != nil {
				return err
			}

			for _, ch := range channels {
				if err := tx.Model(&models.IptvChannel{}).Where("id = ?", ch.ID).Update("sort", ch.ID).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return false
		}
	}
	dao.DB.AutoMigrate(&models.IptvCategory{})
	dao.DB.AutoMigrate(&models.IptvEpg{})
	dao.DB.AutoMigrate(&models.IptvMeals{})
	dao.DB.AutoMigrate(&models.IptvMovie{})
	return true
}

func InitLogo() bool {
	is, err := until.CheckLogo("/config/logo")
	if err != nil || !is {
		os.RemoveAll("/config/logo")             // 删除文件夹
		os.MkdirAll("/config/logo", os.ModePerm) // 创建文件夹
		cmd := exec.Command("cp", "-r", "-f", "./logo", "/config/logo")
		if err := cmd.Run(); err != nil {
			log.Println("cp error:", err)
			return false
		}
	}
	return true
}
