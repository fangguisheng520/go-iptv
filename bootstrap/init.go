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

type IptvCategory struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"unique;column:name" json:"name"`
	Enable       int64  `gorm:"column:enable;default:1" json:"enable"`
	Type         string `gorm:"default:hand;column:type" json:"type"`
	Url          string `gorm:"column:url" json:"url"`
	UA           string `gorm:"column:ua" json:"ua"`
	LatestTime   string `gorm:"column:latesttime" json:"latesttime"`
	AutoCategory int64  `gorm:"column:autocategory" json:"autocategory"`
	Repeat       int64  `gorm:"column:repeat" json:"repeat"`
	Sort         int64  `gorm:"column:sort" json:"sort"`
	Rawcount     int64  `gorm:"column:rawcount;default:0" json:"rawcount"`
}

func (IptvCategory) TableName() string {
	return "iptv_category"
}

type IptvChannel struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Url      string `gorm:"column:url" json:"url"`
	Category string `gorm:"column:category" json:"category"`
}

func (IptvChannel) TableName() string {
	return "iptv_channels"
}

func InitDB() bool {
	dao.DB.AutoMigrate(&models.IptvAdmin{})
	dao.DB.AutoMigrate(&models.IptvUser{})

	has := dao.DB.Migrator().HasColumn(&IptvChannel{}, "sort")
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

	has = dao.DB.Migrator().HasColumn(&models.IptvChannel{}, "category")
	if has {
		dao.DB.Exec("ALTER TABLE iptv_channels DROP COLUMN category;")
	}

	dao.DB.AutoMigrate(&models.IptvChannel{})
	dao.DB.AutoMigrate(&models.IptvCategoryList{})

	has = dao.DB.Migrator().HasColumn(&IptvCategory{}, "url")
	if has {

		var categories []IptvCategory
		dao.DB.Model(&IptvCategory{}).Where("url != ?", "").Find(&categories)
		var list []models.IptvCategoryList
		for _, category := range categories {
			list = append(list, models.IptvCategoryList{
				Name:         category.Name,
				Url:          category.Url,
				Enable:       category.Enable,
				AutoCategory: category.AutoCategory,
				Repeat:       category.Repeat,
				UA:           category.UA,
			})
		}
		if len(list) > 0 {
			dao.DB.Create(&list)
		}
	}

	has = dao.DB.Migrator().HasColumn(&IptvCategory{}, "latesttime")
	if has {
		dao.DB.Exec("ALTER TABLE iptv_category DROP COLUMN url DROP COLUMN latesttime DROP COLUMN autocategory DROP COLUMN repeat;")
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
