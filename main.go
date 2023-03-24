package main

import (
	"crypto/tls"
	"fmt"
	"time"
)

// BuildMessage
// 生成报文
func BuildMessage(headers map[string]string, formData string) ([]byte, []byte) {
	var message = "POST /xlive/revenue/v2/order/createOrder HTTP/1.1\r\n"
	for s := range headers {
		message += fmt.Sprintf("%v: %v\r\n", s, headers[s])
	}
	var MessageByte = []byte(message + "\r\n" + formData)
	return MessageByte[:len(MessageByte)-1], MessageByte[len(MessageByte)-1:]
}

// H1CreateTlsConnection
// 创建连接
func H1CreateTlsConnection(BuyHost string) *tls.Conn {
	var adder = fmt.Sprintf("%v:443", BuyHost)
	var client, _ = tls.Dial("tcp", adder, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         BuyHost,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
		ClientAuth:         tls.RequireAndVerifyClientCert,
	})
	return client
}

// H1SendMessage
// 发送请求
func H1SendMessage(client *tls.Conn, body []byte) {
	_, _ = client.Write(body)
}

// H1ReceiveResponse
// 接收响应
func H1ReceiveResponse(client *tls.Conn, BufLen int64) []byte {
	var result = make([]byte, BufLen)
	var length, _ = client.Read(result)
	return result[:length]
}

func main() {
	fmt.Printf("%v\n", "http1_socket_golang")
	var filePath = GetSettingFilePath()
	var headers, startTime, delayTime, formData = ReaderSetting(filePath)
	var SleepTimeNumber = (float64(delayTime) / 1000) * float64(time.Second)

	var MessageHeader, MessageBody = BuildMessage(headers, formData)

	WaitLocalBiliTimer(startTime, 3)

	var client = H1CreateTlsConnection(headers["host"])
	H1SendMessage(client, MessageHeader)

	WaitServerBiliTimer(startTime, 4)

	time.Sleep(time.Duration(SleepTimeNumber))

	var s = time.Now().UnixNano() / 1e6

	H1SendMessage(client, MessageBody)
	var res = H1ReceiveResponse(client, 1024)

	var e = time.Now().UnixNano() / 1e6

	_ = client.Close()

	fmt.Printf("\n%v\n", string(res))
	fmt.Printf("耗时%vms\n", e-s)
}
