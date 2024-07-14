package znet

import (
	"fmt"
	"io"
	"myzinx/ziface"
	"net"
)

// Connection 是抽象类IConnection的实现，用于定义一个Connection链接器模块
type Connection struct {
	// 当前的socket连接
	Conn net.Conn
	// 链接的ID
	ConnId uint32
	// 当前的链接状态
	IsOpen bool
	// 当前链接所绑定的处理业务的方法api
	HandleAPI ziface.HandleFunc
	// 告知当前链接已经退出/停止的channel
	ExitChan chan bool
}

// NewConnection 初始化一个Connection
func NewConnection(conn net.Conn, connId uint32, callbackAPI ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnId:    connId,
		IsOpen:    true,
		HandleAPI: callbackAPI,
		ExitChan:  make(chan bool, 1),
	}
}

// StartReader 链接的读
func (c *Connection) StartReader() {
	fmt.Println("[START-", c.ConnId, "Connection] Reader goroutine is starting")
	defer fmt.Println("[STOP-", c.ConnId, "Connection] connection is stopped, the remote addr is", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中
		buf := make([]byte, 4096)
		n, err := c.Conn.Read(buf)
		if n == 0 {
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("[ERROR-", c.ConnId, "Connection] Conn read error:", err)
			return
		}

		// 调用当前链接所绑定的HandleAPI
		if err := c.HandleAPI(c.Conn, buf, len(buf)); err != nil {
			fmt.Println("[ERROR-", c.ConnId, "Connection] HandleAPI handle error:", err)
			return
		}
	}
}

// Start 启动链接，即启动链接要处理的业务
func (c *Connection) Start() {
	// 启动从当前链接的读数据的业务
	go c.StartReader()
	// 启动从当前链接写数据的业务

}

// Stop 关闭链接，主要是关闭服务端与客户端的连接，和Connection中的channel
func (c *Connection) Stop() {
	fmt.Println("[STOP-", c.ConnId, "Connection] Connection is stopping")

	// 如果当前链接已经关闭
	if c.IsOpen == false {
		return
	}
	c.IsOpen = false
	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("[ERROR-", c.ConnId, "Connection] Conn close error:", err)
	}
	// 关闭管道
	close(c.ExitChan)
}

// GetConn 获取socket连接
func (c *Connection) GetConn() net.Conn {
	return c.Conn
}

// GetConnId 获取此链接id
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// GetRemoteAddr 获取socket连接中的远程网络地址
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	panic("not implemented") // TODO: Implement
}
