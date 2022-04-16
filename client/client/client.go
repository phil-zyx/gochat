package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gochat/proto"
	"log"
	"net"
	"os"
)

type Client struct {
	ID         int64
	PlayerName string
	ChatRoom   string
	conn       net.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Dial(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Dial tcp err:%v", err)
		return
	}
	c.conn = conn
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SendMsg(msg []byte) {
	//data, err := json.Marshal(msg)
	//if err != nil {
	//	log.Printf("json marshal error: %v", err)
	//	return
	//}
	c.conn.Write(msg)
}

func (c *Client) ReadCmdArgs() {
	buf := make([]byte, 1024)
	for {
		length, err := c.conn.Read(buf)
		if err != nil {
			log.Println(err)
			c.conn.Close()
			return
		}
		log.Println(string(buf[0:length]))

		var msg proto.Message
		err = json.Unmarshal(buf[:length], &msg)
		if err != nil {
			fmt.Printf("json.Unmarshl error: %v", err)
		}
		switch msg.Type {
		case proto.MsgTypeLogin:
			var ack *proto.LoginAck
			err1 := json.Unmarshal([]byte(msg.Data), &ack)
			if err1 == nil {
				c.PlayerName = ack.PlayerName
				c.ID = ack.PlayerID
			}
		case proto.MsgTypeJoinChat:
			var ack *proto.JoinChatRoomAck
			err1 := json.Unmarshal([]byte(msg.Data), &ack)
			if err1 == nil {
				c.ChatRoom = ack.ChatRoom
			}
		case proto.MsgTypeSendChatMsg:
			var ack *proto.SendChatMsgAck
			err1 := json.Unmarshal([]byte(msg.Data), &ack)
			if err1 == nil {
				log.Println("ack:", ack)
			}
		case proto.MsgTypeChatMsg:
			var ntf *proto.ChatMsgNtf
			err1 := json.Unmarshal([]byte(msg.Data), &ntf)
			if err1 == nil {
				log.Printf("Player: %v send a Msg: %v", ntf.Sender, ntf.Msg)
			}
		}
	}
}

// ShowHandleAfterLogin .
func ShowHandleAfterLogin(c *Client) {
	log.Println("----------------login succeed!----------------")
	log.Println("select what you want to do")
	log.Println("1. Join chat room")
	log.Println("2. Send message")
	log.Println("3. Send gm")
	log.Println("4. Exit")
	var key int
	var content string
	var inputReader *bufio.Reader
	var err error
	inputReader = bufio.NewReader(os.Stdin)

	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		// join ChatRoom
		roomName := ""
		log.Println("join chat room name:")
		fmt.Scanf("%s\n", &roomName)
		joinChatRoomMsg := proto.JoinChatRoomReq{PlayerName: c.PlayerName, ChatRoom: roomName}
		data, _ := json.Marshal(joinChatRoomMsg)
		msg1 := proto.Message{Type: proto.MsgTypeJoinChat, Data: string(data)}
		msgData, _ := json.Marshal(msg1)
		c.conn.Write(msgData)
	case 2:
		log.Println("Say something:")
		content, err = inputReader.ReadString('\n')
		if err != nil {
			log.Println("Some error occurred when you input, error:", err)
		}

		joinChatRoomMsg := proto.SendChatMsgReq{PlayerName: c.PlayerName, ChatRoom: c.ChatRoom, Msg: content}
		data, _ := json.Marshal(joinChatRoomMsg)
		msg := proto.Message{Type: proto.MsgTypeSendChatMsg, Data: string(data)}
		msgData, _ := json.Marshal(msg)
		c.conn.Write(msgData)
	case 3:
		log.Println("Send gm")
		gmKey := ""
		param := ""
		log.Println("input gm key:")
		fmt.Scanf("%s\n", &gmKey)
		log.Println("input gm param:")
		fmt.Scanf("%s\n", &param)
		req := proto.GMReq{GM: gmKey, Data: param}
		data, _ := json.Marshal(req)
		msg := proto.Message{Type: proto.MsgTypeGMMsg, Data: string(data)}
		msgData, _ := json.Marshal(msg)
		c.conn.Write(msgData)
	case 4:
		log.Println("Exit...")
		os.Exit(0)
	default:
		log.Print("Selected invalid!\n")
	}
}
