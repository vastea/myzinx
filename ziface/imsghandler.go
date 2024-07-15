package ziface

// IMsgHandler 消息管理抽象层
type IMsgHandler interface {
	// DoMsgHandler 调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体的路由处理
	AddRouter(msgID uint32, router IRouter)
	// StartWorkerPool 启动Worker工作池
	StartWorkerPool()
	// SendMessageToTaskQueue 将消息放入TaskQueue，由worker处理
	SendMessageToTaskQueue(request IRequest)
}
