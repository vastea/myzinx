package znet

import (
	"fmt"
	"github.com/vastea/myzinx/zconf"
	"github.com/vastea/myzinx/ziface"
)

// MsgHandler 是消息处理抽象类IMsgHandler的实现，用于定义消息处理模块
type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter // 存放每个MsgID所对应的处理方法
	TaskQueue      []chan ziface.IRequest    // 负责worker取任务的消息队列
	WorkerPoolSize uint32                    // 业务工作worker池的worker数量，即消息队列的大小
}

// NewMsgHandler 创建MsgHandler
func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: zconf.Conf.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, zconf.Conf.WorkerPoolSize),
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

// StartWorkerPool 启动一个Worker工作池，一个myzinx框架只能有一个worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	// 根据WorkerPoolSize分别开启Worker，每个Worker都分别用goroutine来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 1-给当前worker对应的channel消息队列开辟空间，第n个worker使用第n个channel，每个channel最多接受MaxWorkerTaskLen个任务
		mh.TaskQueue[i] = make(chan ziface.IRequest, zconf.Conf.MaxWorkerTaskLen)
		// 2-启动当前的worker，阻塞等待消息从channel传递进来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandler) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("[START] A worker is start, the workerID is", workerID)

	// 不断阻塞等待对应消息队列的消息
	for request := range taskQueue {
		mh.DoMsgHandler(request)
	}
}

// SendMessageToTaskQueue 将消息放入TaskQueue，由worker处理
func (mh *MsgHandler) SendMessageToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配给不同的worker
	// 根据客户端建立的ConnID进行分配
	workerID := request.GetConnection().GetConnId() % mh.WorkerPoolSize
	// 2 将消息发送给对应worker的TaskQueue
	mh.TaskQueue[workerID] <- request
	fmt.Println("[LOG] Add connID =", request.GetConnection().GetConnId(),
		"request MsgID =", request.GetMsgID(), "to workerID =", workerID)
}
