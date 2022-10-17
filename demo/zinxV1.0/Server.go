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

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("====> DoConnectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnectionBegin")); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Set conn Name, home ....")
	conn.SetProperty("Name", "lxy")
	conn.SetProperty("Home", "fuding")
}

func DoConnectionEnd(conn ziface.IConnection) {
	fmt.Println("===> DoConnectionEnd is Called ...")
	fmt.Println("conn ID = ", conn.GetConnID())

	if value, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name ", value)
	}
	if value, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Home ", value)
	}
}

func main() {
	s := znet.NewServer("demo")
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionEnd)
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
