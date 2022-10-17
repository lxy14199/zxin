package api

import (
	"fmt"
	"zinx/mmo_game/core"
	"zinx/mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"

	"github.com/golang/protobuf/proto"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Talk Handler err: ", err)
		return
	}
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("Get PlayerId err:", err)
		return
	}

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	player.Talk(proto_msg.Content)
}
