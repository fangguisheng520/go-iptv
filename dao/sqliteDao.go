package dao

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		// Logger: logger.New(
		// 	log.New(log.Writer(), "\r\n", log.LstdFlags), // 使用标准日志输出
		// 	logger.Config{
		// 		SlowThreshold:             time.Millisecond, // 慢 SQL 阈值
		// 		LogLevel:                  logger.Info,      // LogLevel 可选: Silent, Error, Warn, Info
		// 		IgnoreRecordNotFoundError: true,             // 忽略记录未找到错误
		// 		Colorful:                  true,             // 彩色输出
		// 	},
		// ),
	})
	if err != nil {
		log.Fatal("无法连接数据库: ", err)
	}
	// Migrate the schema
	// DB.AutoMigrate(&User{})
}

func InitDBDebug(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags), // 使用标准日志输出
			logger.Config{
				SlowThreshold:             time.Millisecond, // 慢 SQL 阈值
				LogLevel:                  logger.Info,      // LogLevel 可选: Silent, Error, Warn, Info
				IgnoreRecordNotFoundError: true,             // 忽略记录未找到错误
				Colorful:                  true,             // 彩色输出
			},
		),
	})
	if err != nil {
		log.Fatal("无法连接数据库: ", err)
	}
	// Migrate the schema
	// DB.AutoMigrate(&User{})
}
