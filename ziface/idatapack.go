package ziface

// IDataPack 封包、拆包
// 直接面向TCP连接中的数据流，用于处理TCP粘包问题
type IDataPack interface {
	// GetHeadLen 获取包的头长度
	GetHeadLen() uint32
	// Pack 封包
	Pack(msg IMessage) ([]byte, error)
	// Unpack 拆包
	Unpack(msgBytes []byte) (IMessage, error)
}
