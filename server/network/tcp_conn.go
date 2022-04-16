package network

import (
	"encoding/json"
	"fmt"
	"github.com/gochat/proto"
	"github.com/gochat/server/player"
	"github.com/gochat/server/redisop"
	"github.com/gochat/server/util"
	"github.com/gochat/server/wordfilter"
	"log"
	"math"
	"net"
)

// AcceptMsgHandle 接受来自客户端的消息
func (cs *ChatServer) AcceptMsgHandle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		log.Println(string(buf[0:length]))
		var msg proto.Message
		err = json.Unmarshal(buf[:length], &msg)
		if err != nil {
			fmt.Println("json.Unmarshl error:", err)
		}
		ack, send := cs.MsgHandler(msg, conn)
		if send {
			ackData, err1 := json.Marshal(ack)
			if err1 != nil {
				fmt.Printf("some error when generate response message, error: %v", err)
			}
			conn.Write(ackData)
		}
	}
}

// MsgHandler 消息处理函数（应该做一个 Register）
func (cs *ChatServer) MsgHandler(msg proto.Message, conn net.Conn) (proto.Message, bool) {
	switch msg.Type {
	case proto.MsgTypeLogin:
		ack := proto.LoginAck{}
		var req *proto.LoginReq
		err := json.Unmarshal([]byte(msg.Data), &req)
		if err != nil {
			log.Println("err", err)
			ack.Err = proto.ErrCodeErrParam
			return proto.PackProtoAck(msg.Type, ack), true
		}
		var p *player.Player
		if req.PlayerID == 0 {
			sameName := cs.PlayerMgr.GetPlayerByName(req.PlayerName)
			if sameName != nil {
				ack.Err = proto.ErrCodeNameRepeated
				return proto.PackProtoAck(msg.Type, ack), true
			}
			p = player.NewPlayer(int64(len(cs.PlayerMgr.Players)+1), req.PlayerName)
			cs.PlayerMgr.AddPlayer(p)
		} else {
			p = cs.PlayerMgr.GetPlayerByName(req.PlayerName)
		}
		if p == nil {
			ack.Err = proto.ErrCodeErrParam
			return proto.PackProtoAck(msg.Type, ack), true
		}
		p.PlayerLogin(conn)
		log.Println("login success")
		ack.PlayerID = p.ID
		ack.PlayerName = p.Name
		return proto.PackProtoAck(msg.Type, ack), true
	case proto.MsgTypeJoinChat:
		ack := proto.JoinChatRoomAck{}
		var req *proto.JoinChatRoomReq
		err := json.Unmarshal([]byte(msg.Data), &req)
		if err != nil {
			log.Println("err", err)
			ack.Err = proto.ErrCodeErrParam
			return proto.PackProtoAck(msg.Type, ack), true
		}
		p := cs.PlayerMgr.GetPlayerByName(req.PlayerName)
		if p == nil {
			ack.Err = proto.ErrCodeErrParam
			return proto.PackProtoAck(msg.Type, ack), true
		}
		addRoom := cs.ChatMgr.JoinRoom(p.ID, req.ChatRoom)
		p.JoinRoom(addRoom.ID)
		ack.RecentMsg = redisop.GetRecentMessage(cs.RedisCli, req.ChatRoom, int64(math.Max(float64(addRoom.MsgIndex-50), 1)), addRoom.MsgIndex)
		ack.ChatRoom = req.ChatRoom
		ack.PlayerName = req.PlayerName
		return proto.PackProtoAck(msg.Type, ack), true
	case proto.MsgTypeSendChatMsg:
		ack := proto.SendChatMsgAck{}
		var req *proto.SendChatMsgReq
		err := json.Unmarshal([]byte(msg.Data), &req)
		if err != nil {
			log.Println("err", err)
			ack.Err = proto.ErrCodeErrParam
			return proto.PackProtoAck(msg.Type, ack), true
		}
		p := cs.PlayerMgr.GetPlayerByName(req.PlayerName)
		if p == nil {
			ack.Err = proto.ErrCodeErrParam
			return proto.PackProtoAck(msg.Type, ack), true
		}
		room := cs.ChatMgr.GetRoom(req.ChatRoom)
		if room == nil {
			ack.Err = proto.ErrCodeErrParam
			return proto.PackProtoAck(msg.Type, ack), true
		}
		checkBadWords := wordfilter.Default.Replace(req.Msg, '*')

		msgID := room.AddMsg()
		// 存入 redis
		redisop.AddMessage(cs.RedisCli, proto.ChatMessage{
			ID:       msgID,
			Sender:   req.PlayerName,
			RoomName: req.ChatRoom,
			Content:  checkBadWords,
			CreateTs: util.NowTs(),
		})

		ntf := proto.ChatMsgNtf{
			Msg:    checkBadWords,
			Sender: p.Name,
		}
		packNtf := proto.PackProtoAck(proto.MsgTypeChatMsg, ntf)
		ntfData, err1 := json.Marshal(packNtf)
		if err1 != nil {
			fmt.Printf("some error when generate response message, error: %v", err)
		}
		// 广播给聊天室
		for _, playerID := range room.Players {
			player := cs.PlayerMgr.GetPlayerByID(playerID)
			if player.Conn != nil {
				player.Conn.Write(ntfData)
			}
		}
		return proto.PackProtoAck(msg.Type, ack), false
	case proto.MsgTypeGMMsg:
		ack := proto.GMAck{}
		var req *proto.GMReq
		err := json.Unmarshal([]byte(msg.Data), &req)
		if err != nil {
			log.Println("err", err)
			ack.Err = proto.ErrCodeErrParam
			return proto.PackProtoAck(msg.Type, ack), true
		}
		switch req.GM {
		case "status":
			// 获取玩家在线时长
			playerName := req.Data
			p := cs.PlayerMgr.GetPlayerByName(playerName)
			log.Printf("player:%v, loginTs:%v, online:%v, roomName:%v",
				playerName, p.LoginTs, util.NowTs()-p.LoginTs, p.JoinRooms)
			log.Println()
			return proto.PackProtoAck(msg.Type, ack), true
		case "popular":
			// 最近十分钟第一高频词
			roomName := req.Data
			room := cs.ChatMgr.GetRoom(roomName)
			word := redisop.GetLast10MinMostPopularWords(cs.RedisCli, roomName, room.MsgIndex, 10)
			log.Println("last 10 min most popular word:", word)
			return proto.PackProtoAck(msg.Type, ack), true
		}
	}
	return proto.Message{}, false
}
