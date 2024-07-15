package ziface

// IRequest 请求模块接口，将客户端请求的数据和对应的链接Connection绑定在一起，形成一个Request
type IRequest interface {
	// GetConnection 得到当前的连接
	GetConnection() IConnection
	// GetMessage 得到当前链接对应的message
	GetMessage() IMessage
	// GetData 得到当前链接对应的message中的具体数据
	GetData() []byte
	// GetMsgID 得到当前链接对应的message的ID
	GetMsgID() uint32
}
