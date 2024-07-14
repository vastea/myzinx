package ziface

// IRouter 路由模块接口，路由可以接受一个指令，然后执行该指令对应的处理方式（不同消息对应不同的处理方式）
type IRouter interface {
	// PreHandle 处理connection业务之前的hook
	PreHandle(request IRequest)
	// Handle 处理connection业务的hook
	Handle(request IRequest)
	// PostHandle 处理connection业务之后的hook
	PostHandle(request IRequest)
}
