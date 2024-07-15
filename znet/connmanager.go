package znet

import (
	"errors"
	"fmt"
	"github.com/vastea/myzinx/ziface"
	"sync"
)

type ConnManager struct {
	connectionMap map[uint32]ziface.IConnection // 链接信息集合
	cmLock        sync.RWMutex                  // 读写锁，用于保护链接map
}

// NewConnManager 创建链接管理
func NewConnManager() ziface.IConnManager {
	return &ConnManager{
		connectionMap: make(map[uint32]ziface.IConnection),
	}
}

// AddConnection 添加链接
func (cm *ConnManager) AddConnection(connection ziface.IConnection) {
	// 保护共享资源map
	cm.cmLock.Lock()
	defer cm.cmLock.Unlock()
	// 将conn加入到map中
	cm.connectionMap[connection.GetConnId()] = connection
	fmt.Println("[LOG] Connection add to connectionMap successfully, the connID is", connection.GetConnId())
}

// RemoveConnection 删除链接
func (cm *ConnManager) RemoveConnection(connection ziface.IConnection) {
	// 保护共享资源map
	cm.cmLock.Lock()
	defer cm.cmLock.Unlock()
	// 将conn加入到map中
	delete(cm.connectionMap, connection.GetConnId())
	fmt.Println("[LOG] Connection delete connection from connectionMap successfully, the connID is", connection.GetConnId())
}

// GetConnection 根据ConnID获取链接
func (cm *ConnManager) GetConnection(connectionID uint32) (ziface.IConnection, error) {
	// 保护共享资源map
	cm.cmLock.RLock()
	defer cm.cmLock.RUnlock()
	if connnection, ok := cm.connectionMap[connectionID]; ok {
		return connnection, nil
	}
	return nil, errors.New("connection not found")
}

// GetConnectionNum 得到当前链接总数
func (cm *ConnManager) GetConnectionNum() int {
	return len(cm.connectionMap)
}

// ClearConnection 清除并终止所有链接
func (cm *ConnManager) ClearConnection() {
	// 保护共享资源map
	cm.cmLock.Lock()
	defer cm.cmLock.Unlock()
	// 删除connection并停止connection的工作
	for connID, conn := range cm.connectionMap {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connectionMap, connID)
	}
	fmt.Println("[END] Clear Connection successfully")
}
