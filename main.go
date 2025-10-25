package main

import (
	"flag"
	"go-iptv/bootstrap"
	"go-iptv/crontab"
	"go-iptv/dao"
	"go-iptv/router"
	"go-iptv/until"
	"log"
	"os"
)

func main() {
	if until.CheckRam() {
		log.Println("可用内存不足512MB，无法运行")
		return
	}

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

	var debug bool = false
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	log.Println("初始化EPG缓存...")
	cache, err := dao.NewFileCache("/config/cache/", true)
	if err != nil {
		log.Println("初始化缓存失败:", err)
		return
	}
	dao.Cache = cache
	if dao.Cache.Clear() != nil {
		log.Println("初始化清除缓存失败:", err)
		return
	}

	if !until.Exists("/config/iptv.db") || !until.Exists("/config/config.yml") || !until.Exists("/config/install.lock") {
		bootstrap.Installed = false
		log.Println("检测到未安装，开始安装...")
		log.Println("启动接口...")
		router := router.InitRouter(debug)
		router.Run(":" + *port)
	} else {
		bootstrap.Installed = true
	}

	log.Println("加载数据库...")
	if debug {
		debug = true
		dao.InitDBDebug("/config/iptv.db")
	} else {
		dao.InitDB("/config/iptv.db")
	}

	if !bootstrap.InitDB() {
		log.Println("数据库初始化失败,请删除/config/iptv.db重新安装")
		return
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
	go crontab.EpgCron()
	go until.InitCacheRebuild()

	if !debug {
		bootstrap.InitJwtKey() // 初始化JWTkey
		if !bootstrap.BuildAPK() {
			log.Println("APK编译错误")
			return
		}
	}

	log.Println("启动接口...")
	router := router.InitRouter(debug)
	router.Run(":" + *port)
}
