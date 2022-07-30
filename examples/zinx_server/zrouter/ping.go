package zrouter

import (
	"github.com/simondong1979/iTreansfer/network/ziface"
	"github.com/simondong1979/iTreansfer/network/zlog"
	"github.com/simondong1979/iTreansfer/network/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {

	zlog.Debug("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	zlog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping[FromServer]"))
	if err != nil {
		zlog.Error(err)
	}
}
