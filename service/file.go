package service

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"test/common"
	"test/dao"
	"test/model"
	"test/serializer"
	"time"
)

//type UploadFileRequestModel struct {
//	SoftwareBranch string `json:"SoftwareBranch" form:"SoftwareBranch" description:"软件分支"`
//	SoftwareVerNum string `json:"SoftwareVerNum" form:"SoftwareVerNum" description:"软件版本号"`
//	CalMainNum     string `json:"CalMainNum" form:"CalMainNum" description:"cal_main号"`
//	CalSubNum      string `json:"CalSubNum" form:"CalSubNum" description:"cal_sub号"`
//	ChangeLog      string `json:"ChangeLog" form:"ChangeLog" description:"变更日志"`
//	CalCmd         string `json:"CalCmd" form:"CalCmd" description:"cal_cmd"`
//	BuildFileMD5   string `json:"BuildFileMD5" form:"BuildFileMD5" description:"build文件md5"`
//	A2LFileMD5     string `json:"A2LFileMD5" form:"A2LFileMD5" description:"a2l文件md5"`
//	//BuildFileName  string `json:"BuildFileName" form:"BuildFileName" description:"build文件名"`
//	//A2LFileName    string `json:"A2LFileName" form:"A2LFileName" description:"a2l文件名"`
//	//Files          []*multipart.FileHeader `json:"file" form:"file" description:"文件"`
//}

//func UploadECUFile(c *gin.Context) {
//	var data UploadFileRequestModel
//	if err := c.ShouldBind(&data); err != nil {
//		c.JSON(400, serializer.Response{
//			Code: 400,
//			Msg:  "参数绑定时出错",
//		})
//		return
//	}
//	fmt.Println("11111111111111111111111111111111111111111111111111")
//	// 获取文件
//	form, _ := c.MultipartForm()
//	buildFile := form.File["BuildFile"]
//	a2lFile := form.File["A2LFile"]
//	if len(buildFile) == 0 || len(a2lFile) == 0 || len(buildFile) > 1 || len(a2lFile) > 1 {
//		c.JSON(400, serializer.Response{
//			Code: 400,
//			Msg:  "文件不能为空,或文件数量不能大于1",
//		})
//		return
//	}
//	//files := data.Files
//	fmt.Println("222222222222222222222222222222222222222222222222222")
//
//	adminIdAny, _ := c.Get("userId")
//	adminId := adminIdAny.(uint)
//	adminName := dao.FindUserName(adminId)
//
//	common.WriteLog(adminId, fmt.Sprintf("上传文件   软件分支：%s   软件版本号：%s   cal_main号：%s   cal_sub号：%s   变更日志：%s   cal_cmd：%s   build文件名：%s   a2l文件名：%s", data.SoftwareBranch, data.SoftwareVerNum, data.CalMainNum, data.CalSubNum, data.ChangeLog, data.CalCmd, buildFile[0].Filename, a2lFile[0].Filename))
//
//	if data.SoftwareBranch == "" || data.SoftwareVerNum == "" || data.CalMainNum == "" || data.CalSubNum == "" || data.ChangeLog == "" { // || data.BuildFileName == "" || data.A2LFileName == ""
//		c.JSON(400, serializer.Response{
//			Code: 400,
//			Msg:  "参数不能为空",
//		})
//		return
//	}
//
//	branch, _ := strconv.ParseUint(data.SoftwareBranch, 10, 64)
//	version, _ := strconv.ParseUint(data.SoftwareVerNum, 10, 64)
//	calmain, _ := strconv.ParseUint(data.CalMainNum, 10, 64)
//	calsub, _ := strconv.ParseUint(data.CalSubNum, 10, 64)
//
//	fmt.Println("333333333333333333333333333333333333333333333333333333")
//	// 查询记录
//	RowsAffected := dao.DBCheckECUFileMapRecordExists(uint(branch), uint(version))
//
//	// 查到相同的记录
//	if RowsAffected != 0 {
//		c.JSON(400, serializer.Response{
//			Code: 400,
//			Msg:  "记录已存在，请更换分支或版本",
//		})
//		return
//	}
//
//	// 数据库插入记录
//	t := time.Now().Format("2006-01-02 15:04:05")
//	ecuFileMap := model.EcuFileMap{
//		UserName:   adminName,
//		CreateTime: t,
//		Branch:     uint(branch),
//		Version:    uint(version),
//		BuildFile:  buildFile[0].Filename,
//		A2lFile:    a2lFile[0].Filename,
//		CalMain:    uint(calmain),
//		CalSub:     uint(calsub),
//		ChangeLog:  data.ChangeLog,
//		CalCmd:     data.CalCmd,
//	}
//	err := dao.DBCreateECUFileMapRecord(adminId, ecuFileMap)
//	if err != nil {
//		c.JSON(400, serializer.Response{
//			Code: 400,
//			Msg:  err.Error(),
//		})
//		return
//	}
//
//	// 将文件添加到指定目录
//	//var fileNames []string
//	//for _, file := range files {
//	//fmt.Println("=====上传的文件名：", file.Filename)
//	//fileNames = append(fileNames, "upload/"+file.Filename)
//	buildFilePath := "upload/" + buildFile[0].Filename // 设置保存文件的路径
//	a2lFilePath := "upload/" + a2lFile[0].Filename     // 设置保存文件的路径
//
//	// 查看file_path是否存在
//	_, err = os.Stat(buildFilePath)
//	// file_path路径下不存在此文件，保存
//	if err != nil {
//		err = c.SaveUploadedFile(buildFile[0], buildFilePath) // 保存文件
//		if err != nil {
//			c.JSON(400, serializer.Response{
//				Code: 400,
//				Msg:  "请求失败",
//			})
//			return
//		}
//	}
//	_, err = os.Stat(a2lFilePath)
//	// file_path路径下不存在此文件，保存
//	if err != nil {
//		err = c.SaveUploadedFile(a2lFile[0], a2lFilePath) // 保存文件
//		if err != nil {
//			c.JSON(400, serializer.Response{
//				Code: 400,
//				Msg:  "请求失败",
//			})
//			return
//		}
//	}
//	//}
//
//	c.JSON(200, serializer.Response{
//		Code: 200,
//		Msg:  "success",
//	})
//}

type UploadFileRequestModel struct {
	SoftwareBranch string                `json:"SoftwareBranch" form:"SoftwareBranch" description:"软件分支"`
	SoftwareVerNum string                `json:"SoftwareVerNum" form:"SoftwareVerNum" description:"软件版本号"`
	CalMainNum     string                `json:"CalMainNum" form:"CalMainNum" description:"cal_main号"`
	CalSubNum      string                `json:"CalSubNum" form:"CalSubNum" description:"cal_sub号"`
	ChangeLog      string                `json:"ChangeLog" form:"ChangeLog" description:"变更日志"`
	ChangeLogEN    string                `json:"ChangeLogEN" form:"ChangeLogEN" description:"变更日志英文"`
	CalCmd         string                `json:"CalCmd" form:"CalCmd" description:"cal_cmd"`
	BuildFileMD5   string                `json:"BuildFileMD5" form:"BuildFileMD5" description:"build文件md5"`
	A2LFileMD5     string                `json:"A2LFileMD5" form:"A2LFileMD5" description:"a2l文件md5"`
	BuildFile      *multipart.FileHeader `json:"BuildFile" form:"BuildFile" description:"build文件"`
	A2LFile        *multipart.FileHeader `json:"A2LFile" form:"A2LFile" description:"a2l文件"`
	//Files          []*multipart.FileHeader `json:"file" form:"file" description:"文件"`
}

type DownloadFileRequestModel struct {
	SoftwareBranch string `json:"SoftwareBranch" form:"SoftwareBranch" description:"软件分支"`
	SoftwareVerNum string `json:"SoftwareVerNum" form:"SoftwareVerNum" description:"软件版本号"`
}

func UploadECUFile(c *gin.Context) {
	var data UploadFileRequestModel
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

	common.WriteLog(adminId, fmt.Sprintf("上传文件   软件分支：%s   软件版本号：%s   cal_main号：%s   cal_sub号：%s   变更日志：%s   cal_cmd：%s   build文件名：%s   a2l文件名：%s   build文件md5: %s   a2l文件md5: %s", data.SoftwareBranch, data.SoftwareVerNum, data.CalMainNum, data.CalSubNum, data.ChangeLog, data.CalCmd, data.BuildFile.Filename, data.A2LFile.Filename, data.BuildFileMD5, data.A2LFileMD5))

	if data.SoftwareBranch == "" || data.SoftwareVerNum == "" || data.CalMainNum == "" || data.CalSubNum == "" || data.ChangeLog == "" || data.BuildFile == nil || data.A2LFile == nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数不能为空",
		})
		return
	}

	// 校验文件类型
	if strings.HasSuffix(data.BuildFile.Filename, "a2l") || !strings.HasSuffix(data.A2LFile.Filename, "a2l") || data.A2LFile.Filename == data.BuildFile.Filename {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "上传文件类型错误",
		})
		return
	}

	// 验证a2l文件md5
	a2lFile, _ := data.A2LFile.Open()
	defer a2lFile.Close()
	// 创建一个足够大的缓冲区来存储文件内容
	buffer := new(bytes.Buffer)
	io.Copy(buffer, a2lFile)
	fileContent := buffer.Bytes()
	if MD5(fileContent) != data.A2LFileMD5 {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "a2l文件md5校验失败，可能为网络传输出错，请重新上传",
		})
		return
	}
	// 验证build文件md5
	buildFile, _ := data.BuildFile.Open()
	defer buildFile.Close()
	buffer = new(bytes.Buffer)
	io.Copy(buffer, buildFile)
	fileContent = buffer.Bytes()
	if MD5(fileContent) != data.BuildFileMD5 {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "build文件md5校验失败，可能为网络传输出错，请重新上传",
		})
		return
	}

	// 验证a2l文件和build文件的文件名相同
	s1 := strings.Split(data.A2LFile.Filename, ".")
	s2 := strings.Split(data.BuildFile.Filename, ".")
	if s1[0] != s2[0] {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "a2l文件和build文件的文件名不相同，请重新上传",
		})
		return
	}

	branch, _ := strconv.ParseUint(data.SoftwareBranch, 10, 64)
	version, _ := strconv.ParseUint(data.SoftwareVerNum, 10, 64)
	calmain, _ := strconv.ParseUint(data.CalMainNum, 10, 64)
	calsub, _ := strconv.ParseUint(data.CalSubNum, 10, 64)

	var ecuRecord model.EcuFileMap

	// 查询记录
	RowsAffected := dao.DBCheckECUFileMapRecordExists(uint(branch), uint(version))
	// 查到相同的记录
	if RowsAffected != 0 {
		err := dao.DBFindECUFileMapRecord(adminId, uint(branch), uint(version), &ecuRecord)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  err.Error(),
			})
			return
		}
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  fmt.Sprintf("记录已存在, %s 在时间 %s 时上传，修改记录为：%s", ecuRecord.UserName, ecuRecord.CreateTime, ecuRecord.ChangeLog),
		})
		return
	}

	// 将文件添加到指定目录
	//var fileNames []string
	//for _, file := range files {
	//fmt.Println("=====上传的文件名：", file.Filename)
	//fileNames = append(fileNames, "upload/"+file.Filename)
	buildFilePath := "upload/" + data.BuildFile.Filename // 设置保存文件的路径
	a2lFilePath := "upload/" + data.A2LFile.Filename     // 设置保存文件的路径

	_, err := os.Stat(buildFilePath)
	if err == nil {
		err = os.Remove(buildFilePath)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "服务器删除文件时失败",
			})
			return
		}
	}

	err = c.SaveUploadedFile(data.BuildFile, buildFilePath)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "服务器写入文件时失败",
		})
		return
	} // 保存文件

	_, err = os.Stat(a2lFilePath)
	if err == nil {
		err = os.Remove(a2lFilePath)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "服务器删除文件时失败",
			})
			return
		}
	}
	err = c.SaveUploadedFile(data.A2LFile, a2lFilePath)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "服务器写入文件时失败",
		})
		return
	}
	//}

	// 文件压缩
	files := []string{buildFilePath, a2lFilePath}
	output := "ECUSoftware/" + s1[0] + ".zip"

	// 判断upload文件夹中是否存在指定的zip文件
	_, err = os.Stat(output)
	if err != nil {
		// 文件不存在，则将文件压缩打包
		if err = common.ZipFiles(output, files); err != nil {
			panic(err)
		}
	} else {
		err = dao.DBWhereBuildFileFindRecord(path.Base(output), &ecuRecord)
		if err != nil {
			// 模拟文件保存和插入数据库两个操作绑定为原子操作
			os.Remove(buildFilePath)
			os.Remove(a2lFilePath)
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  err.Error(),
			})
			return
		}
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  fmt.Sprintf("文件已存在, %s 在时间 %s 时上传了 %s，修改记录为：%s", ecuRecord.UserName, ecuRecord.CreateTime, ecuRecord.BuildFile, ecuRecord.ChangeLog),
		})
		return
	}

	// 数据库插入记录
	t := time.Now().Format("2006-01-02 15:04:05")
	ecuFileMap := model.EcuFileMap{
		UserName:    adminName,
		CreateTime:  t,
		Branch:      uint(branch),
		Version:     uint(version),
		BuildFile:   path.Base(output),
		CalMain:     uint(calmain),
		CalSub:      uint(calsub),
		CalCmd:      data.CalCmd,
		ChangeLog:   data.ChangeLog,
		ChangeLogEn: data.ChangeLogEN,
	}
	err = dao.DBCreateECUFileMapRecord(adminId, ecuFileMap)
	if err != nil {
		// 模拟文件保存和插入数据库两个操作绑定为原子操作
		os.Remove(buildFilePath)
		os.Remove(a2lFilePath)
		os.Remove(output)
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

func DownloadFile(c *gin.Context) {
	var data DownloadFileRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}

	branch, _ := strconv.ParseUint(data.SoftwareBranch, 10, 64)
	version, _ := strconv.ParseUint(data.SoftwareVerNum, 10, 64)
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	common.WriteLog(adminId, fmt.Sprintf("下载文件   软件分支：%s   软件版本号：%s", data.SoftwareBranch, data.SoftwareVerNum))

	var ecuFileMap model.EcuFileMap

	RowsAffected := dao.DBCheckECUFileMapRecordExists(uint(branch), uint(version))
	if RowsAffected == 0 {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "没有找到对应的记录，请检查参数是否正确",
		})
		return
	}
	err := dao.DBFindECUFileMapRecordNoAuth(uint(branch), uint(version), &ecuFileMap)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	//定义文件路径
	filePath := "ECUSoftware/" + ecuFileMap.BuildFile

	_, err = os.Stat(filePath)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "文件不存在，请检查参数是否正确",
		})
		return
	}

	// 设置请求头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+path.Base(filePath))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	// 下载文件
	c.File(filePath)
}

func ECUSoftwareCheckNewVer(c *gin.Context) {
	var data DownloadFileRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}

	branch, _ := strconv.ParseUint(data.SoftwareBranch, 10, 64)
	version, _ := strconv.ParseUint(data.SoftwareVerNum, 10, 64)
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)

	common.WriteLog(adminId, fmt.Sprintf("检查文件新版本   软件分支：%s   软件版本号：%s", data.SoftwareBranch, data.SoftwareVerNum))

	raw := dao.DBCheckECUFileMapRecordExists(uint(branch), uint(version))
	if raw == 0 {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "没有找到对应的记录，请检查参数是否正确",
		})
		return
	}

	var ecuRecord []model.EcuFileMap

	err := dao.DBCheckGtCurrentVersion(uint(branch), uint(version), &ecuRecord)
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
			"NewVerCnt":   len(ecuRecord),
			"NewVersions": ecuRecord,
		},
	})
}
