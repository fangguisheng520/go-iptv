package html

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.HTML(200, "admin_login.html", gin.H{
		"CurrentURL": strings.TrimSuffix(c.Request.URL.String(), "/"), // 当前URL
	})
}
