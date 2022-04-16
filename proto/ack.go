package proto

type ErrCode int32

type LoginAck struct {
	Err        ErrCode
	PlayerID   int64
	PlayerName string
}

type JoinChatRoomAck struct {
	Err        ErrCode
	PlayerName string
	ChatRoom   string
	RecentMsg []ChatMessage
}

type SendChatMsgAck struct {
	Err ErrCode
}

type GMAck struct {
	Err ErrCode
}