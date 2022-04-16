package player

import (
	"github.com/gochat/server/util"
	"net"
)

type Mgr struct {
	Players []*Player
	Online  []int64 // 在线玩家
}

func NewMgr() *Mgr {
	return &Mgr{}
}

func (mgr *Mgr) AddPlayer(p *Player) {
	mgr.Players = append(mgr.Players, p)
}

func (mgr *Mgr) GetPlayerByID(id int64) *Player {
	for _, p := range mgr.Players {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func (mgr *Mgr) GetPlayerByName(name string) *Player {
	for _, p := range mgr.Players {
		if p.Name == name {
			return p
		}
	}
	return nil
}

type Player struct {
	ID        int64
	Name      string
	JoinRooms []int64
	Conn      net.Conn
	LoginTs   int64
}

func NewPlayer(id int64, name string) *Player {
	return &Player{
		ID:   id,
		Name: name,
	}
}

func (p *Player) PlayerLogin(conn net.Conn) {
	p.Conn = conn
	p.LoginTs = util.NowTs()
}

func (p *Player) JoinRoom(roomID int64) {
	if util.FindInt64(p.JoinRooms, roomID) != -1 {
		return
	}
	p.JoinRooms = append(p.JoinRooms, roomID)
}

func (p *Player) ExitRoom(roomID int64) {
	for idx, id := range p.JoinRooms {
		if id == roomID {
			p.JoinRooms = append(p.JoinRooms[:idx], p.JoinRooms[idx+1:]...)
			return
		}
	}
}
