// Package znet 此包是myzinx对ziface包中抽象接口的具体实现
package znet

import (
	"errors"
	"fmt"
	"io"
	"myzinx/ziface"
	"net"
	"os"
)

// Server 是抽象类IServer的实现，用于定义一个Server服务器模块
type Server struct {
	// ServerName 服务器的名称
	ServerName string
	// IPVersion 服务器绑定的ip版本("tcp4"...)
	Network string
	// IP 服务器绑定的ip地址
	IP string
	// Port 服务器绑定的端口
	Port int
}

// NewServer 初始化Server
func NewServer(name, network, ip string, port int) ziface.IServer {
	return &Server{
		ServerName: name,
		Network:    network,
		IP:         ip,
		Port:       port,
	}
}

// CallBackToClient 定义当前客户端连接所绑定的HandleAPI
// 回写函数，是ziface.iconnection.HandleFunc的实现
func CallBackToClient(conn net.Conn, data []byte, dn int) error {
	n, err := io.Copy(os.Stdout, conn)
	if n == 0 {
		return nil
	}
	if err != nil && err != io.EOF {
		fmt.Println("[ERROR] Bytes copy error:", err)
		return errors.New("[Error] The implement of HandleFunc CallBackToClient error")
	}
	return nil
}

// Start 启动一个Server
func (s *Server) Start() {
	ipAddr := fmt.Sprintf("%s:%d", s.IP, s.Port)

	go func() {
		fmt.Printf("[START] Server Listener at ipAddr is %s starting\n", ipAddr)
		listener, err := net.Listen(s.Network, ipAddr)
		if err != nil {
			fmt.Printf("[ERROR] net.Listen error, network is %s, ip is %s, port is %d, error is %v\n",
				"tcp", s.IP, s.Port, err)
			return
		}

		fmt.Println("[START] Start myzinx server is successful, the serverName is: ", s.ServerName)

		var cid uint32 = 0 // 定义connection的ConnId
		// 接收监听
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("[ERROR] Accept listener error:", err)
				return
			}
			fmt.Println("与客户端连接成功", conn.RemoteAddr())
			// 与客户端建立连接成功之后，进行业务处理
			// 将conn交给connection去处理
			connection := NewConnection(conn, cid, CallBackToClient)
			go connection.Start()
			cid++
		}
	}()
}

// Stop 该方法用于将一些服务器的资源、状态或一些已经开辟的连接信息 进行停止或关闭
func (s *Server) Stop() {

}

// Serve 暴露给框架使用者，用于启动服务器，并可以在启动之后做一些自定义的业务逻辑
func (s *Server) Serve() {
	// 启动Server
	s.Start()
	// 阻塞
	select {}
}
