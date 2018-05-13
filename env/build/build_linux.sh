#!/bin/bash


if [ ! -d "$GOPATH/bin/chat" ]; then
   mkdir $GOPATH/bin/chat
fi

if [ ! -d "$GOPATH/bin/chat/linux_amd64" ]; then
   mkdir $GOPATH/bin/chat/linux_amd64
fi

if [ ! -d "$GOPATH/bin/chat/linux_amd64/config" ]; then
   mkdir $GOPATH/bin/chat/linux_amd64/config
fi

go build  -race --o $GOPATH/bin/chat/linux_amd64/spawn $GOPATH/src/github.com/azraid/pasque/bus/spawn/main.go 
go build  -race --o $GOPATH/bin/chat/linux_amd64/router $GOPATH/src/github.com/azraid/pasque/bus/router/main.go $GOPATH/src/github.com/azraid/pasque/bus/router/router.go
go build  -race --o $GOPATH/bin/chat/linux_amd64/sgate $GOPATH/src/github.com/azraid/pasque/bus/sgate/main.go $GOPATH/src/github.com/azraid/pasque/bus/sgate/gate.go
go build -race -o $GOPATH/bin/chat/linux_amd64/tcgate $GOPATH/src/github.com/azraid/pasque/bus/tcgate/main.go $GOPATH/src/github.com/azraid/pasque/bus/tcgate/gate.go $GOPATH/src/github.com/azraid/pasque/bus/tcgate/stub.go
go build --o $GOPATH/bin/chat/linux_amd64/logsrv $GOPATH/src/github.com/azraid/pasque/bus/logsrv/main.go 
go build  -race --o $GOPATH/bin/chat/linux_amd64/sesssrv $GOPATH/src/github.com/azraid/pasque/services/auth/sesssrv/main.go $GOPATH/src/github.com/azraid/pasque/services/auth/sesssrv/db.go  $GOPATH/src/github.com/azraid/pasque/services/auth/sesssrv/grid.go  $GOPATH/src/github.com/azraid/pasque/services/auth/sesssrv/txn.go 

go build -race -o $GOPATH/bin/chat/linux_amd64/chatroomsrv $GOPATH/src/github.com/azraid/chat/services/chat/chatroomsrv/main.go $GOPATH/src/github.com/azraid/chat/services/chat/chatroomsrv/grid.go  $GOPATH/src/github.com/azraid/chat/services/chat/chatroomsrv/txn.go
go build -race -o $GOPATH/bin/chat/linux_amd64/chatusersrv $GOPATH/src/github.com/azraid/chat/services/chat/chatusersrv/main.go $GOPATH/src/github.com/azraid/chat/services/chat/chatusersrv/grid.go  $GOPATH/src/github.com/azraid/chat/services/chat/chatusersrv/txn.go

cp -rf $GOPATH/src/github.com/azraid/chat/env/config/system_linux.json $GOPATH/bin/chat/linux_amd64/config/system.json
cp -rf $GOPATH/src/github.com/azraid/chat/env/run/run_linux.sh $GOPATH/bin/chat/linux_amd64/run.sh
cp -rf $GOPATH/src/github.com/azraid/pasque/env/config/userauthdb.json $GOPATH/bin/chat/linux_amd64/config/userauthdb.json
