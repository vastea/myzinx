package znet

import (
	"fmt"
	"io"
	"myzinx/zconf"
	"myzinx/ziface"
	"myzinx/zpack"
	"net"
	"testing"
	"time"
)

func TestServer_Serve(t *testing.T) {
	go newClient()

	s := NewServer()
	// 注册钩子函数
	s.SetOnConnectionStart(func(connection ziface.IConnection) {
		fmt.Println("[HOOK] OnConnectionStart invoke successfully")
	})
	s.SetOnConnectionStop(func(connection ziface.IConnection) {
		fmt.Println("[HOOK] OnConnectionStop invoke successfully")
	})
	// 注册路由
	s.AddRouter(0, &myRouter{})
	s.Serve()
}

type myRouter struct {
	BaseRouter
}

func (br *myRouter) Handle(request ziface.IRequest) {
	fmt.Println("[PROCESS] Handle is starting, the role is write the same content of the client send to client...")
	err := request.GetConnection().SendMessage(0, []byte("Hello Client, you send content for me is "+string(request.GetData())))
	if err != nil {
		fmt.Println("[ERROR client1] Client1 SendMessage error", err)
	}
}

func newClient() {
	ipAddress := fmt.Sprintf("%s:%d", zconf.Conf.Host, zconf.Conf.Port)

	// 3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	// 模拟一个连接
	conn, err := net.Dial(zconf.Conf.Network, ipAddress)
	if err != nil {
		fmt.Println("[ERROR client1] Client create error", err)
	}
	fmt.Println("[Start client1] Client is created")

	// 监听服务器端传过来的数据
	go func() {
		for {
			// 读取客户端的MessageHead
			buf := make([]byte, zpack.DataPackInstance.GetHeadLen())
			n, err := io.ReadFull(conn, buf)
			if n == 0 {
				fmt.Println("[EOF] The connect read is EOF")
				return
			}
			if err != nil && err != io.EOF {
				if err != nil {
					fmt.Println("[ERROR client1] Client read error", err)
				}
				return
			}
			// 拆包 获取msgLen和msgId
			msg, err := zpack.DataPackInstance.Unpack(buf)
			if err != nil {
				if err != nil {
					fmt.Println("[ERROR client1] Client unpack error", err)
				}
				return
			}
			// 根据msgLen读取MessageData
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				n, err := io.ReadFull(conn, data)
				if n == 0 {
					return
				}
				if err != nil && err != io.EOF {
					if err != nil {
						fmt.Println("[ERROR client1] Client read error", err)
					}
					return
				}
			}
			fmt.Println("收到服务端信息: ", string(data))
		}
	}()

	// 向服务器写数据
	msgData := []byte("Hello")
	msg := &zpack.Message{
		ID:      0,
		DataLen: uint32(len(msgData)),
		Data:    msgData,
	}
	dataBytes, err := zpack.DataPackInstance.Pack(msg)
	if err != nil {
		fmt.Println("[ERROR client1] Client Pack msg error", err)
	}

	// 向服务器写数据
	msgData2 := []byte("World!")
	msg2 := &zpack.Message{
		ID:      0,
		DataLen: uint32(len(msgData2)),
		Data:    msgData2,
	}
	dataBytes2, err := zpack.DataPackInstance.Pack(msg2)
	if err != nil {
		fmt.Println("[ERROR client1] Client Pack msg error", err)
	}

	dataBytes = append(dataBytes, dataBytes2...)
	_, err = conn.Write(dataBytes)
	if err != nil {
		fmt.Println("[ERROR client1] Client write error", err)
	}

	time.Sleep(10 * time.Second)
	conn.Close()
	fmt.Println("[END client1] Client is end")
}
