package api

import (
	"fmt"
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRssUrl(c *gin.Context) {
	_, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	c.Request.ParseForm()
	params := c.Request.PostForm
	id := params.Get("id")

	c.JSON(200, service.GetRssUrl(id))
}

func GetTXTRss(c *gin.Context) {
	if token, ok := c.GetQuery("token"); !ok {
		c.String(200, "token 参数不存在")
		return
	} else {
		if token == "" {
			c.String(200, "token 参数不存在")
			return
		}
		scheme := GetClientScheme(c)

		host := c.Request.Host
		if !until.IsValidHost(host) {
			c.String(200, "参数不合法")
			return
		}
		host = fmt.Sprintf("%s://%s", scheme, host)

		c.String(200, service.GetRss(token, host))
	}
}

func GetTXTRssEpg(c *gin.Context) {
	if token, ok := c.GetQuery("token"); !ok {
		c.String(200, "token 参数不存在")
		return
	} else {
		if token == "" {
			c.String(200, "token 参数不存在")
			return
		}
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		host := c.Request.Host
		if !until.IsValidHost(host) {
			c.String(200, "参数不合法")
			return
		}
		host = fmt.Sprintf("%s://%s", scheme, host)

		c.XML(200, service.GetRssEpg(token, host))
	}
}

func GetClientScheme(c *gin.Context) string {
	// 1) X-Forwarded-Proto（可能是 "https" 或 "http"，也可能是逗号分隔的列表）
	if xf := c.Request.Header.Get("X-Forwarded-Proto"); xf != "" {
		// 取第一个值，移除空格，小写
		parts := strings.Split(xf, ",")
		if len(parts) > 0 {
			return strings.ToLower(strings.TrimSpace(parts[0]))
		}
	}
	if xf := c.Request.Header.Get("X-Forwarded-Scheme"); xf != "" {
		// 取第一个值，移除空格，小写
		parts := strings.Split(xf, ",")
		if len(parts) > 0 {
			return strings.ToLower(strings.TrimSpace(parts[0]))
		}
	}

	// 2) Forwarded: 表示形式如: Forwarded: for=192.0.2.60;proto=https;by=203.0.113.43
	if f := c.Request.Header.Get("Forwarded"); f != "" {
		// 简单查找 proto= 后面的值（更严格的解析可用正则或更完整解析）
		// 例如 "for=..., proto=https; ..." 或 ";proto=https"
		if i := strings.Index(strings.ToLower(f), "proto="); i != -1 {
			// 从 proto= 后面截取到下一个分号或逗号或结尾
			v := f[i+len("proto="):]
			end := len(v)
			for j, ch := range v {
				if ch == ';' || ch == ',' {
					end = j
					break
				}
			}
			return strings.ToLower(strings.TrimSpace(v[:end]))
		}
	}

	// 3) X-Forwarded-SSL: on 表示 https（一些旧代理会设置）
	if xfs := strings.ToLower(c.Request.Header.Get("X-Forwarded-SSL")); xfs == "on" {
		return "https"
	}

	// 4) 回退：检查当前连接是否使用 TLS（适用于没有代理或直连）
	if c.Request.TLS != nil {
		return "https"
	}
	return "http"
}
