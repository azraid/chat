package main

import (
	"fmt"
	"os"

	. "github.com/azraid/chat/services/chat"
	"github.com/azraid/pasque/app"
	n "github.com/azraid/pasque/core/net"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) chatroomsrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	cli := n.NewClient(eid)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(GetRoomMsg{}), OnGetRoom)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(JoinRoomMsg{}), OnJoinRoom)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(SendChatMsg{}), OnSendChat)

	toplgy := n.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "RoomID",
		FederatedApis: cli.ListGridApis()}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
