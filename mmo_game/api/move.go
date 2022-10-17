package api

import (
	"fmt"
	"zinx/mmo_game/core"
	"zinx/mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"

	"google.golang.org/protobuf/proto"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal err:", err)
		return
	}
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty Pid err: ", err)
		return
	}

	fmt.Printf("Player Pid %d, move(%f, %f, %f, %f)\n", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
