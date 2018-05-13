package chat

import (
	"time"

	co "github.com/azraid/pasque/core"
)

type CreateRoomMsg struct {
	UserID co.TUserID
}

type CreateRoomMsgR struct {
	RoomID string
}

type ListMyRoomsMsg struct {
	UserID co.TUserID
}

type ListMyRoomsMsgR struct {
	Rooms []struct {
		RoomID string
		Lasted time.Time
	}
}

type SendChatMsg struct {
	UserID   co.TUserID
	RoomID   string
	ChatType int
	Msg      string
}

type SendChatMsgR struct {
}

type RecvChatMsg struct {
	UserID     co.TUserID
	ChatUserID co.TUserID
	RoomID     string
	ChatType   int
	Msg        string
}

type RecvChatMsgR struct {
}
