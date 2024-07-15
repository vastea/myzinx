package ziface

// IMessage 请求中消息的抽象接口
type IMessage interface {
	// GetMSgID 获取消息ID
	GetMSgID() uint32
	// GetDataLen 获取消息长度
	GetDataLen() uint32
	// GetData 获取消息内容
	GetData() []byte
	// SetMsgID 设置消息ID
	SetMsgID(msgID uint32)
	// SetDataLen 设置消息长度
	SetDataLen(dataLen uint32)
	// SetData 设置消息内容
	SetData(data []byte)
}
