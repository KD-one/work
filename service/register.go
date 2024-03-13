package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"test/common"
	"test/dao"
	"test/model"
)

func ToRegister(c *gin.Context) {
	c.HTML(200, "loginRegisterView/register.html", nil)
}

func Register(c *gin.Context) {
	var user model.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.String(500, "绑定数据出错：%v", err)
	}

	name := user.Name
	password := user.Password
	appAuth := user.AppAuth
	paraAuth := user.ParaAuth

	if len(name) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户名不能为空！",
		})
		return
	}

	if len(appAuth) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "应用权限不能为空！",
		})
		return
	}

	if len(paraAuth) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "参数权限不能为空！",
		})
		return
	}

	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码长度不能小于6位！",
		})
		return
	}

	err = dao.DBUserRegister(name, password, appAuth, paraAuth, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": err.Error(),
		})
		return
	}

	//err = common.DB.AutoMigrate(&model.User{})
	//if err != nil {
	//	c.JSON(500, gin.H{
	//		"message": "注册时添加字段失败！：" + err.Error(),
	//	})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
	})
	fmt.Printf("成功注册！\n用户名: %s \n密码: %s \n", name, password)
	common.UserRecord.Printf(" [info] 注册成功 用户名: %s 密码: %s \n", name, password)
}
