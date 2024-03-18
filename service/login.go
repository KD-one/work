package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginService struct {
	Name      string `gorm:"varchar(20);not null" json:"name" form:"name"`
	Password  string `gorm:"size:255;not null" json:"password" form:"password"`
	UserLevel int    `gorm:"not null" json:"user_level" form:"user_level"`
	AppAuth   string `gorm:"varchar(255);not null" json:"app_auth" form:"app_auth"`
	ParaAuth  string `gorm:"varchar(255);not null" json:"para_auth" form:"para_auth"`
}

// ToLogin 展示用户登录前端页面模板
func ToLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "loginRegisterView/login.html", nil)
}

//c.SetCookie("auth", v, 84600, "/", "127.0.0.1", false, true)
//fmt.Println("---------成功设置cookie----------------")

// Info 返回用户信息   （需要身份认证通过后才能访问）
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	//将用户信息返回
	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{"user": user},
	})
}
