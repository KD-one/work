package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test/common"
)

//Casbin 权限管理

type CasbinInfo struct {
	Path   string `json:"path" form:"path"`
	Method string `json:"method" form:"method"`
}
type CasbinCreateRequest struct {
	RoleId      string       `json:"role_id" form:"role_id" description:"角色ID"`
	CasbinInfos []CasbinInfo `json:"casbin_infos" description:"权限模型列表"`
}

type CasbinGroup struct {
	Subject string `json:"subject" form:"subject" `
	Role    string `json:"role" form:"role"`
}

func AddPolicy(c *gin.Context) {
	common.UserRecord.Printf("==========")
	var params CasbinCreateRequest
	c.ShouldBind(&params)

	for _, v := range params.CasbinInfos {
		common.UserRecord.Println(params.RoleId, v.Path, v.Method)
		_, err := common.CasBin.AddPolicy(params.RoleId, v.Path, v.Method)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": err,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "添加策略成功",
	})
}
func AddGroupingPolicy(c *gin.Context) {
	common.UserRecord.Printf("==========")
	var cg CasbinGroup
	c.ShouldBind(&cg)

	common.UserRecord.Println("g", cg.Subject, cg.Role)
	_, err := common.CasBin.AddGroupingPolicy(cg.Subject, cg.Role)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "用户 - 角色 映射成功",
	})
}
