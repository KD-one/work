package api

import (
	"github.com/gin-gonic/gin"
	"test/service"
)

func Update(c *gin.Context) {
	var data service.UpdateService
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, gin.H{
			"msg": "参数错误",
		})
		return
	}
	res := data.Update(c)
	c.JSON(200, res)
}

func ToUpdate(c *gin.Context) {
	c.HTML(200, "update/update.html", nil)
}
