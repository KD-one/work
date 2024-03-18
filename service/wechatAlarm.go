package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Message struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list"`
}

var once sync.Once
var (
	cpuLastWarningTime    *time.Time
	memoryLastWarningTime *time.Time
	netLastWarningTime    time.Time
)

// SendWeChatAlarm 发送企业微信告警
func SendWeChatAlarm(url, content string, atUsers ...string) error {
	// 如果需要@所有人，需要在atUsers中添加 "@all",
	// 需要@某个人直接添加其企业微信userid即可(一般为名字的英文拼写,全部小写)
	var message Message

	param := make(map[string]interface{})
	param["msgtype"] = "text"

	message.Content = content
	message.MentionedList = atUsers

	param["text"] = message

	postData, er := json.Marshal(param)
	if er != nil {
		return errors.New("json marshal出错")
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(postData)))
	if err != nil {
		return errors.New("发送消息出错")
	}
	defer resp.Body.Close()

	// 如果需要，可以进一步检查响应内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("读取响应内容出错")
	}
	fmt.Println("发送消息后的响应: ", string(bodyBytes))
	return nil
}

// CompareMemoryUsedPercent 检查内存使用率
// 参数: limit: 阈值, limitDuration: 警告间隔
func CompareMemoryUsedPercent(limit float64, limitDuration time.Duration) error {

	v, err := mem.VirtualMemory()
	if err != nil {
		return errors.New(fmt.Sprint("获取内存使用率出错: ", err.Error()))
	}
	err = ValidAlert(viper.GetString("robot.monitor"), "内存使用率", limit, v.UsedPercent, limitDuration, memoryLastWarningTime)
	if err != nil {
		return err
	}
	return nil
}

// CompareCpuUsedPercent 检查cpu使用率
// 参数: limit: 阈值, limitDuration: 警告间隔
func CompareCpuUsedPercent(limit float64, limitDuration time.Duration) error {

	percent, err := cpu.Percent(0, false)
	if err != nil {
		return errors.New(fmt.Sprint("获取cpu使用率出错: ", err.Error()))
	}
	err = ValidAlert(viper.GetString("robot.monitor"), "cpu使用率", limit, percent[0], limitDuration, cpuLastWarningTime)
	if err != nil {
		return err
	}
	return nil
}

// CompareNetUsedPercent 检查网络带宽使用率
// 参数: limit: 阈值, limitDuration: 警告间隔
func CompareNetUsedPercent(t time.Duration) error {
	// 获取所有网络接口的流量信息
	interfaces, err := net.IOCounters(false)
	if err != nil {
		return errors.New(fmt.Sprint("获取网络接口流量信息出错: ", err.Error()))
	}

	// 遍历所有接口
	for _, iface := range interfaces {
		// 提取出接口名称和流量信息
		name := iface.Name
		receivedBytes := iface.BytesRecv
		transmittedBytes := iface.BytesSent

		// 输出相关信息
		fmt.Printf("Interface: %s\nReceived Bytes: %d\nTransmitted Bytes: %d\n",
			name, receivedBytes, transmittedBytes)
	}
	return nil
}

func NotFound(c *gin.Context) {
	c.HTML(404, "errView/404.html", nil)
}

func ValidAlert(robotName, alert string, limit, validData float64, limitDuration time.Duration, lastWarningTime *time.Time) error {
	now := time.Now()
	// 内存使用率超过报警线
	//fmt.Println("--------", *lastWarningTime, "---------")
	if validData > limit {
		// 第一次触发警告,重复计时
		if lastWarningTime == nil {
			lastWarningTime = &now
			err := SendWeChatAlarm(robotName, fmt.Sprintf("%s高于阈值: %.2f", alert, limit))
			if err != nil {
				return err
			}
			// 距离上次警告已经超过30分钟，重复计时
		} else if now.Sub(*lastWarningTime) >= limitDuration {
			fmt.Println("--------", *lastWarningTime, "---------")
			*lastWarningTime = now
			fmt.Println("--------", *lastWarningTime, "---------")
			err := SendWeChatAlarm(robotName, fmt.Sprintf("%s持续30分钟高于阈值: %.2f", alert, limit))
			if err != nil {
				return err
			}
		}
		// 内存使用率在limit以下
	} else if lastWarningTime != nil && now.Sub(*lastWarningTime) < limitDuration {
		// 内存使用率从大于limit降到了limit以下，重置最后警告时间
		lastWarningTime = nil
	}
	return nil
}

func CheckSystemStatus() {
	for {
		_ = CompareMemoryUsedPercent(60, 10*time.Minute)
		_ = CompareCpuUsedPercent(60, 10*time.Minute)
		time.Sleep(time.Second * 1)
	}
}
