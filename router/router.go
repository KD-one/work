package router

import (
	"github.com/gin-gonic/gin"
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
	// appAuth
	router.POST("/addChangeAppAuth", middleware.ParseToken(), service.AddChangeAppAuth)
	router.GET("/getAppAuthTable", middleware.ParseToken(), service.GetAppAuthTable)
	router.POST("/getAppAuthByName", middleware.ParseToken(), service.GetAppAuthByName)
	// ECU
	router.POST("/ECUProjectAddChange", middleware.ParseToken(), service.ECUProjectAddChange)
	router.POST("/VerRecordAdd", middleware.ParseToken(), service.VerRecordAdd)
	router.POST("/VerRecordChange", middleware.ParseToken(), service.VerRecordChange)
	router.POST("/VerRecordRelease", middleware.ParseToken(), service.VerRecordRelease)
	//  file
	router.POST("/ECUSoftwareUpload", middleware.ParseToken(), service.UploadECUFile)
	router.POST("/ECUSoftwareDownload", middleware.ParseToken(), service.DownloadFile)
	router.POST("/ECUSoftwareCheckNewVer", middleware.ParseToken(), service.ECUSoftwareCheckNewVer)
	// log
	router.GET("/getAdminLog", middleware.ParseToken(), service.GetAdminLog)
	router.GET("/getClientLog", middleware.ParseToken(), service.GetClientLog)
	router.POST("/postClientLog", middleware.ParseToken(), service.PostClientLog)
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
	// 查看所有在线用户（需要权限）
	router.GET("/onlineUsers", service.OnlineUsers)

	// 获取数据库中所有表信息
	//router.GET("/getTables", middleware.AuthMiddleware(), service.GetTables)
	//// 获取表全量数据
	//router.GET("/getTableData", middleware.AuthMiddleware(), service.GetTableData)
	//// 创建数据表（全量请求）
	//router.POST("/createTable", middleware.AuthMiddleware(), service.CreateTable)
	//// 插入数据或修改(全量数据请求与获取)
	//router.POST("/insertOrUpdate", middleware.AuthMiddleware(), service.InsertOrUpdate)
}
