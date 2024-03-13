package FileRelated

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// ToDownload 展示所有文件列表
func ToDownload(c *gin.Context) {
	c.HTML(200, "uploadDownloadView/redictWithToken.html", nil)
}

func ShowFileList(c *gin.Context) {
	// 读取upload目录下的所有文件
	files, err := os.ReadDir("upload")
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 构造文件列表
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	c.HTML(http.StatusOK, "uploadDownloadView/download.html", fileNames)
}
