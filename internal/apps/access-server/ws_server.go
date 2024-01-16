package access_server

import (
	_conn "Aurora/internal/apps/access-server/conn"
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

func (s *Server) StartWSServer() {
	http.HandleFunc("/ws", wsHandler(s))
	err := http.ListenAndServe(s.Config.Host+":"+s.Config.Port, nil)
	if err != nil {
		panic(err)
	}
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
