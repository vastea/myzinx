package ziface

// IMsgHandler 消息管理抽象层
type IMsgHandler interface {
	// DoMsgHandler 调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体的路由处理
	AddRouter(msgID uint32, router IRouter)
}
