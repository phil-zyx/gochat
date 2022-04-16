package proto

import "encoding/json"

const (
	MsgTypeLogin       = "login"
	MsgTypeJoinChat    = "join_chat"
	MsgTypeChatMsg     = "chat_msg"
	MsgTypeSendChatMsg = "send_chat_msg"
	MsgTypeGMMsg       = "gm"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func PackProtoAck(msgTyp string, data interface{}) Message {
	jsonData, _ := json.Marshal(data)
	return Message{
		Type: msgTyp,
		Data: string(jsonData),
	}
}

type ChatMessage struct {
	ID       int64
	Sender   string
	RoomName string
	Content  string
	CreateTs int64
}
