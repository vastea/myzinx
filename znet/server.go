// Package znet 此包是myzinx对ziface包中抽象接口的具体实现
package znet

import (
	"fmt"
	"myzinx/zconf"
	"myzinx/ziface"
	"net"
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
	// 当前Server的Router
	MsgHandler ziface.IMsgHandler
	// 当前Server的链接管理器
	ConnectionManager ziface.IConnManager
	// 该server创建链接之后自动调用的Hook函数
	OnConnectionStart func(connection ziface.IConnection)
	// 该server销毁链接之前自动调用的Hook函数
	OnConnectionStop func(connection ziface.IConnection)
}

// NewServer 初始化Server
func NewServer() ziface.IServer {
	s := &Server{
		ServerName:        zconf.Conf.Name,
		Network:           zconf.Conf.Network,
		IP:                zconf.Conf.Host,
		Port:              zconf.Conf.Port,
		MsgHandler:        NewMsgHandler(),
		ConnectionManager: NewConnManager(),
	}
	zconf.Conf.Show()
	return s
}

// Start 启动一个Server
func (s *Server) Start() {
	ipAddr := fmt.Sprintf("%s:%d", s.IP, s.Port)

	go func() {
		// 启动server时开启worker工作池
		go func() {
			s.MsgHandler.StartWorkerPool()
		}()

		fmt.Printf("[START] Server %s Listener at ipAddr is starting\n", zconf.Conf.Name)
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
			// 如果超过最大连接数量，那么关闭此连接
			if s.ConnectionManager.GetConnectionNum() >= zconf.Conf.MaxConnection {
				conn.Close()
				continue
			}

			// 与客户端建立连接成功之后，进行业务处理
			// 将conn交给connection去处理
			connection := NewConnection(s, conn, cid, s.MsgHandler)
			go connection.Start()
			cid++
		}
	}()
}

// Stop 该方法用于将一些服务器的资源、状态或一些已经开辟的连接信息 进行停止或关闭
func (s *Server) Stop() {
	fmt.Println("[STOP] Myzinx server is stopping...")
	s.ConnectionManager.ClearConnection()
	fmt.Println("[STOP] All connection were cleared")
}

// Serve 暴露给框架使用者，用于启动服务器，并可以在启动之后做一些自定义的业务逻辑
func (s *Server) Serve() {
	// 启动Server
	s.Start()
	// 阻塞
	select {}
}

// AddRouter 传入一个Router，服务端去建立链接处理消息时按此Router的规则进行处理
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

// GetConnectionManager 返回当前server拥有的链接管理器
func (s *Server) GetConnectionManager() ziface.IConnManager {
	return s.ConnectionManager
}

// SetOnConnectionStart 注册OnConnectionStart钩子函数的方法
func (s *Server) SetOnConnectionStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnectionStart = hookFunc
}

// SetOnConnectionStop 注册OnConnectionStop钩子函数的方法
func (s *Server) SetOnConnectionStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnectionStop = hookFunc
}

// CallOnConnectionStart 调用OnConnectionStart钩子函数的方法
func (s *Server) CallOnConnectionStart(connection ziface.IConnection) {
	if s.OnConnectionStart != nil {
		fmt.Println("[HOOK] OnConnectionStart is invoked")
		s.OnConnectionStart(connection)
	}
}

// CallOnConnectionStop 调用OnConnectionStop钩子函数的方法
func (s *Server) CallOnConnectionStop(connection ziface.IConnection) {
	if s.OnConnectionStop != nil {
		fmt.Println("[HOOK] OnConnectionStop is invoked")
		s.OnConnectionStop(connection)
	}
}
