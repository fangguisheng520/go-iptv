package bootstrap

import (
	"go-iptv/dao"
	"go-iptv/until"
	"log"
	"os"
	"os/exec"
)

var Installed bool = false

func Install() (bool, string) {

	if !until.Exists("/config") {
		log.Println("请映射config文件夹到容器/config中")
		return false, "请映射config文件夹到容器/config中"
	}

	if !until.Exists("/app/database/sqlite.sql") || !until.Exists("/app/config.yml") {
		log.Println("缺少必要的文件")
		return false, "缺少必要的文件"
	}

	os.RemoveAll("/config")
	if err := os.MkdirAll("/config", 0755); err != nil {
		log.Println("创建/config文件夹失败:", err)
		return false, err.Error()
	}

	if err := until.CopyFile("/app/config.yml", "/config/config.yml"); err != nil {
		log.Println("复制配置文件失败:", err)
		return false, "复制配置文件失败:" + err.Error()
	}

	cmd := exec.Command("sqlite3", "/config/iptv.db")
	cmd.Stdin, _ = os.Open("/app/database/sqlite.sql") // 把 SQL 文件内容传给标准输入

	if err := cmd.Run(); err != nil {
		log.Println("初始化数据库失败:", err)
		return false, "初始化数据库失败:" + err.Error()
	}
	log.Println("初始化数据库完成")
	log.Println("加载数据库...")
	dao.InitDB("/config/iptv.db")
	log.Println("初始化EPG缓存...")
	cache, err := dao.NewFileCache("/config/cache/", true)
	if err != nil {
		log.Println("初始化缓存失败:", err)
		return false, "初始化缓存失败:" + err.Error()
	}
	dao.Cache = cache

	if !InitLogo() {
		log.Println("初始化Logo失败")
		return false, "初始化Logo失败"
	}

	dao.CONFIG_PATH = "/config/config.yml"
	dao.LoadConfigFile()

	if !dao.LoadConfig() {
		log.Println("conf加载错误")
		return false, "conf加载错误"
	}
	file, err := os.Create("/config/install.lock") // 创建文件
	if err != nil {
		log.Println("创建install.lock失败:", err)
		return false, "创建install.lock失败:" + err.Error()
	}
	defer file.Close()
	InitJwtKey()
	return true, "success"
}
