package bootstrap

import (
	"go-iptv/until"
	"time"
)

func InitJwtKey() {
	// 读取配置文件
	hostname, _ := until.GetContainerID()
	until.JwtKey = []byte(until.Md5(hostname + time.Now().Format("2006-01-02 15:04:05")))
	until.RssKey = []byte(until.Md5(hostname))
}
