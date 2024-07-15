package znet

import (
	"errors"
	"fmt"
	"io"
	"myzinx/zconf"
	"myzinx/ziface"
	"myzinx/zpack"
	"net"
)

// Connection 是抽象类IConnection的实现，用于定义一个Connection链接器模块
type Connection struct {
	// 当前connection隶属于哪个Server
	Server ziface.IServer
	// 当前的socket连接
	Conn net.Conn
	// 链接的ID
	ConnId uint32
	// 当前的链接状态
	IsOpen bool
	// 当前Connection对应的MsgHandler
	MsgHandler ziface.IMsgHandler
	// 无缓冲管道，用于读写goroutine之间的通信
	MsgChan chan []byte
}

// NewConnection 初始化一个Connection
func NewConnection(server ziface.IServer, conn net.Conn, connId uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Server:     server,
		Conn:       conn,
		ConnId:     connId,
		IsOpen:     true,
		MsgHandler: msgHandler,
		MsgChan:    make(chan []byte),
	}

	// 将connection加入到ConnectionManager中
	c.Server.GetConnectionManager().AddConnection(c)

	return c
}

// StartReader 链接的读
func (c *Connection) StartReader() {
	fmt.Println("[START-", c.ConnId, "Connection] Reader goroutine is starting")
	defer fmt.Println("[STOP-", c.ConnId, "Connection] connection is stopped, the remote addr is", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的MessageHead
		buf := make([]byte, zpack.DataPackInstance.GetHeadLen())
		n, err := io.ReadFull(c.Conn, buf)
		if n == 0 {
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("[ERROR-", c.ConnId, "Connection] Conn read error:", err)
			return
		}
		// 拆包 获取msgLen和msgId
		msg, err := zpack.DataPackInstance.Unpack(buf)
		if err != nil {
			fmt.Println("[ERROR] The server unpack error:", err)
			return
		}
		// 根据msgLen读取MessageData
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			n, err := io.ReadFull(c.Conn, data)
			if n == 0 {
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("[ERROR-", c.ConnId, "Connection] Conn read error:", err)
				return
			}
		}
		msg.SetData(data)

		// 得到当前Connection数据对应的Request数据
		req := &Request{
			connection: c,
			msg:        msg,
		}

		if zconf.Conf.WorkerPoolSize > 0 {
			fmt.Println("[START] WorkerPool Mode is start...")
			// 已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMessageToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

// StartWriter 链接的写
func (c *Connection) StartWriter() {
	fmt.Println("[START-", c.ConnId, "Connection] Writer goroutine is starting")
	for msgBytes := range c.MsgChan {
		// 将数据发送给客户端
		if _, err := c.Conn.Write(msgBytes); err != nil {
			fmt.Println("[ERROR] The Connection write dataBytes error")
		}
	}
}

// Start 启动链接，即启动链接要处理的业务
func (c *Connection) Start() {
	// 启动从当前链接的读数据的业务
	go c.StartReader()
	// 启动从当前链接写数据的业务
	go c.StartWriter()
	// 按照开发者传递进来的创建链接之后需要调用的处理业务，执行对应的hook函数
	c.Server.CallOnConnectionStart(c)
}

// Stop 关闭链接，主要是关闭服务端与客户端的连接，和Connection中的channel
func (c *Connection) Stop() {
	fmt.Println("[STOP-", c.ConnId, "Connection] Connection is stopping")
	c.Server.CallOnConnectionStop(c)

	// 如果当前链接已经关闭
	if c.IsOpen == false {
		return
	}
	c.IsOpen = false

	// 将当前链接从链接管理器中删除
	c.Server.GetConnectionManager().RemoveConnection(c)

	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("[ERROR-", c.ConnId, "Connection] Conn close error:", err)
	}
	// 关闭管道
	close(c.MsgChan)
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

// SendMessage 将服务端发给客户端的数据先封包再发送
func (c *Connection) SendMessage(msgId uint32, data []byte) error {
	// 判断连接状态
	if c.IsOpen == false {
		return errors.New("[ERROR] Connection already closed")
	}
	// 对数据封包
	dp := zpack.DataPackInstance
	msg := &zpack.Message{
		ID:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
	dataBytes, err := dp.Pack(msg)
	if err != nil {
		return errors.New("[ERROR] The Connection Pack data error")
	}
	// 将数据发送给客户端
	c.MsgChan <- dataBytes
	return nil
}
