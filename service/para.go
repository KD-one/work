package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"test/common"
	"test/dao"
	"test/model"
	"test/serializer"
)

type GetParaTableModel struct {
	TableName string `json:"TableName" form:"TableName" description:"para_auth表名"`
}

type DeleteParaModel struct {
	TableName string `json:"TableName" form:"TableName" description:"para_auth表名"`
	Version   int    `json:"Version" form:"Version" description:"版本号"`
	ChangeLog string `json:"ChangeLog" form:"ChangeLog" description:"变更日志"`
}

type AddChangeParaModel struct {
	TableName string           `json:"TableName" form:"TableName" description:"para_auth表名"`
	ParaData  []model.Paraauth `json:"ParaData" form:"ParaData" description:"para_auth数据"`
	Version   int              `json:"Version" form:"Version" description:"版本号"`
	ChangeLog string           `json:"ChangeLog" form:"ChangeLog" description:"变更日志"`
}

func GetParaTable(c *gin.Context) {
	var data GetParaTableModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}

	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	common.WriteLog(adminId, "获取参数表信息   表名："+data.TableName)

	var paraData []model.Paraauth

	err := dao.DBParaAuthGetTable(adminId, data.TableName, &paraData)
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
			"paraTableData": paraData,
		},
	})
}

func GetParaAuthList(c *gin.Context) {
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	common.WriteLog(adminId, "获取参数权限表列表")

	var paraTableNames []string

	err := dao.DBParaAuthGetTableList(adminId, &paraTableNames)
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
	var data AddChangeParaModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	common.WriteLog(adminId, fmt.Sprintf("添加或修改参数权限表   表名：%s   版本号：%d   修改记录：%s   数据：%+v", data.TableName, data.Version, data.ChangeLog, data.ParaData))

	err := dao.DBParaAuthAddUpdateTable(adminId, data.TableName, data.ParaData, data.Version, data.ChangeLog)
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
	var data DeleteParaModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	common.WriteLog(adminId, fmt.Sprintf("删除参数权限表   表名：%s   版本号：%d   修改记录：%s", data.TableName, data.Version, data.ChangeLog))

	err := dao.DBParaAuthDeleteTable(adminId, data.TableName, data.Version, data.ChangeLog)
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
