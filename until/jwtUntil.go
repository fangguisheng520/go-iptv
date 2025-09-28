package until

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("test")

// ----------------------
// 生成 JWT
// ----------------------
func GenerateJWT(username string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ----------------------
// 验证 JWT
// ----------------------
func VerifyJWT(tokenString string) (jwt.MapClaims, error, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp := int64(claims["exp"].(float64))
		now := time.Now().Unix()
		if exp-now < 600 {
			return claims, nil, true
		}
		return claims, nil, false
	}

	return nil, jwt.ErrTokenMalformed, false
}

func GetAuthName(c *gin.Context) (string, bool) {
	auth, exists := c.Get("auth")
	if !exists {
		return "", false
	}
	claims, ok := auth.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return "", false
	}
	return claims["username"].(string), true
}

func GetAuthExp(claims jwt.MapClaims) int64 {
	return int64(claims["exp"].(float64))
}
