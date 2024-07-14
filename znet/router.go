package znet

import "myzinx/ziface"

// BaseRouter 实现router时先嵌入BaseRouter基类，然后根据需要对这个基类的方法进行重写
// BaseRouter 中实现的IRouter的方法，不需要写实际的逻辑，只是起到一个将接口实现的逻辑，
// 因为有的Router并不希望实现PreHandle或者PostHandle，那么这样只需要继承这个BaseRouter就可以了
type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request ziface.IRequest) {
}

func (br *BaseRouter) Handle(request ziface.IRequest) {
}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {
}
