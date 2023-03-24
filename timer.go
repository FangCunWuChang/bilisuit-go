package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

var NowTime int64
var startTimer = true

var mutex = sync.Mutex{}

type NetTimer struct {
	Message []byte
	client  *tls.Conn
}

// init
// 初始化
func (receiver *NetTimer) init() *NetTimer {
	var MessageList = []string{
		"GET /x/report/click/now HTTP/1.1\r\nhost: api.bilibili.com", "Connection: keep-alive",
		"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/105.0",
	}
	receiver.Message = []byte(strings.Join(MessageList, "\r\n") + "\r\n\r\n")
	return receiver
}

// updateClient
// 更新连接
func (receiver *NetTimer) updateClient() {
	receiver.client, _ = tls.Dial("tcp", "api.bilibili.com:443", &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "api.bilibili.com",
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
		ClientAuth:         tls.RequireAndVerifyClientCert,
	})
}

// GetBiliTime
// 获取b站时间
func (receiver *NetTimer) getBiliTime() int64 {
	_, _ = receiver.client.Write(receiver.Message)
	var buf = make([]byte, 1024)
	var length, _ = receiver.client.Read(buf)
	var rec = string(buf[:length])
	if len(rec) == 0 {
		receiver.updateClient()
		return receiver.getBiliTime()
	}
	var JsonData = make(map[string]map[string]int64)
	var SplitBody = strings.Split(rec, "\r\n\r\n")
	var Body = SplitBody[len(SplitBody)-1]
	var _ = json.Unmarshal([]byte(Body), &JsonData)
	return JsonData["data"]["now"]
}

// UpdateServerTime
// 更新b站服务器时间
func (receiver *NetTimer) UpdateServerTime() {
	receiver.updateClient()
	for startTimer {
		var biliTime = receiver.getBiliTime()
		mutex.Lock()
		if biliTime > NowTime {
			NowTime = biliTime
		}
		mutex.Unlock()
	}
	_ = receiver.client.Close()
}

// WaitLocalBiliTimer
// 计时器人口
// saleTime: 开售时间
// jump: 跳出时间
func WaitLocalBiliTimer(saleTime, jump int64) {
	var localTime, JumpTime float64
	localTime = float64(time.Now().UnixNano()) / 1e9
	JumpTime = float64(saleTime) - float64(jump)
	for JumpTime > localTime {
		fmt.Printf("\r%f", localTime)
		localTime = float64(time.Now().UnixNano()) / 1e9
	}
}

func WaitServerBiliTimer(saleTime, number int64) {
	for i := 0; i < int(number); i++ {
		var timer = new(NetTimer).init()
		go timer.UpdateServerTime()
	}
	for NowTime < saleTime {
		fmt.Printf("\r%v", NowTime)
	}
	startTimer = false
}
