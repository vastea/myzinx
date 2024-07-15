package ziface

import "net"

// IConnection 链接模块接口
type IConnection interface {
	// Start 启动链接，让当前的链接准备开始工作
	Start()
	// Stop 停止链接，结束当前链接的工作
	Stop()
	// GetConn 获取当前链接所绑定的socket
	GetConn() net.Conn
	// GetConnId 获取当前链接模块的链接ID
	GetConnId() uint32
	// GetRemoteAddr 获取远程客户端的TCP状态 IP Port
	GetRemoteAddr() net.Addr
	// SendMessage 发送数据，将数据发送给远程的客户端
	SendMessage(msgId uint32, data []byte) error
}

// HandleFunc 定义一个处理链接的方法，返回当前的链接、当前链接的内容、内容的长度
type HandleFunc func(net.Conn, []byte, int) error
