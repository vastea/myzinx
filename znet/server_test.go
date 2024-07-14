package znet

import (
	"fmt"
	"myzinx/ziface"
	"net"
	"testing"
	"time"
)

const (
	ServerName  = "myzinx-v0.1"
	Network     = "tcp"
	IP          = "127.0.0.1"
	Port        = 8888
	SendMessage = "HELLO"
)

func TestServer(t *testing.T) {
	go newClient()

	s := NewServer(
		ServerName,
		Network,
		IP,
		Port)
	s.AddRouter(&myRouter{})
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

	// 监听服务器端传过来的数据
	go func() {
		for {
			msg := make([]byte, 512)
			n, err := conn.Read(msg)
			for p, v := range msg {
				if v == 0 {
					msg = msg[:p]
					break
				}
			}
			if err != nil {
				fmt.Println("[ERROR client1] Client write error", err)
				return
			}
			if n == 0 {
				fmt.Println("[END] Client reader is end")
				return
			}
			fmt.Println("收到服务端信息: ", string(msg))
		}
	}()

	time.Sleep(time.Second)
	// 向服务器写数据
	_, err = conn.Write([]byte(SendMessage + "\n"))
	if err != nil {
		fmt.Println("[ERROR client1] Client write error", err)
	}
}

type myRouter struct {
	BaseRouter
}

func (br *myRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("[PROCESS] PreHandle is starting...")
}

func (br *myRouter) Handle(request ziface.IRequest) {
	fmt.Println("[PROCESS] Handle is starting, the role is write the same content of the client send to client...")
	_, err := request.GetConnection().GetConn().Write([]byte("Hello Client"))
	if err != nil {
		fmt.Println("[ERROR client1] Client write error", err)
	}
}

func (br *myRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("[PROCESS] PostHandle is starting...")
}
