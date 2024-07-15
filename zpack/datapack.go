// Package zpack 封包拆包的实现
package zpack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"myzinx/zconf"
	"myzinx/ziface"
	"sync"
)

var packOnce sync.Once
var DataPackInstance ziface.IDataPack

type DataPack struct {
}

// NewDataPack 拆包封包，单例赋值
func NewDataPack() ziface.IDataPack {
	packOnce.Do(func() {
		DataPackInstance = new(DataPack)
	})
	return DataPackInstance
}

// GetHeadLen 获取包的头长度
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 + ID uint32 = 8 bytes
	return 8
}

// Pack 封包
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuff中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// 将MsgID写进dataBuff中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMSgID()); err != nil {
		return nil, err

	}
	// 将data数据写进dataBuff中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err

	}
	return dataBuf.Bytes(), nil
}

// Unpack 拆包
func (dp *DataPack) Unpack(msgBytes []byte) (ziface.IMessage, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuf := bytes.NewBuffer(msgBytes)

	msg := &Message{}
	// 从dataBuf中读取数据到msg.DataLen
	// 因为要修改msg.DataLen的值，所以需要传入指针
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 从dataBuf中读取数据到msg.ID
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	// 数据包是否超长
	if zconf.Conf.MaxPackageSize > 0 && msg.GetDataLen() > zconf.Conf.MaxPackageSize {
		return nil, errors.New("[ERROR] The packageSize is overflow")
	}

	// 这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}

func init() {
	NewDataPack()
}
