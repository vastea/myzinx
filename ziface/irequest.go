package ziface

// IRequest 请求模块接口，将客户端请求的数据和对应的链接Connection绑定在一起，形成一个Request
type IRequest interface {
	// GetConnection 得到当前的链接
	GetConnection() IConnection
	// GetData 得到请求的消息数据
	GetData() []byte
}
