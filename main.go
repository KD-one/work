package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"test/config"
	"test/dao"
	"test/router"
	"test/service"
)

func main() {
	config.Init()

	// 初始化数据库
	dao.InitDB()

	// 初始化用户列表
	err := config.InitUserList(&service.UserList)
	if err != nil {
		panic("初始化用户列表失败")
	}

	// 每秒检查 UserList 中的用户是否过期
	go service.CheckUsersExpiration()

	// 每秒检查系统状态
	go service.CheckSystemStatus()

	// 每周统计数据
	go service.InitReportStatistics()

	// 创建默认路由
	r := gin.Default()

	// 注册路由
	router.RouterList(r)

	// 在指定地址端口开启服务
	r.Run(viper.GetString("port"))
}
