package main

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/azraid/chat/services/chat"
	"github.com/azraid/pasque/app"
	. "github.com/azraid/pasque/core"
	n "github.com/azraid/pasque/core/net"
	"github.com/azraid/pasque/services/auth"
)

func OnCreateRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body CreateRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	roomID := GenerateGuid().String()
	if r, err := cli.SendReq(SpnChatRoom, "JoinRoom", JoinRoomMsg{RoomID: roomID, UserID: body.UserID}); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gridData
	} else if r.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, r.Header.GetError(), nil)
		return gridData
	}

	gd := getGridData(ToUserID(req.Header.Key), gridData)
	gd.Rooms[roomID] = ChatRoom{Lasted: time.Now()}

	cli.SendRes(req, CreateRoomMsgR{RoomID: roomID})
	return gd
}

func OnJoinRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body JoinRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if r, err := cli.SendReq(SpnChatRoom, "JoinRoom", JoinRoomMsg{RoomID: body.RoomID, UserID: body.UserID}); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gridData
	} else if r.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, r.Header.GetError(), nil)
		return gridData
	}

	gd := getGridData(ToUserID(req.Header.Key), gridData)
	gd.Rooms[body.RoomID] = ChatRoom{Lasted: time.Now()}

	cli.SendRes(req, JoinRoomMsgR{})
	return gd
}

//ListRooms 사용자가 채팅중인 방 리스트를 보여준다.
func OnListMyRooms(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {

	var body ListMyRoomsMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	gd := getGridData(ToUserID(req.Header.Key), gridData)

	res := ListMyRoomsMsgR{}
	res.Rooms = make([]struct {
		RoomID string
		Lasted time.Time
	}, len(gd.Rooms))

	i := 0
	for k, _ := range gd.Rooms {
		res.Rooms[i].RoomID = k
		i++
	}

	if err := cli.SendRes(req, res); err != nil {
		app.ErrorLog(err.Error())
	}

	return gd
}

//SendChatMsg 채팅 메세지를 전송한다.
func OnSendChat(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {

	var body SendChatMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	userID := ToUserID(req.Header.Key)
	gd := getGridData(userID, gridData)

	if v, ok := gd.Rooms[body.RoomID]; !ok {
		app.ErrorLog("RoomID[%s] not found", body.RoomID)
		cli.SendResWithError(req, RaiseNError(NErrorChatNotFoundRoomID), nil)
		return gd
	} else {
		v.Lasted = time.Now()
	}

	chatroomReq := SendChatMsg{UserID: userID, RoomID: body.RoomID, ChatType: 1, Msg: body.Msg}

	res, err := cli.SendReq(SpnChatRoom, "SendChat", chatroomReq)
	if err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gd
	}

	if res.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, res.Header.GetError(), nil)
		return gd
	}

	if err := cli.SendRes(req, SendChatMsgR{}); err != nil {
		app.ErrorLog(err.Error())
	}
	return gd
}

//RecvChatMsg 채팅 메세지를 수신한다.
func OnRecvChat(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body RecvChatMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	userID := ToUserID(req.Header.Key)
	gd := getGridData(userID, gridData)
	if v, ok := gd.Rooms[body.RoomID]; ok {
		v.Lasted = time.Now()
	}

	res, err := cli.SendReq(SpnSession, "GetUserLocation", auth.GetUserLocationMsg{UserID: userID,
		Spn: GameSpn})
	if err != nil {
		app.DebugLog("no user session at OnRecvChat")
		return gd
	}

	var rbody auth.GetUserLocationMsgR
	if err := json.Unmarshal(res.Body, &rbody); err != nil {
		app.ErrorLog(err.Error())
		return gd
	}

	cli.SendReqDirect(GameSpn, rbody.GateEid, rbody.Eid, "RecvChat", body)

	fmt.Printf("%s:%s-%s\r\n", body.ChatUserID, body.Msg, time.Now().Format(time.RFC3339))

	return gd
}
