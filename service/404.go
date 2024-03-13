package service

import "github.com/gin-gonic/gin"

func NotFound(c *gin.Context) {
	c.HTML(404, "errView/404.html", nil)
}
