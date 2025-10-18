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
	dao.DB.AutoMigrate(&models.IptvEpgList{})
	var epgList []models.IptvEpgList
	if err := dao.DB.Model(&models.IptvEpgList{}).Find(&epgList).Error; err != nil {
		return false
	}
	if len(epgList) == 0 {
		dao.DB.Where("name like ?", "51zmt-%").Delete(&models.IptvEpg{})
		var update = models.IptvEpgList{
			Name:    "51zmt",
			Remarks: "51zmt",
			Url:     "http://epg.51zmt.top:8000/e.xml",
			Status:  1,
		}
		dao.DB.Model(&models.IptvEpgList{}).Save(&update)
		if !until.UpdataEpgListOne(update.ID) {
			log.Println("初始化51zmt失败")
		}
	}
	dao.DB.AutoMigrate(&models.IptvMeals{})
	dao.DB.AutoMigrate(&models.IptvMovie{})
	
	// 确保 rawcount 字段存在（类型为 integer，默认值为 0）
	if !dao.DB.Migrator().HasColumn(&models.IptvCategory{}, "rawcount") {
		log.Printf("rawcount 字段不存在,开始添加")
		// 添加字段
		if err := dao.DB.Migrator().AddColumn(&models.IptvCategory{}, "rawcount"); err != nil {
			log.Printf("添加 rawcount 字段失败: %v", err)
			return false
		}
	}
	
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
