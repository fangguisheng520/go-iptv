package router

import (
	"go-iptv/api"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func ApkRouter(r *gin.Engine, path string) {

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	router := r.Group(path)
	{
		router.GET("/weather", api.GetWeather)
		router.GET("/getepg", api.GetEpg)
		router.GET("/getver", api.Getver)
		router.GET("/bg", api.GetBg)
		router.POST("/login", api.ApkLogin)
		router.POST("/channels", api.GetChannels)

	}
}
