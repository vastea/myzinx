package znet

import (
	"fmt"
	"myzinx/ziface"
)

// MsgHandler 是消息处理抽象类IMsgHandler的实现，用于定义消息处理模块
type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个MsgID所对应的处理方法
}

// NewMsgHandler 创建MsgHandler
func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}

}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	router, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("[ERROR] The MsgHandler is not exist, msgID is: ", request.GetMsgID())
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

// AddRouter 为消息添加具体的路由处理
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 如果已经存在，就不用重复加入了
	if _, ok := mh.Apis[msgID]; ok {
		msgErr := fmt.Sprintf("repeated api , msgID = %+v\n", msgID)
		panic(msgErr)
	}
	mh.Apis[msgID] = router
}

func init() {
	NewMsgHandler()
}
