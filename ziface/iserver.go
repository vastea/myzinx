// Package ziface 此包是myzinx框架的抽象接口层
package ziface

// IServer 服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 给当前的服务注册一个路由，供客户端的连接处理使用
	AddRouter(msgID uint32, router IRouter)
	// GetConnectionManager 返回当前server拥有的链接管理器
	GetConnectionManager() IConnManager
	// SetOnConnectionStart 注册OnConnectionStart钩子函数的方法
	SetOnConnectionStart(hookFunc func(connection IConnection))
	// SetOnConnectionStop 注册OnConnectionStop钩子函数的方法
	SetOnConnectionStop(hookFunc func(connection IConnection))
	// CallOnConnectionStart 调用OnConnectionStart钩子函数的方法
	CallOnConnectionStart(connection IConnection)
	// CallOnConnectionStop 调用OnConnectionStop钩子函数的方法
	CallOnConnectionStop(connection IConnection)
}
