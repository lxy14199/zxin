package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	fmt.Println("recv from client msgId = ", request.GetMsgId(),
		", data = ", string(request.GetData()))

	if err := request.GetConnection().SendMsg(1, []byte("hello zinx")); err != nil {
		fmt.Println("sendMsg err:", err)
	}
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	fmt.Println("recv from client msgId = ", request.GetMsgId(),
		", data = ", string(request.GetData()))

	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...")); err != nil {
		fmt.Println("sendMsg err:", err)
	}
}

func main() {
	s := znet.NewServer("demo")
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}
