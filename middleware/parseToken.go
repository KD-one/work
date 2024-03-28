package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test/common"
)

// ParseToken 该中间件用于判断token是否有效，并将有效的token解析
func ParseToken() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// validate token formate
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足，请先去登录后再访问！"})
			ctx.Abort()
			return
		}

		// 解析token
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足，请先去登录后再访问！"})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.UserId)

		ctx.Next()
	}
}
