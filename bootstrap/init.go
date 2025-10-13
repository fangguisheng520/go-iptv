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
		err1 := os.RemoveAll("/config/logo") // 删除文件夹
		if err1 != nil {
			log.Println("删除logo失败:", err1)
			return false
		}
		err2 := os.MkdirAll("/config/logo", os.ModePerm) // 创建文件夹
		if err2 != nil {
			log.Println("创建logo失败:", err2)
			return false
		}
		cmd := exec.Command("bash", "-c", "cp -rf ./logo/* /config/logo")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("复制logo失败: %v --- %s\n", err, string(output))
			return false
		}
		// if err := cmd.Run(); err != nil {
		// 	log.Println("复制logo失败:", err)
		// 	return false
		// }
	}
	return true
}
