package access_server

import (
	"github.com/gorilla/websocket"
	"net/http"
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
