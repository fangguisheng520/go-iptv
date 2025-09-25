package until

import (
	"crypto/md5"
	"encoding/hex"
)

var PANEL_MD5_KEY = "tvkey_"

// HashPassword 使用 PANEL_MD5_KEY + 密码 做 md5
func HashPassword(password string) string {
	h := md5.New()
	h.Write([]byte(PANEL_MD5_KEY + password))
	return hex.EncodeToString(h.Sum(nil))
}
