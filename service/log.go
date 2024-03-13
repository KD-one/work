package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

// ListLogFiles 展示log文件列表
func ListLogFiles(c *gin.Context) {
	dirPath := "./log/uploadRecord"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".log" { // 只显示.log文件
			fileNames = append(fileNames, file.Name())
		}
	}

	c.HTML(http.StatusOK, "logView/log.html", gin.H{"files": fileNames})
}

// ViewLogFile 显示具体log文件内容
func ViewLogFile(c *gin.Context) {
	filename := c.Param("filename")
	dirPath := "./log/uploadRecord"
	fullPath := filepath.Join(dirPath, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", data) // 假设文件是纯文本格式，否则需要调整MIME类型
}
