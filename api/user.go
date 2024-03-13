package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"test/service"
)

func GetUserList(c *gin.Context) {
	var data service.RequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	res := data.GetUserList(c)
	c.JSON(200, res)
}

func Login(c *gin.Context) {
	var data service.LoginService
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	res := data.Login(c)
	c.JSON(200, res)
}

func AddChangeUser(c *gin.Context) {
	var data service.RequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	fmt.Println("data: ", data)
	res := data.AddChangeUser(c)
	c.JSON(200, res)
}
