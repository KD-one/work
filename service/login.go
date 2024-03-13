package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"test/common"
	"test/dao"
	"test/serializer"
	"time"
)

// User 保存用户登录相关信息
type OnlineUserInfo struct {
	Username string
	Password string
	Token    string // 自定义生成的token
	Expires  time.Time
}

type LoginService struct {
	Name      string `gorm:"varchar(20);not null" json:"name" form:"name"`
	Password  string `gorm:"size:255;not null" json:"password" form:"password"`
	UserLevel int    `gorm:"not null" json:"user_level" form:"user_level"`
	AppAuth   string `gorm:"varchar(255);not null" json:"app_auth" form:"app_auth"`
	ParaAuth  string `gorm:"varchar(255);not null" json:"para_auth" form:"para_auth"`
}

// OnlineUserMap 用户在线状态存储
var OnlineUserMap sync.Map

// ToLogin 展示用户登录前端页面模板
func ToLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "loginRegisterView/login.html", nil)
}

func (l *LoginService) Login(c *gin.Context) serializer.Response {

	name := l.Name
	password := l.Password

	if len(name) == 0 {
		return serializer.Response{
			Code: 422,
			Msg:  "用户名不能为空",
		}
	}

	if len(password) < 6 {
		return serializer.Response{
			Code: 422,
			Msg:  "密码长度不能小于6",
		}
	}

	var userId uint
	err := dao.DBUserLogin(name, password, &userId)
	if err != nil {
		return serializer.Response{
			Code: 422,
			Msg:  err.Error(),
		}
	}

	//发放token
	token, err := common.ReleaseToken(userId)
	if err != nil {
		return serializer.Response{
			Code: 500,
			Msg:  "发放token失败",
		}
	}

	// 通过协程添加用户到在线用户中
	//go addOnlineUser(name, password, token)
	//
	//// 启动协程持续检查用户是否过期
	//go func() {
	//	for {
	//		cleanExpiredUsers()
	//		time.Sleep(time.Second * 20) // 每20秒检查一次
	//	}
	//}()
	common.UserRecord.Println(" [info] 登录成功", name)
	//返回结果
	return serializer.Response{
		Code: 200,
		Msg:  "登录成功",
		Data: gin.H{
			"token": token,
		},
	}

	//c.SetCookie("auth", v, 84600, "/", "127.0.0.1", false, true)
	//fmt.Println("---------成功设置cookie----------------")

}

// Info 返回用户信息   （需要身份认证通过后才能访问）
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	//将用户信息返回
	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{"user": user},
	})
}

// 添加用户到在线用户中
func addOnlineUser(name, password, token string) {
	u := &OnlineUserInfo{
		Username: name,
		Password: password,
		Expires:  time.Now().Add(time.Minute * 1),
		Token:    token,
	}
	OnlineUserMap.Store(u.Username, u)
}

// 清理过期用户
func cleanExpiredUsers() {
	OnlineUserMap.Range(func(key, value interface{}) bool {
		user := value.(*OnlineUserInfo)
		if user.Expires.Before(time.Now()) {
			OnlineUserMap.Delete(key.(string)) // 移除过期用户
		}
		return true
	})
}
