package service

import (
	"github.com/gin-gonic/gin"
	"test/dao"
	"test/model"
	"test/serializer"
)

var UserList []model.UserList

type RequestModel struct {
	Name      string `json:"name" form:"name" description:"用户名"`
	Password  string `json:"password" form:"password" description:"密码"`
	UserLevel int    `json:"user_level" form:"user_level" description:"用户权限"`
	AppAuth   string `json:"app_auth" form:"app_auth" description:"应用权限"`
	ParaAuth  string `json:"para_auth" form:"para_auth" description:"参数权限"`
	Version   int    `json:"version" form:"version" description:"版本号"`
	ChangeLog string `json:"changelog" form:"changelog" description:"变更日志"`
}

func OnlineUsers(c *gin.Context) {
	OnlineUserMap.Range(func(k, v interface{}) bool {
		//fmt.Println("key>>>>>>>>: ", k, "             value>>>>>>>>: ", v.(*User))
		c.JSON(200, gin.H{
			"key":   k.(string),
			"value": v.(*OnlineUserInfo),
		})
		return true
	})
	//c.JSON(200, gin.H{
	//	"message": "success",
	//	"data":    OnlineUserMap,
	//})
}

func (r *RequestModel) AddChangeUser(c *gin.Context) serializer.Response {
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	//var user model.User
	//err := json.Unmarshal([]byte(c.PostForm("user")), &user)
	//if err != nil {
	//	c.JSON(404, gin.H{
	//		"message": "fail",
	//		"data":    "解析json数据时出错",
	//		"err":     err.Error(),
	//	})
	//	return
	//}
	//versionString := c.PostForm("version")
	//version, _ := strconv.Atoi(versionString)
	//changelog := c.PostForm("changelog")
	user := model.User{
		Name:      r.Name,
		Password:  r.Password,
		UserLevel: r.UserLevel,
		AppAuth:   r.AppAuth,
		ParaAuth:  r.ParaAuth,
	}
	err := dao.DBUserAddUpdate(adminId, user, r.Version, r.ChangeLog)
	if err != nil {
		return serializer.Response{
			Code: 422,
			Msg:  "dao层出错",
			Data: gin.H{
				"err": err.Error(),
			},
		}
	}
	return serializer.Response{
		Code: 200,
		Msg:  "success",
	}
}
func (r *RequestModel) GetUserList(c *gin.Context) serializer.Response {
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	var users []model.User
	err := dao.DBUserGetTable(adminId, &users)
	if err != nil {
		return serializer.Response{
			Code: 422,
			Msg:  err.Error(),
		}
	}
	for _, user := range users {
		exist := false
		for i := 0; i < len(UserList); i++ {
			if user.Name == UserList[i].Name {
				UserList[i].AppAuth = user.AppAuth
				UserList[i].ParaAuth = user.ParaAuth
				exist = true
				break
			}
		}
		if !exist {
			UserList = append(UserList, model.UserList{
				Name:     user.Name,
				AppAuth:  user.AppAuth,
				ParaAuth: user.ParaAuth,
			})
		}
	}
	return serializer.Response{
		Code: 200,
		Data: UserList,
		Msg:  "success",
	}
}
