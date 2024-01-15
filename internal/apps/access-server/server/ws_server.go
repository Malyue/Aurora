package server

import (
	_conn "Aurora/internal/apps/access-server/internal/conn"
	"Aurora/internal/apps/access-server/internal/message"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 65536,
	CheckOrigin: func(r *http.Request) bool {
		// TODO add subprotocol
		return true
	},
}

func StartWSServer(address string) {
	http.HandleFunc("/ws", wsHandler)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// block to wait the authorization
	_, data, err := wsConn.ReadMessage()
	if err != nil {
		// TODO close conn
		return
	}

	msg, err := message.MessageAuth(data)
	if err != nil {
		return
	}

	// valid token
	conn := _conn.NewConn(wsConn, "")

	//DoConn(conn)
	fmt.Println(conn)
}

func DoConn(conn *_conn.Conn) {
	for {
		err := conn.WS.SetReadDeadline(time.Now().Add(12 * time.Minute))
		if err != nil {
			return
		}
		_, data, err := conn.WS.ReadMessage()
		if err != nil {

		}

		conn.HandleMessage(data)
	}
}
