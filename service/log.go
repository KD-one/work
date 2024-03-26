package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"test/common"
	"test/dao"
	"test/model"
	"test/serializer"
)

// PostClientLogModel 接受请求参数
type PostClientLogModel struct {
	Time      string `json:"Time" form:"Time"`
	PCName    string `json:"PCName" form:"PCName"`
	ChangeLog string `json:"ChangeLog" form:"ChangeLog"`
}

//type GetClientLogResponseModel struct {
//	Id        uint   `json:"id" form:"id"`
//	UserName  string `json:"UserName" form:"UserName"`
//	Time      string `json:"Time" form:"Time"`
//	PCName    string `json:"PCName" form:"PCName"`
//	ChangeLog string `json:"ChangeLog" form:"ChangeLog"`
//}

// ListLogFiles 展示log文件列表
func ListLogFiles(c *gin.Context) {
	dirPath := "./log/uploadRecord"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".log" { // 只显示.log文件
			fileNames = append(fileNames, file.Name())
		}
	}

	c.HTML(http.StatusOK, "logView/log.html", gin.H{"files": fileNames})
}

// ViewLogFile 显示具体log文件内容
func ViewLogFile(c *gin.Context) {
	filename := c.Param("filename")
	dirPath := "./log/uploadRecord"
	fullPath := filepath.Join(dirPath, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", data) // 假设文件是纯文本格式，否则需要调整MIME类型
}

func GetAdminLog(c *gin.Context) {
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	common.WriteLog(adminId, "尝试获取管理员版本日志")
	var tablever []model.Tablever
	err := dao.DBGetLimitVersionTable(adminId, &tablever)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: gin.H{"table": tablever},
	})
}

func GetClientLog(c *gin.Context) {
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	common.WriteLog(adminId, "尝试获取客户端版本日志")
	var cLog []model.ClientLog
	err := dao.DBGetLimitClientLogTable(adminId, &cLog)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: gin.H{"table": cLog},
	})
}

func PostClientLog(c *gin.Context) {
	var data PostClientLogModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	adminName := dao.FindUserName(adminId)

	common.WriteLog(adminId, fmt.Sprintf("用户名：%s   主机名：%s   时间：%s   更改日志：%s", adminName, data.PCName, data.Time, data.ChangeLog))

	if data.PCName == "" || data.Time == "" || data.ChangeLog == "" {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数不能为空",
		})
	}

	clientlog := model.ClientLog{
		UserName:  adminName,
		PCName:    data.PCName,
		ChangeLog: data.ChangeLog,
		Time:      data.Time,
	}

	err := dao.DBClientLogAdd(clientlog)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
	})

}
