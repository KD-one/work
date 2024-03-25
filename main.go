package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"test/common"
	"test/config"
	"test/router"
	"test/service"
)

func main() {
	config.Init()

	// 初始化数据库
	common.InitDB()

	// 初始化casbin数据库
	//common.InitCasbinDB()

	// 日志相关
	common.Log()
	defer common.F.Close()

	// 初始化用户列表
	err := config.InitUserList(&service.UserList)
	if err != nil {
		panic("初始化用户列表失败")
	}

	// 每秒检查系统状态
	go service.CheckUsersExpiration()

	// 创建默认路由
	r := gin.Default()

	//// 注册全局模板
	//r.LoadHTMLGlob("template/**/*")
	//
	//// 配置静态文件服务
	//r.Static("/images", "./images")

	// 注册路由
	router.RouterList(r)

	// 在指定地址端口开启服务
	r.Run(viper.GetString("port"))
}
