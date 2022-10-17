package core

import (
	"fmt"
	"math/rand"
	"sync"
	"zinx/mmo_game/pb"
	"zinx/ziface"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	Pid  int32
	Conn ziface.IConnection
	X    float32
	Y    float32
	Z    float32
	V    float32
}

var PidGen int32 = 1
var IdLock sync.Mutex

func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}
	return p
}

func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error!")
	}
	return
}

func (p *Player) SyncPid() {
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	p.SendMsg(1, proto_msg)
}

func (p *Player) BroadCastStartPosition() {
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, proto_msg)
}

func (p *Player) Talk(content string) {
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	players := WorldMgrObj.GetAllPlayers()

	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) SyncSurrounding() {
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Y)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	players_proto_msg := make([]*pb.Player, 0, len(players))

	for _, player := range players {
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}

	SyncPlayer_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg,
	}
	p.SendMsg(202, SyncPlayer_proto_msg)
}

func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {
	p.X, p.Y, p.Z, p.V = x, y, z, v
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	players := p.GetSuroundingPlayer()

	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

}

func (p *Player) GetSuroundingPlayer() []*Player {
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Y)
	players := make([]*Player, 0, len(pids))

	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}

func (p *Player) Offline() {
	players := p.GetSuroundingPlayer()
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)

}
