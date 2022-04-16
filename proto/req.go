package proto

type LoginReq struct {
	PlayerID int64
	PlayerName string
}

type JoinChatRoomReq struct {
	PlayerName string
	ChatRoom string
}

type SendChatMsgReq struct {
	PlayerName string
	ChatRoom string
	Msg string
}

type GMReq struct {
	GM string
	Data string
}