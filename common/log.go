package common

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

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
	_, err := f.WriteString(data)
	if err != nil {
		fmt.Println("write file error")
	}
}
