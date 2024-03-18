package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"test/common"
	"test/dao"
	"test/model"
	"test/serializer"
)

type ParaRequestModel struct {
	TableName string           `json:"table_name" form:"table_name" description:"para_auth表名"`
	ParaData  []model.Paraauth `json:"para_data" form:"para_data" description:"para_auth数据"`
	Version   int              `json:"version" form:"version" description:"版本号"`
	ChangeLog string           `json:"changelog" form:"changelog" description:"变更日志"`
}

func GetParaTable(c *gin.Context) {
	var data ParaRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 422,
			Msg:  "参数绑定时出错",
		})
		return
	}

	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	tableName := "表名：" + data.TableName
	logData := common.DataToString(adminId, "获取参数表信息", tableName)
	err := common.WriteStringToLog(logData)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	err = dao.DBParaAuthGetTable(adminId, data.TableName, &data.ParaData)
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
		Data: gin.H{
			"paraTableData": data.ParaData,
		},
	})
}

func GetParaAuthList(c *gin.Context) {
	var data ParaRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 422,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	logData := common.DataToString(adminId, "获取参数权限表列表")
	err := common.WriteStringToLog(logData)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	var paraTableNames []string

	err = dao.DBParaAuthGetTableList(adminId, &paraTableNames)
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
		Data: gin.H{
			"paraAuthTables": paraTableNames,
		},
	})
}

func AddChangePara(c *gin.Context) {
	var data ParaRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 422,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	tableName := "表名：" + data.TableName
	version := fmt.Sprintf("版本号：%d", data.Version)
	userChangeLog := "修改记录：" + data.ChangeLog
	tableData := "数据：" + fmt.Sprintf("%+v", data.ParaData)
	logData := common.DataToString(adminId, "添加或修改参数权限表", tableName, version, userChangeLog, tableData)
	err := common.WriteStringToLog(logData)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	err = dao.DBParaAuthAddUpdateTable(adminId, data.TableName, data.ParaData, data.Version, data.ChangeLog)
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

func DeletePara(c *gin.Context) {
	var data ParaRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 422,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	tableName := "表名：" + data.TableName
	version := fmt.Sprintf("版本号：%d", data.Version)
	userChangeLog := "修改记录：" + data.ChangeLog
	logData := common.DataToString(adminId, "删除参数权限表", tableName, version, userChangeLog)
	err := common.WriteStringToLog(logData)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	err = dao.DBParaAuthDeleteTable(adminId, data.TableName, data.Version, data.ChangeLog)
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
