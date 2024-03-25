package common

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	F              *os.File
	err            error
	UserRecord     *log.Logger
	DownloadRecord *log.Logger
)

// Log 初始化日志
func Log() {
	uploadFileLog()
	userRecordLog()
	downloadFileLog()
}

// 记录上传文件的信息
func uploadFileLog() {
	t := time.Now().Format("2006_01_02")
	// 初始化日志
	F, err = os.OpenFile("log/uploadRecord/"+t+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	//defer F.Close()
	// 将文件输出流和控制台输出流整合到一个io.Writer上
	multiWriter := io.MultiWriter(os.Stdout, F)
	// 设置日志输出位置
	log.SetOutput(multiWriter)
	// 设置输出内容，除时间外增加打印文件名和行号
	log.SetFlags(log.Ldate | log.Ltime)
}

// 记录用户的登陆注册信息
func userRecordLog() {
	t := time.Now().Format("2006_01_02")
	f, er := os.OpenFile("log/userRecord/"+t+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if er != nil {
		return
	}
	UserRecord = log.New(f, "", log.LstdFlags)
}

// 记录下载文件的信息
func downloadFileLog() {
	t := time.Now().Format("2006_01_02")
	f, er := os.OpenFile("log/downloadRecord/"+t+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if er != nil {
		return
	}
	DownloadRecord = log.New(f, "", log.LstdFlags)
}

func DataToString(userId uint, describe string, parameter ...string) string {
	userIdString := strconv.Itoa(int(userId))
	t := time.Now().Format("2006-01-02 15:04:05")
	data := "[ " + userIdString + " ]   " + t + "   " + describe
	for _, v := range parameter {
		data += "   " + v
	}
	data += "\n"
	return data
}

func WriteStringToLog(data string) error {
	t := time.Now().Format("2006_01_02")
	f, er := os.OpenFile("log/"+t+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if er != nil {
		return errors.New("open file error")
	}
	defer f.Close()
	_, err = f.WriteString(data)
	if err != nil {
		return errors.New("write file error")
	}
	return nil
}

func WriteLog(userId uint, describe string) {
	userIdString := strconv.Itoa(int(userId))
	t := time.Now().Format("2006-01-02 15:04:05")
	data := "[ " + userIdString + " ]   " + t + "   " + describe
	//for _, v := range parameter {
	//	data += "   " + v
	//}
	data += "\n"
	t = time.Now().Format("2006_01_02")
	f, er := os.OpenFile("log/"+t+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if er != nil {
		fmt.Println("open file error")
	}
	defer f.Close()
	_, err = f.WriteString(data)
	if err != nil {
		fmt.Println("write file error")
	}
}
