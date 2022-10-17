package main

import (
	"fmt"
	"zinx/mmo_game/api"
	"zinx/mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	player := core.NewPlayer(conn)
	fmt.Println("Player Pid = ", player.Pid, " is New")
	player.SyncPid()
	player.BroadCastStartPosition()
	core.WorldMgrObj.AddPlayer(player)
	conn.SetProperty("pid", player.Pid)
	player.SyncSurrounding()
	fmt.Println("====> Player Pid = ", player.Pid)
}

func OnConnectionLost(conn ziface.IConnection) {
	pid, _ := conn.GetProperty("pid")
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.Offline()
}

func main() {
	s := znet.NewServer("MMO Game Zinx")
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	s.AddRouter(2, &api.WorldChatApi{})
	s.AddRouter(3, &api.MoveApi{})
	s.Serve()
}
