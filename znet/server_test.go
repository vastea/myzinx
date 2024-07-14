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

	// 3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	// 模拟一个连接
	conn, err := net.Dial(Network, ipAddress)
	if err != nil {
		fmt.Println("[ERROR client1] Client create error", err)
	}
	fmt.Println("[Start client1] Client is created")

	// 向服务器写数据
	time.Sleep(time.Second)
	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		fmt.Println("[ERROR client1] Client write error", err)
	}
	fmt.Printf("[client1]向服务器写了%d个字符\n", n)

}
