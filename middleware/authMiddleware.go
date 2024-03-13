package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"test/common"
	"test/model"
)

// AuthMiddleware 该中间件用于判断token是否有效，并将有效的token解析
func AuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足，请先去登录后再访问！"})
			ctx.Abort()
			return
		}

		//提取token的有效部分（"Bearer "共占7位)
		tokenString = tokenString[7:]

		// 解析token
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足，请先去登录后再访问！"})
			ctx.Abort()
			return
		}

		// 验证通过后获取claim 中的userId
		userId := claims.UserId
		var user model.User
		common.DB.First(&user, userId)

		// 用户不存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足，请先去登录后再访问！"})
			ctx.Abort()
			return
		}

		// 将用户名写入上下文中
		ctx.Set("userName", user.Name)
		ctx.Set("userId", user.ID)

		ctx.Next()
	}
}

//// Auth cookie权限验证中间件
//func Auth(c *gin.Context) {
//
//	if uname := getUsernameFromCookie(c); uname == "" {
//		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录，无法访问！"})
//		// 跳转路由
//		//c.Redirect(401, "127.0.0.1:9999/toLogin")
//		c.Abort()
//		return
//	}
//	c.Next()
//}
//
//// 从cookie中获取用户名
//func getUsernameFromCookie(c *gin.Context) string {
//	for _, cookie := range strings.Split(c.Request.Header.Get("cookie"), ";") {
//		arr := strings.Split(cookie, "=")
//		key := strings.TrimSpace(arr[0])
//		value := strings.TrimSpace(arr[1])
//		if key == "auth" {
//			if uname, ok := LoggedIn[value]; ok {
//				return uname
//			}
//		}
//	}
//	return ""
//}
