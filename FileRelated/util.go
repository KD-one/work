package FileRelated

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"test/common"
	"test/model"
)

// DownloadFileFromParam 根据路径参数不同，下载不同文件
func DownloadFileFromParam(c *gin.Context) {
	fileName := c.Param("file")
	filePath := filepath.Join("upload", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		c.String(http.StatusNotFound, "File Not Found")
		return
	}
	defer file.Close()

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")

	// 将文件内容写入响应
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	common.DownloadRecord.Println("[info] 下载文件：", fileName, "成功")
}

// ZipFiles 将多个文件打包成一个zip文件
func ZipFiles(filename string, files []string) error {
	fmt.Println("start zip file......")
	//创建输出文件目录
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	//创建空的zip档案，可以理解为打开zip文件，准备写入
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	fmt.Println("zip file created success!")
	return nil
}

// AddFileToZip 添加单个文件到zip文件中
func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	//打开要压缩的文件
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	//获取文件的描述
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	//FileInfoHeader返回一个根据fi填写了部分字段的Header，可以理解成是将fileInfo转换成zip格式的文件信息
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filename
	/*
	   预定义压缩算法。
	   archive/zip包中预定义的有两种压缩方式。一个是仅把文件写入到zip中。不做压缩。一种是压缩文件然后写入到zip中。默认的Store模式。就是只保存不压缩的模式。
	   Store   unit16 = 0  //仅存储文件
	   Deflate unit16 = 8  //压缩文件
	*/
	header.Method = zip.Deflate
	//创建压缩包头部信息
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	//将源复制到目标，将fileToZip 写入writer   是按默认的缓冲区32k循环操作的，不会将内容一次性全写入内存中,这样就能解决大文件的问题
	_, err = io.Copy(writer, fileToZip)
	return err
}

// FileTableValidatAndCreate 验证文件表记录并插入到文件数据库中
func FileTableValidatAndCreate(filename, size, t, projectNumber, versionNumber string) {
	if filename != "" {
		tmp := strings.Split(filename, ".")
		filename = tmp[0]
		fileType := tmp[1]
		f3 := common.DB.Where(map[string]interface{}{
			"name": filename,
			"type": fileType,
		}).First(&model.File{})

		// 数据库中没有查到此纪录，则插入
		if f3.RowsAffected == 0 {
			// 创建文件记录
			common.DB.Model(&model.File{}).Create(map[string]interface{}{
				"name":           filename,
				"type":           fileType,
				"size":           size,
				"created_at":     t,
				"project_number": projectNumber,
				"version_number": versionNumber,
			})
			//log.Printf("上传的文件名：%s", filename)
		}
	}
}

// InitPageInfo 初始化分页相关信息
func InitPageInfo(c *gin.Context) {
	//p := c.Query("page")
	//page, _ := strconv.Atoi(p)
	//pageSize := 15
	//start, end := SlicePage(page, pageSize, len(fileNames))
	//totalPages := int(math.Ceil(float64(len(fileNames)) / float64(pageSize)))
	//
	//// 计算并检查下一页和末页是否有效
	//nextLink := ""
	//if page+1 <= totalPages {
	//	nextLink = "/download?page=" + strconv.Itoa(page+1)
	//} else {
	//	nextLink = ""
	//}
	//lastPageLink := ""
	//if page < totalPages {
	//	lastPageLink = "/download?page=" + strconv.Itoa(totalPages)
	//}
	//
	//c.HTML(http.StatusOK, "uploadDownloadView/download.html", gin.H{
	//	"files":        fileNames[start:end],
	//	"page":         page,
	//	"size":         pageSize,
	//	"totalPages":   totalPages,
	//	"nums":         make([]int, end-start),
	//	"nextLink":     nextLink,
	//	"lastPageLink": lastPageLink,
	//})
}

// SlicePage 根据当前页和每页大小，切分文件列表
func SlicePage(page, pageSize, nums int) (sliceStart, sliceEnd int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 20
	}

	if pageSize > nums {
		return 0, nums
	}

	// 总页数
	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}
