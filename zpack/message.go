package zpack

// Message 是抽象类IMessage的实现，用于定义在Request中的消息的具体结构
// 是封包和拆包中，一个包的内容
type Message struct {
	ID      uint32 // 消息的ID
	DataLen uint32 // 消息的长度
	Data    []byte // 消息内容
}

// GetMSgID 获取消息ID
func (m *Message) GetMSgID() uint32 {
	return m.ID
}

// GetDataLen 获取消息长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// GetData 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgID 设置消息ID
func (m *Message) SetMsgID(msgID uint32) {

}

// SetDataLen 设置消息长度
func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}

// SetData 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
