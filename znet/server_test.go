package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

const (
	ServerName = "myzinx-v0.1"
	Network    = "tcp"
	IP         = "127.0.0.1"
	Port       = 8888
)

func TestServer(t *testing.T) {
	go newClient()

	s := NewServer(
		ServerName,
		Network,
		IP,
		Port)
	s.Serve()
}

func newClient() {
	ipAddress := fmt.Sprintf("%s:%d", IP, Port)

	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	fmt.Println("[Start] Client is creating")
	// 模拟一个连接
	conn, err := net.Dial(Network, ipAddress)
	if err != nil {
		fmt.Println("[ERROR] Client create error", err)
	}
	// 向服务器写数据
	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		fmt.Println("[ERROR] Client write error", err)
	}
	fmt.Printf("向服务器写了%d个字符\n", n)

}
