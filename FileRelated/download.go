package FileRelated

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"test/common"
	"test/model"
)

func ParamToDownload(c *gin.Context) {
	c.HTML(http.StatusOK, "uploadDownloadView/paramDownload.html", nil)
}

func DownloadFile(c *gin.Context) {
	//  查询一些必要的参数 进行一些必要的验证
	//attachmentId := c.Query("attachment_id")
	attachmentName := c.Query("attachment_name")

	var data []byte
	// 获取要返回的文件数据流
	// 看你文件存在哪里了，本地就直接os.Open就可以了，总之是要获取一个[]byte
	fileContent, err := os.Open("./upload")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "下载文件失败,请联系管理员"})
		return
	}

	fileContent.Read(data)
	// 设置返回头并返回数据
	fileContentDisposition := "attachment;filename=\"" + attachmentName + "\""
	c.Header("Content-Type", "application/zip") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "contentType", data)
}

// UrlFileDownloadService 根据url参数，寻找文件并下载
func UrlFileDownloadService(c *gin.Context) {
	filePath := c.Query("url")
	//打开文件
	fileTmp, errByOpenFile := os.Open(filePath)
	defer fileTmp.Close()

	//获取文件的名称
	fileName := path.Base(filePath)

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	if filePath == "" || fileName == "" || errByOpenFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "下载出错", "data": errByOpenFile.Error()})
		//c.Redirect(http.StatusFound, "/404")
		log.Println(" [error] 下载失败:  ", errByOpenFile)
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")

	c.File(filePath)
	log.Println(" [info] 成功下载：  ", fileName)
	return
}

// BuildFileDownloadService 根据参数，寻找.srz文件或.hex文件下载
func BuildFileDownloadService(c *gin.Context) {
	// 解析出请求中的项目号、软件名、软件版本号、type参数
	SoftwareVersionBranch := c.Query("SoftwareVersionBranch")
	SoftwareVersionNumber := c.Query("SoftwareVersionNumber")
	name := c.Query("name")
	//t := c.Query("type")
	var f model.Filemap
	f.SoftwareVersionBranch = SoftwareVersionBranch
	f.SoftwareVersionNumber = SoftwareVersionNumber
	f.SoftwareName = name

	// 捕获查询数据库时产生的panic
	defer func() {
		if err := recover(); err != nil {
			c.JSON(400, gin.H{"status": "panic", "message": "没找到要下载的记录，请检查参数并重新下载"})
			return
		}
	}()
	// 查询数据库中是否存在，没找到记录会产生panic
	err := common.DB.Where(map[string]interface{}{"software_version_branch": SoftwareVersionBranch, "software_version_number": SoftwareVersionNumber, "software_name": name}).Find(&f).Error
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}

	//定义文件路径
	filePath := "upload/" + f.BuildFile
	//打开文件
	fileTmp, errByOpenFile := os.Open(filePath)
	defer fileTmp.Close()

	//获取文件的名称
	fileName := path.Base(filePath)

	if filePath == "" || fileName == "" || errByOpenFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "下载出错", "data": errByOpenFile.Error()})
		//c.Redirect(http.StatusFound, "/404")
		log.Println(" [error] 下载失败:  ", errByOpenFile)
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.File(filePath)
	log.Println(" [info] 成功下载：  ", fileName)
	return
}

// A2lFileDownloadService 根据参数，寻找.a2l文件下载
func A2lFileDownloadService(c *gin.Context) {
	// 解析出请求中的项目号、软件名、软件版本号、type参数
	SoftwareVersionBranch := c.Query("SoftwareVersionBranch")
	SoftwareVersionNumber := c.Query("SoftwareVersionNumber")
	name := c.Query("name")
	//t := c.Query("type")
	var f model.Filemap
	f.SoftwareVersionBranch = SoftwareVersionBranch
	f.SoftwareVersionNumber = SoftwareVersionNumber
	f.SoftwareName = name

	// 捕获查询数据库时产生的panic
	defer func() {
		if err := recover(); err != nil {
			c.JSON(400, gin.H{"status": "panic", "message": "没找到要下载的记录，请检查参数并重新下载"})
			return
		}
	}()
	// 查询数据库中是否存在，没找到记录会产生panic
	err := common.DB.Where(map[string]interface{}{"software_version_branch": SoftwareVersionBranch, "software_version_number": SoftwareVersionNumber, "software_name": name}).Find(&f).Error
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}

	//定义文件路径
	filePath := "upload/" + f.A2lFile
	//打开文件
	fileTmp, errByOpenFile := os.Open(filePath)
	defer fileTmp.Close()

	//获取文件的名称
	fileName := path.Base(filePath)

	if filePath == "" || fileName == "" || errByOpenFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "下载出错", "data": errByOpenFile.Error()})
		//c.Redirect(http.StatusFound, "/404")
		log.Println(" [error] 下载失败:  ", errByOpenFile)
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.File(filePath)
	log.Println(" [info] 成功下载：  ", fileName)
	return
}

// FileDownloadService 根据请求参数，寻找文件并下载
func FileDownloadService(c *gin.Context) {
	// 解析出请求中的项目号、软件名、软件版本号、type参数
	SoftwareVersionBranch := c.Query("SoftwareVersionBranch")
	SoftwareVersionNumber := c.Query("SoftwareVersionNumber")
	name := c.Query("name")

	var f model.Filemap
	f.SoftwareVersionBranch = SoftwareVersionBranch
	f.SoftwareVersionNumber = SoftwareVersionNumber
	f.SoftwareName = name

	// 捕获查询数据库时产生的panic
	defer func() {
		if err := recover(); err != nil {
			c.JSON(400, gin.H{"status": "panic", "message": "参数输入有误或者当前记录缺少必备文件，请检查参数并重新下载", "data": err})
			return
		}
	}()
	// 查询数据库中是否存在，没找到记录会产生panic
	err := common.DB.Where(map[string]interface{}{"software_version_branch": SoftwareVersionBranch, "software_version_number": SoftwareVersionNumber, "software_name": name}).Find(&f).Error
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}

	//定义文件路径
	buildFilePath := "upload/" + f.BuildFile
	a2lFilePath := "upload/" + f.A2lFile
	var optionalFilePath []string
	if f.OptionalFile != "" {
		optionalFilePath = strings.Split(f.OptionalFile, ",")
	}

	//打开文件
	fileTmp, errByOpenFile := os.Open(buildFilePath)
	defer fileTmp.Close()
	if errByOpenFile != nil {
		c.JSON(400, gin.H{"status": "error", "message": errByOpenFile.Error()})
		return
	}

	//获取文件的名称
	//buildFileName := path.Base(buildFilePath)
	//a2lFileName := path.Base(a2lFilePath)

	// 验证
	if f.BuildFile == "" || f.A2lFile == "" || errByOpenFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "下载出错", "data": errByOpenFile.Error()})
		//c.Redirect(http.StatusFound, "/404")
		common.DownloadRecord.Println(" [error] 下载失败:  ", errByOpenFile)
		return
	}

	//要压缩成一个zip的多个文件的路径
	files := []string{buildFilePath, a2lFilePath}
	if optionalFilePath != nil {
		for _, ofp := range optionalFilePath {
			files = append(files, "upload/"+ofp)
		}
	}

	//now := strconv.FormatInt(time.Now().UnixNano(), 10)
	//设置输出的zip的路径
	output := "upload/" + SoftwareVersionBranch + "_" + SoftwareVersionNumber + "_" + name + ".zip"

	filename := path.Base(output)

	// 判断upload文件夹中是否存在指定的zip文件
	_, err = os.Stat(output)
	if err != nil {
		// 文件不存在，则将文件压缩打包
		if err := ZipFiles(output, files); err != nil {
			panic(err)
		}
	}

	// 设置请求头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	// 下载文件
	c.File(output)
	c.String(http.StatusOK, "下载成功")
	common.DownloadRecord.Println(" [info] 成功下载：  ", filename)
}

// VisitorFileList 展示游客可以查看的文件列表
func VisitorFileList(c *gin.Context) {
	// 读取upload目录下的所有文件
	files, err := os.ReadDir("upload")
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 构造文件列表
	var fileNames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".hex") || strings.HasSuffix(file.Name(), ".srz") || strings.HasSuffix(file.Name(), ".a2l") {
			fileNames = append(fileNames, file.Name())
		}
	}

	c.HTML(http.StatusOK, "uploadDownloadView/download.html", fileNames)
}
