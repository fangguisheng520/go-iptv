package crontab

import (
	"go-iptv/until"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func EpgCron() {
	c := cron.New(cron.WithSeconds()) // 支持秒级别 cron 表达式

	// 每天凌晨1点执行
	// cron 表达式格式: 秒 分 时 日 月 星期
	// 下面表示每天 01:00:00
	c.AddFunc("0 0 1 * * *", func() {
		log.Println("更新EPG列表任务开始执行:", time.Now().Format("2006-01-02 15:04:05"))
		// 在这里写你的任务逻辑
		if until.UpdataEpgList() {
			log.Println("更新EPG列表任务执行成功:", time.Now().Format("2006-01-02 15:04:05"))
		} else {
			log.Println("更新EPG列表任务执行失败:", time.Now().Format("2006-01-02 15:04:05"))
		}
	})
	c.Start()
}
