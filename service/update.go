package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"strings"
	"test/FileRelated"
	"test/dao"
	"test/model"
	"test/serializer"
	"time"
)

type UpdateService struct {
	ProjectName   string `json:"project_name" form:"project_name"`
	VersionNumber string `json:"version_number" form:"version_number"`
	SoftwareName  string `json:"software_name" form:"software_name"`
	A2lFile       string `json:"a2l_file" form:"a2l_file"`
	BuildFile     string `json:"build_file" form:"build_file"`
	OptionalFile  string `json:"optional_file" form:"optional_file"`
}

func (service *UpdateService) Update(c *gin.Context) serializer.Response {

	var a2lFile string          // a2l文件名
	var buildFile string        // 编译文件名
	var optionalFile []string   // 可选文件
	var a2lFileSize string      // a2l文件大小
	var buildFileSize string    // 编译文件大小
	var optionalFileSize string // 可选文件大小

	form, _ := c.MultipartForm()
	files := form.File["file"] // 获取文件

	// 将文件名获取并赋值
	for _, file := range files {
		if strings.HasSuffix(file.Filename, "a2l") {
			a2lFile = file.Filename
			a2lFileSize = strconv.FormatInt(file.Size, 10)
		} else if strings.HasSuffix(file.Filename, "hex") || strings.HasSuffix(file.Filename, "srz") {
			buildFile = file.Filename
			buildFileSize = strconv.FormatInt(file.Size, 10)
		} else {
			optionalFile = append(optionalFile, file.Filename)
			optionalFileSize = strconv.FormatInt(file.Size, 10)
		}
	}

	// 三个参数都不能为空
	if service.ProjectName == "" || service.VersionNumber == "" || service.SoftwareName == "" {
		return serializer.Response{
			Code: 400,
			Msg:  "参数不能为空，请填写好参数后重新上传",
		}
	}

	// 验证是否上传了必选文件
	if a2lFile == "" || buildFile == "" {
		return serializer.Response{
			Code: 400,
			Msg:  "缺少必要文件，请重新上传！！！ （提示：缺少a2l或srz文件或hex文件）",
		}
	}

	t := time.Now().Format("2006-01-02 15:04:05")

	// 文件数据库中创建新纪录或更新旧记录
	err := dao.CreateOrUpdateFile(a2lFile, a2lFileSize, t, service.ProjectName, service.VersionNumber)
	if err != nil {
		return serializer.Response{
			Code: 500,
			Msg:  "a2l文件名已存在！",
			Data: err,
		}
	}
	err = dao.CreateOrUpdateFile(buildFile, buildFileSize, t, service.ProjectName, service.VersionNumber)
	if err != nil {
		return serializer.Response{
			Code: 500,
			Msg:  "hex或srz文件名已存在！",
			Data: err,
		}
	}
	for _, f := range optionalFile {
		err = dao.CreateOrUpdateFile(f, optionalFileSize, t, service.ProjectName, service.VersionNumber)
		if err != nil {
			return serializer.Response{
				Code: 500,
				Msg:  "可选文件名已存在！",
				Data: err,
			}
		}
	}

	// 获取旧有文件
	filemap, err := dao.GetFileMapByProjectAndVersion(service.ProjectName, service.VersionNumber)
	if err != nil {
		return serializer.Response{
			Code: 500,
			Msg:  "文件更新失败,项目名错误或版本号错误！",
			Data: err,
		}
	}

	oldZipFile := filemap.SoftwareVersionBranch + "_" + filemap.SoftwareVersionNumber + "_" + filemap.SoftwareName + ".zip"
	oldZipFilePath := "upload/" + oldZipFile

	// 删除旧有zip文件
	err = os.Remove(oldZipFilePath)
	if err != nil {
		return serializer.Response{
			Code: 500,
			Msg:  "删除旧有文件失败！请检查文件列表中是否存在该文件！（文件命名：项目名_版本号_软件名.zip）",
			Data: err,
		}
	}

	// 根据命名格式（文件命名：项目名_版本号_软件名.zip）将上传文件添加到upload目录，并打包成zip添加到upload目录
	var fileNames []string
	// 将文件添加到指定目录
	for _, file := range files {
		fmt.Println(file.Filename)
		fileNames = append(fileNames, "upload/"+file.Filename)
		file_path := "upload/" + file.Filename // 设置保存文件的路径

		// 查看file_path是否存在
		_, err := os.Stat(file_path)
		// file_path路径下已经存在此文件，跳过保存
		if err == nil {
			//c.JSON(400, gin.H{"message": "文件已存在", "data": err.Error()})
			//return
			continue
		}
		// file_path路径下不存在此文件
		err = c.SaveUploadedFile(file, file_path) // 保存文件
		if err != nil {
			log.Printf("[error] 上传文件失败：%v", err)
			return serializer.Response{
				Code: 500,
				Msg:  "文件上传失败",
				Data: err,
			}
		}
		v, _ := c.Get("userName")
		log.Printf("[info] 用户 %v 上传文件 %s 到项目 %s 版本 %s 软件名 %s", v, file.Filename, service.ProjectName, service.VersionNumber, service.SoftwareName)

	}

	// zip
	tmp := "upload/" + service.ProjectName + "_" + service.VersionNumber + "_" + service.SoftwareName + ".zip"
	err = FileRelated.ZipFiles(tmp, fileNames)
	if err != nil {
		fmt.Println("上传文件时zip打包出错！！！ ", err)
	}

	of := strings.Join(optionalFile, ",")
	// 将需要修改的字段填入
	fileMap := model.Filemap{
		SoftwareName: service.SoftwareName,
		A2lFile:      a2lFile,
		BuildFile:    buildFile,
		OptionalFile: of,
	}

	if err := dao.UpdateFileMapByProjectAndVersion(&fileMap, service.ProjectName, service.VersionNumber); err != nil {
		return serializer.Response{
			Code: 500,
			Msg:  "文件更新失败",
			Data: err,
		}
	}

	return serializer.Response{
		Code: 200,
		Msg:  "文件更新成功",
	}
}
