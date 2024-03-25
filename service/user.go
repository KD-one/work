package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"sync"
	"test/common"
	"test/dao"
	"test/model"
	"test/serializer"
	"time"
)

// OnlineUserInfo 保存用户登录相关信息
type OnlineUserInfo struct {
	Username string
	Password string
	Token    string // 自定义生成的token
	Expires  time.Time
}

// OnlineUserMap 用户在线状态存储
var OnlineUserMap sync.Map

// 锁
var mutex sync.Mutex

// UserList 用户列表
var UserList []model.UserList

// ClientList 客户端列表
var ClientList sync.Map

type LoginRequestModel struct {
	Name     string `json:"Name" form:"Name" description:"用户名"`
	Password string `json:"Password" form:"Password" description:"密码"`
	HostName string `json:"HostName" form:"HostName" description:"主机名"`
}

type DeleteUserModel struct {
	Name      string `json:"Name" form:"Name" description:"用户名"`
	Version   int    `json:"Version" form:"Version" description:"版本号"`
	ChangeLog string `json:"ChangeLog" form:"ChangeLog" description:"变更日志"`
}

type SendCommandModel struct {
	AdminMechineName    string `json:"AdminMechineName" form:"AdminMechineName" description:"管理员机器名"`
	Command             string `json:"Command" form:"Command" description:"命令"`
	ToClient            string `json:"ToClient" form:"ToClient" description:"客户端ip"`
	ToClientMachineName string `json:"ToClientMachineName" form:"ToClientMachineName" description:"客户端机器名"`
}

type AddChangeUserModel struct {
	Name     string `json:"Name" form:"Name" description:"用户名"`
	Password string `json:"Password" form:"Password" description:"密码"`
	//UserLevel int    `json:"UserLevel" form:"UserLevel" description:"用户权限"`
	AppAuth   string `json:"AppAuth" form:"AppAuth" description:"应用权限"`
	ParaAuth  string `json:"ParaAuth" form:"ParaAuth" description:"参数权限"`
	Version   int    `json:"Version" form:"Version" description:"版本号"`
	ChangeLog string `json:"ChangeLog" form:"ChangeLog" description:"变更日志"`
}

type ClientRequestModel struct {
	LabviewVersion    int    `json:"LabviewVersion" form:"LabviewVersion" description:"labview版本号"`
	EcuVersion        int    `json:"EcuVersion" form:"EcuVersion" description:"ecu版本号"`
	SystemStatus      string `json:"SystemStatus" form:"SystemStatus" description:"系统状态"`
	StartCount        int    `json:"StartCount" form:"StartCount" description:"启动次数"`
	FirstErrCode      int    `json:"FirstErrCode" form:"FirstErrCode" description:"第一次错误代码"`
	InstructionId     uint   `json:"InstructionId" form:"InstructionId" description:"指令id"`
	InstructionResult string `json:"InstructionResult" form:"InstructionResult" description:"指令结果"`
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

func Login(c *gin.Context) {

	//// TODO： 前端传输过来的数据为加密数据并置于body中，取出加密内容并解密后，将解密内容塞回body中（暂时用不到）
	//// 从c.Request.Body读出body数据并存入bodyBytes变量中
	//bodyBytes, _ := io.ReadAll(c.Request.Body)
	//// 将bodyBytes变量内容转换为字符串
	//bodyString := string(bodyBytes)
	//// 解密bodyString字符串
	//// ------此处执行解密函数------
	//// 使用解密字符串的[]byte类型，重置body
	//c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(bodyString)))

	var data LoginRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	name := data.Name
	password := data.Password

	fmt.Printf("用户请求登录   用户名：%s 密码：%s \n", name, password)
	common.WriteLog(0, fmt.Sprintf("用户请求登录   用户名：%s", name))

	if len(name) == 0 {
		c.JSON(422, serializer.Response{
			Code: 422,
			Msg:  "用户名不能为空",
		})
		return
	}

	if len(password) < 6 {
		c.JSON(422, serializer.Response{
			Code: 422,
			Msg:  "密码不能小于6位",
		})
		return
	}

	var userId uint
	err := dao.DBUserLogin(name, password, &userId)
	if err != nil {
		c.JSON(422, serializer.Response{
			Code: 422,
			Msg:  err.Error(),
		})
		return
	}

	//发放token
	token, err := common.ReleaseToken(userId)
	if err != nil {
		c.JSON(500, serializer.Response{
			Code: 500,
			Msg:  "token生成失败",
		})
		return
	}

	expireTime := viper.GetDuration("user.expiration")
	mutex.Lock()
	for i := 0; i < len(UserList); i++ {
		if UserList[i].Name == name {

			UserList[i].Online = true
			UserList[i].Expiration = time.Now().Add(expireTime * time.Second)
			UserList[i].HostName = data.HostName

			break
		}
	}
	mutex.Unlock()

	// 通过协程添加用户到在线用户中
	//go addOnlineUser(name, password, token)
	//
	//// 启动协程持续检查用户是否过期
	//go func() {
	//	for {
	//		cleanExpiredUsers()
	//		time.Sleep(time.Second * 20) // 每20秒检查一次
	//	}
	//}()
	common.UserRecord.Println(" [info] 登录成功", name)
	//返回结果
	c.JSON(200, serializer.Response{
		Code: 200,
		Data: gin.H{
			"token": token,
		},
		Msg: "success",
	})
}

func GetUserList(c *gin.Context) {
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 日志记录
	common.WriteLog(adminId, "获取用户列表")

	var users []model.User
	err := dao.DBUserGetTable(adminId, &users)
	if err != nil {
		c.JSON(422, serializer.Response{
			Code: 422,
			Msg:  err.Error(),
		})
		return
	}
	// 记录新列表users中的所有用户名
	tempMap := make(map[string]bool)
	// 对UserList中没有的元素进行添加并更新已有用户信息
	mutex.Lock()
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
		tempMap[user.Name] = true
	}
	mutex.Unlock()
	// 对UserList中多出来的元素进行删除
	var newUserList []model.UserList
	for _, user := range UserList {
		if tempMap[user.Name] { // 存在于新列表中，保留
			newUserList = append(newUserList, user)
		}
	}
	mutex.Lock()
	UserList = newUserList
	mutex.Unlock()

	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: gin.H{
			"userList": UserList,
		},
	})
}

func AddChangeUser(c *gin.Context) {
	var data AddChangeUserModel
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
	common.WriteLog(adminId, fmt.Sprintf("增加或修改用户   用户名：%s   密码：%s   操作权限：%s   参数权限：%s   版本号：%d   修改记录：%s", data.Name, data.Password, data.AppAuth, data.ParaAuth, data.Version, data.ChangeLog))

	user := model.User{
		Name:     data.Name,
		Password: data.Password,
		AppAuth:  data.AppAuth,
		ParaAuth: data.ParaAuth,
	}
	err := dao.DBUserAddUpdate(adminId, user, data.Version, data.ChangeLog)
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

func DeleteUser(c *gin.Context) {
	var data DeleteUserModel
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
	common.WriteLog(adminId, fmt.Sprintf("删除用户   用户名：%s   版本号：%d   修改记录：%s", data.Name, data.Version, data.ChangeLog))

	err := dao.DBUserDelete(adminId, data.Name, data.Version, data.ChangeLog)
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

// 添加用户到在线用户中
func addOnlineUser(name, password, token string) {
	u := &OnlineUserInfo{
		Username: name,
		Password: password,
		Expires:  time.Now().Add(time.Minute * 1),
		Token:    token,
	}
	OnlineUserMap.Store(u.Username, u)
}

// 清理过期用户
func cleanExpiredUsers() {
	OnlineUserMap.Range(func(key, value interface{}) bool {
		user := value.(*OnlineUserInfo)
		if user.Expires.Before(time.Now()) {
			OnlineUserMap.Delete(key.(string)) // 移除过期用户
		}
		return true
	})
}

func GetDBTableVersion(c *gin.Context) {
	var version int
	err := dao.DBTableGetLatestVer(&version)
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
			"version": version,
		},
	})
}

// SendCommand 管理端发送指令
func SendCommand(c *gin.Context) {
	var data SendCommandModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	// 记录日志
	logData := fmt.Sprintf("发送命令   发送命令；%s   管理端机器名：%s   客户：%s   客户端机器名：%s", data.Command, data.AdminMechineName, data.ToClient, data.ToClientMachineName)
	common.WriteLog(adminId, logData)

	// 构造指令记录存入数据库
	adminName := dao.FindUserName(adminId)
	instruction := model.Instruction{
		CreateTime:        time.Now().Format("2006-01-02 15:04:05"),
		AdminName:         adminName,
		AdminMechineName:  data.AdminMechineName,
		ClientName:        data.ToClient,
		ClientMachineName: data.ToClientMachineName,
	}
	id, err := dao.InsertInstructionReturnId(adminId, instruction)
	if err != nil {
		c.JSON(422, serializer.Response{
			Code: 422,
			Msg:  err.Error(),
		})
		return
	}
	// 将信息存放全局客户列表中
	ClientList.Store(data.ToClient, model.ClientList{
		InstructionId:     id,
		InstructionResult: data.Command,
	})

	fmt.Println("ClientList: ------------", ClientList, "----------------")
	fmt.Println("UserList: ------------", UserList, "----------------")

	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
	})
}

// Keepalive 客户端发送心跳并携带数据
func Keepalive(c *gin.Context) {
	var data ClientRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	clientIdAny, _ := c.Get("userId")
	clientId := clientIdAny.(uint)
	// 记录日志
	logData := fmt.Sprintf("keepAlive   labview版本号：%d   ECU版本号；%d   系统状态：%s   启动次数：%d   第一次错误代码：%d", data.LabviewVersion, data.EcuVersion, data.SystemStatus, data.StartCount, data.FirstErrCode)
	common.WriteLog(clientId, logData)

	// 根据指令id更新记录的指令结果
	var instruction model.Instruction
	if data.InstructionId == 0 {
		common.WriteLog(clientId, "指令id不存在")
	} else {
		err := dao.FindInstruction(data.InstructionId, &instruction)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  err.Error(),
			})
			return
		}
		instruction.Result = data.InstructionResult
		instruction.ResultCreateTime = time.Now().Format("2006-01-02 15:04:05")
		err = dao.UpdateInstructionRecord(data.InstructionId, instruction)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  err.Error(),
			})
			return
		}
	}
	clientName := dao.FindUserName(clientId)
	expireTime := viper.GetDuration("user.expiration")
	fmt.Println("UserList: ------------", UserList, "----------------")
	// 更新客户端过期时间
	mutex.Lock()
	for i := 0; i < len(UserList); i++ {
		if UserList[i].Name == clientName {

			UserList[i].Online = true
			UserList[i].Expiration = time.Now().Add(expireTime * time.Second)

			break
		}
	}
	mutex.Unlock()

	// 查找全局client列表中是否存在该client，不存在直接返回空内容，存在则返回内容并删除client列表中的该client
	i, ok := ClientList.Load(clientName)
	if !ok {
		c.JSON(200, serializer.Response{
			Code: 200,
			Msg:  "success",
			Data: model.ClientList{
				InstructionId:     0,
				InstructionResult: "",
			},
		})
		return
	}
	ins := i.(model.ClientList)
	ClientList.Delete(clientName)
	fmt.Println("ClientList: ------------", ClientList, "----------------")
	fmt.Println("UserList: ------------", UserList, "----------------")
	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: model.ClientList{
			InstructionId:     ins.InstructionId,
			InstructionResult: ins.InstructionResult,
		},
	})
}

// CheckUsersExpiration 不断检查 UserList 中的用户是否过期
func CheckUsersExpiration() {
	for {
		// 获取当前时间
		now := time.Now()

		mutex.Lock()
		// 遍历用户列表
		for i := range UserList {
			user := &UserList[i]
			// 如果用户在线且已经过期
			if user.Online && user.Expiration.Before(now) {
				// 将 Online 字段设为 false
				user.Online = false
			}
		}
		mutex.Unlock()

		// 延迟一段时间后再检查
		time.Sleep(1 * time.Second)
	}
}
