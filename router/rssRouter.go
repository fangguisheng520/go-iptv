package router

import (
	"go-iptv/api"

	"github.com/gin-gonic/gin"
)

func RssRouter(r *gin.Engine, path string) {
	router := r.Group(path)
	{
		router.GET("/getRss", api.GetTXTRss)
	}
}
