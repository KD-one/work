package middleware

import (
	"github.com/gin-gonic/gin"
	"test/common"
)

// casbin中间件
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.RequestURI() //
		method := c.Request.Method
		common.UserRecord.Println(path, method)
		//验证url权限,只有root用户可以通过验证
		roleId := "root"
		ok, _ := common.CasBin.Enforce(roleId, path, method)
		if ok {
			c.Next()
		} else {
			c.Abort()
			c.JSON(200, gin.H{
				"msg": "很遗憾,权限验证没有通过",
			})
			return
		}
	}
}
