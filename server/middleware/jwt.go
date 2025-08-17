package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"net/http"
)

var Secret []byte

func InitSecret() {
	Secret = []byte(viper.GetString("app.Jwt") + "mundo")
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	// StandardClaims 已经弃用，使用 RegisteredClaims
	jwt.RegisteredClaims
}

// JWTAuthMiddleware 是一个Gin中间件，用于验证JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": http.StatusUnauthorized,
					"msg":  "请求未携带token，无权限访问",
				})
				c.Abort()
				return
			}
		} else if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "无效的token形式",
			})
			c.Abort()
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中，以便后续处理使用
		c.Set("uid", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// ParseToken 验证 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})

	if err != nil {
		// 区分具体错误类型
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token已过期")
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, errors.New("token签名无效")
			} else {
				return nil, errors.New("token解析错误: " + err.Error())
			}
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
