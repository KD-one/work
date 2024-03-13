package FileRelated

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"test/common"
	"test/model"
	"time"
)

// ToFileUpload 展示文件上传前端页面模板
func ToFileUpload(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "uploadDownloadView/fileUpload.html", nil)
}

// UploadFile 上传单个文件
func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Printf("上传的文件名：%s", file.Filename)
	time_unix := strconv.FormatInt(time.Now().Unix(), 10) // 获取时间戳并转成字符串
	file_path := "upload/" + time_unix + file.Filename    // 设置保存文件的路径，不要忘了后面的文件名
	err := c.SaveUploadedFile(file, file_path)
	if err != nil {
		fmt.Printf("SaveUploadedFile,err=%v", err)
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	c.String(200, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

// UploadFiles 上传多个文件
func UploadFiles(c *gin.Context) {
	var a2lFile string             // a2l文件名
	var buildFile string           // 编译文件名
	var optionalFile []string      // 可选文件
	var a2lFileSize string         // a2l文件大小
	var buildFileSize string       // 编译文件大小
	var optionalFileSize string    // 可选文件大小
	branch := c.PostForm("branch") // 项目号
	number := c.PostForm("number") // 版本号
	name := c.PostForm("name")     // 软件名

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

	//捕获查询数据库时产生的panic
	defer func() {
		if err := recover(); err != nil {
			c.JSON(400, gin.H{"status": "panic", "message": "上传文件有重复，请修改后重新尝试", "data": err})
			return
		}
	}()

	// 查询记录
	m := common.DB.Where(map[string]interface{}{
		"software_version_branch": branch,
		"software_version_number": number,
		"software_name":           name,
	}).First(&model.Filemap{})
	fmt.Printf("查询到%d条记录\n", m.RowsAffected)

	// 查到三个参数相同的记录
	if m.RowsAffected != 0 {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "参数重复，请检查参数后重新上传",
		})
		return
	}

	// 三个参数都不能为空
	if branch == "" || number == "" || name == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "参数不能为空，请填写好参数后重新上传",
		})
		return
	}

	// 验证是否上传了必选文件
	if a2lFile == "" || buildFile == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "缺少必要文件，请重新上传！！！ （提示：缺少a2l或srz文件或hex文件）",
		})
		return
	}

	t := time.Now().Format("2006-01-02 15:04:05")

	// 验证文件表记录并插入到文件表中
	FileTableValidatAndCreate(a2lFile, a2lFileSize, t, branch, number)
	FileTableValidatAndCreate(buildFile, buildFileSize, t, branch, number)
	for _, f := range optionalFile {
		FileTableValidatAndCreate(f, optionalFileSize, t, branch, number)
	}

	// 将可选文件全部用逗号隔开，并拼接成字符串
	of := strings.Join(optionalFile, ",")
	// 创建记录
	common.DB.Model(&model.Filemap{}).Create(map[string]interface{}{
		"created_at":              t,
		"updated_at":              t,
		"software_version_branch": branch,
		"software_version_number": number,
		"software_name":           name,
		"a2l_file":                a2lFile,
		"build_file":              buildFile,
		"optional_file":           of,
	})

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
			c.String(http.StatusBadRequest, "请求失败")
			return
		}

		v, _ := c.Get("userName")
		log.Printf("[info] 用户 %v 上传文件 %s 到项目 %s 版本 %s 软件名 %s", v, file.Filename, branch, number, name)
	}

	// zip
	tmp := "upload/" + branch + "_" + number + "_" + name + ".zip"
	err := ZipFiles(tmp, fileNames)
	if err != nil {
		fmt.Println("上传文件时zip打包出错！！！ ", err)
	}

	// 请求响应
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "上传成功",
	})
}
