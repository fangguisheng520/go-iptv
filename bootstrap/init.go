package bootstrap

import (
	"go-iptv/dao"
	"go-iptv/models"
	"go-iptv/until"
	"log"
	"os"
	"os/exec"
)

func InitDB() {
	dao.DB.AutoMigrate(&models.IptvAdmin{})
	dao.DB.AutoMigrate(&models.IptvUser{})
	dao.DB.AutoMigrate(&models.IptvChannel{})
	dao.DB.AutoMigrate(&models.IptvCategory{})
	dao.DB.AutoMigrate(&models.IptvEpg{})
	dao.DB.AutoMigrate(&models.IptvMeals{})
	dao.DB.AutoMigrate(&models.IptvMovie{})
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
