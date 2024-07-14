package znet

import "myzinx/ziface"

// Request 是抽象类IRequest的实现，用于定义一个Request
type Request struct {
	// 已经和客户端建立好的链接
	connection ziface.IConnection
	// 客户端请求的数据
	data []byte
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.connection
}

func (r *Request) GetData() []byte {
	return r.data
}
