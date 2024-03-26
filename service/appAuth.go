package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"test/common"
	"test/dao"
	"test/model"
	"test/serializer"
)

type AddChangeAppAuthModel struct {
	AuthName  string `json:"AuthName" form:"AuthName" description:"权限名称"`
	H2Valves  int    `json:"H2Valves" form:"H2Valves"`
	H2SP      int    `json:"H2SP" form:"H2SP"`
	AirCmp    int    `json:"AirCmp" form:"AirCmp"`
	AirThrot  int    `json:"AirThrot" form:"AirThrot"`
	CoolFan   int    `json:"CoolFan" form:"CoolFan"`
	CoolET    int    `json:"CoolET" form:"CoolET"`
	CoolHT    int    `json:"CoolHT" form:"CoolHT"`
	Version   int    `json:"Version" form:"Version" description:"版本号"`
	ChangeLog string `json:"ChangeLog" form:"ChangeLog" description:"变更日志"`
}

type GetAppAuthByNameModel struct {
	AuthName string `json:"AuthName" form:"AuthName" description:"权限名称"`
}

func AddChangeAppAuth(c *gin.Context) {
	var data AddChangeAppAuthModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	//da
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	common.WriteLog(adminId, fmt.Sprintf("增加或修改app权限   用户名：%s   H2阀门：%d   H2SP：%d   AirCmp：%d   AirThrot：%d   CoolFan：%d   CoolET：%d   CoolHT：%d   版本号：%d   变更日志：%s", data.AuthName, data.H2Valves, data.H2SP, data.AirCmp, data.AirThrot, data.CoolFan, data.CoolET, data.CoolHT, data.Version, data.ChangeLog))

	app := model.Appauth{
		AuthName: data.AuthName,
		CoolET:   data.CoolET,
		AirCmp:   data.AirCmp,
		CoolFan:  data.CoolFan,
		AirThrot: data.AirThrot,
		H2Valves: data.H2Valves,
		H2SP:     data.H2SP,
		CoolHT:   data.CoolHT,
	}
	err := dao.DBAppAuthAddUpdate(adminId, app, data.Version, data.ChangeLog)
	if err != nil {
		c.JSON(422, serializer.Response{
			Code: 422,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
	})
}

func GetAppAuthTable(c *gin.Context) {
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	common.WriteLog(adminId, "尝试获取所有app权限信息")
	var appAuths []model.Appauth

	err := dao.DBAppAuthGetTable(adminId, &appAuths)
	if err != nil {
		c.JSON(422, serializer.Response{
			Code: 422,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: gin.H{
			"appAuthTable": appAuths,
		},
	})
}

func GetAppAuthByName(c *gin.Context) {
	var data GetAppAuthByNameModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	common.WriteLog(adminId, "尝试获取当前用户的app权限信息")

	var appAuth model.Appauth
	err := dao.DBFindAppAuthByName(adminId, data.AuthName, &appAuth)
	if err != nil {
		c.JSON(422, serializer.Response{
			Code: 422,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: gin.H{
			"appAuth": appAuth,
		},
	})
}
