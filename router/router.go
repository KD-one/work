package router

import (
	"github.com/gin-gonic/gin"
	"test/FileRelated"
	"test/api"
	"test/middleware"
	"test/service"
)

func RouterList(router *gin.Engine) {

	// user
	router.POST("/login", service.Login)
	router.GET("/getUserList", middleware.ParseToken(), service.GetUserList)
	router.POST("/addChangeUser", middleware.ParseToken(), service.AddChangeUser)
	router.POST("/deleteUser", middleware.ParseToken(), service.DeleteUser)
	router.POST("/sendCommand", middleware.ParseToken(), service.SendCommand)
	router.POST("/keepalive", middleware.ParseToken(), service.Keepalive)
	// paraAuth
	router.GET("/getParaList", middleware.ParseToken(), service.GetParaAuthList)
	router.POST("/getParaTable", middleware.ParseToken(), service.GetParaTable)
	router.POST("/addChangePara", middleware.ParseToken(), service.AddChangePara)
	router.POST("/deletePara", middleware.ParseToken(), service.DeletePara)
	// 获取最新表操作记录
	router.GET("/getDBTableVersion", middleware.ParseToken(), service.GetDBTableVersion)

	//router.GET("/compareMemoryUsedPercent", service.CompareMemoryUsedPercent)
	//router.GET("/compareCpuUsedPercent", service.CompareCpuUsedPercent)
	//router.GET("/compareNetUsedPercent", service.CompareNetUsedPercent)
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	// 首页
	router.GET("/index", service.Index)
	router.GET("/visitorIndex", service.VisitorIndex)

	// 用户登录（发放权限）
	router.GET("/", service.ToLogin)

	// 注册用户（需要权限）
	router.GET("/toRegister", service.ToRegister)
	router.POST("/register", middleware.AuthMiddleware(), service.Register)

	// 多文件上传（需要权限）
	router.GET("/toUpload", FileRelated.ToFileUpload)
	router.POST("/uploads", middleware.AuthMiddleware(), FileRelated.UploadFiles)

	// 所有文件列表页面（需要权限）
	router.GET("/toDownload", FileRelated.ToDownload)
	router.GET("/showFiles", middleware.AuthMiddleware(), FileRelated.ShowFileList)

	// 查看所有在线用户（需要权限）
	router.GET("/onlineUsers", service.OnlineUsers)

	// 更新文件（需要权限）
	router.GET("/toUpdate", api.ToUpdate)
	router.POST("/update", middleware.AuthMiddleware(), api.Update)

	// 获取数据库中所有表信息
	router.GET("/getTables", middleware.AuthMiddleware(), service.GetTables)
	// 获取表全量数据
	router.GET("/getTableData", middleware.AuthMiddleware(), service.GetTableData)
	// 创建数据表（全量请求）
	router.POST("/createTable", middleware.AuthMiddleware(), service.CreateTable)
	// 插入数据或修改(全量数据请求与获取)
	router.POST("/insertOrUpdate", middleware.AuthMiddleware(), service.InsertOrUpdate)

	// 游客有权限查看到的文件列表
	router.GET("/visitorToDownload", FileRelated.VisitorFileList)

	// 文件更改日志
	router.GET("/toListLogFiles", service.ListLogFiles)
	router.GET("/viewLogFile/:filename", service.ViewLogFile)

	// 通过请求参数下载文件
	router.GET("/paramToDownload", FileRelated.ParamToDownload)
	router.GET("/download", FileRelated.FileDownloadService)

	// 通过项目名查询历史版本
	router.GET("/toVersion", service.BranchToVersion1)
	router.GET("/versions", service.BranchToVersion2)

	// 通用路由
	// 点击文件下载
	router.GET("/:file", FileRelated.DownloadFileFromParam)

	// 测试身份认证中间件
	router.GET("/without_auth", service.ServiceWithoutAuth)
	// 查看用户信息
	router.GET("/info", middleware.AuthMiddleware(), service.Info)
	//// 单文件上传
	//router.POST("/upload", FileRelated.UploadFile)
	//// 添加策略
	//router.POST("/addPolicy", middleware.CasbinMiddleware(), service.AddPolicy)
	//// 添加组策略
	//router.POST("/addGroupPolicy", service.AddGroupingPolicy)
	//// 根据请求参数下载文件
	//router.GET("/urlDownload", FileRelated.UrlFileDownloadService)
	//// 根据路径参数下载文件
	//router.GET("/api/ecu/bf", FileRelated.BuildFileDownloadService)
	//router.GET("/api/ecu/a2l", FileRelated.A2lFileDownloadService)
	// 错误路由
	//router.GET("/404", service.NotFound)
	// 路由分组
	//authRouter := router.Group("auth")
	//{
	//	//以下的接口，都使用Authorize()中间件身份验证
	//	authRouter.Use(service.Authorize())
	//	// 测试身份认证中间件
	//	authRouter.GET("/with_auth", service.ServiceWithAuth)
	//}
}
