package chatmgr

import "github.com/gochat/server/util"

type ChatMgr struct {
	Rooms []*Room
}

type Room struct {
	ID int64
	Name string
	Players []int64
	MsgIndex int64
}

func NewMgr() *ChatMgr {
	return &ChatMgr{}
}

func (mgr *ChatMgr) GetRoom(roomName string) *Room {
	for _, r := range mgr.Rooms {
		if r.Name == roomName {
			return r
		}
	}
	return nil
}

func (mgr *ChatMgr) AddRooms(room *Room) {
	mgr.Rooms = append(mgr.Rooms, room)
}

func (mgr *ChatMgr) JoinRoom(playerID int64, roomName string) *Room {
	for _, r := range mgr.Rooms {
		if r.Name == roomName {
			if util.FindInt64(r.Players, playerID) != -1 {
				return r
			}
			r.AddPlayer(playerID)
			return r
		}
	}
	newRoom := NewRoom(int64(len(mgr.Rooms)+1), roomName)
	mgr.AddRooms(newRoom)
	newRoom.AddPlayer(playerID)
	return newRoom
}

func NewRoom(id int64, name string) *Room{
	return &Room{
		ID: id,
		Name: name,
	}
}

func (r *Room) AddPlayer(playerID int64) {
	r.Players = append(r.Players, playerID)
}

func (r *Room) AddMsg() int64 {
	r.MsgIndex++
	return r.MsgIndex
}