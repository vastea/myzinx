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
}
