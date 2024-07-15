package ziface

// IConnManager 链接管理模块接口
type IConnManager interface {
	// AddConnection 添加链接
	AddConnection(connection IConnection)
	// 删除链接
	RemoveConnection(connection IConnection)
	// GetConnection 根据ConnID获取链接
	GetConnection(connectionID uint32) (IConnection, error)
	// GetConnectionNum 得到当前链接总数
	GetConnectionNum() int
	// ClearConnection 清除并终止所有链接
	ClearConnection()
}
