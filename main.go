package main

import (
	"flag"
	"go-iptv/bootstrap"
	"go-iptv/crontab"
	"go-iptv/dao"
	"go-iptv/router"
	"go-iptv/until"
	"log"
)

func main() {
	port := flag.String("port", "80", "启动端口 eg: 80")
	flag.Parse()
	if !until.CheckPort(*port) {
		return
	}

	if !until.CheckJava() {
		log.Println("请安装Java JDK 1.8环境")
		return
	}

	if !until.CheckApktool() {
		log.Println("请安装apktool环境")
		return
	}

	log.Println("初始化EPG缓存...")
	cache, err := dao.NewFileCache("/config/cache/", true)
	if err != nil {
		log.Println("初始化缓存失败:", err)
		return
	}
	dao.Cache = cache

	bootstrap.InitJwtKey() // 初始化JWTkey

	if !until.Exists("/config/iptv.db") || !until.Exists("/config/config.yml") || !until.Exists("/config/install.lock") {
		bootstrap.Installed = false
		log.Println("检测到未安装，开始安装...")
		log.Println("启动接口...")
		router := router.InitRouter()
		router.Run(":" + *port)
	} else {
		bootstrap.Installed = true
	}

	log.Println("加载数据库...")
	dao.InitDB("/config/iptv.db")
	if !bootstrap.InitDB() {
		log.Println("数据库初始化失败,请删除/config/iptv.db重新安装")
	}

	dao.CONFIG_PATH = "/config/config.yml"
	dao.LoadConfigFile()

	if !dao.LoadConfig() {
		log.Println("conf加载错误")
		return
	}

	if !bootstrap.InitLogo() {
		log.Println("logo目录初始化错误")
		return
	}

	go crontab.Crontab()

	if bootstrap.Installed {
		if !bootstrap.BuildAPK() {
			log.Println("APK编译错误")
			return
		}
	}

	log.Println("启动接口...")
	router := router.InitRouter()
	router.Run(":" + *port)
}
