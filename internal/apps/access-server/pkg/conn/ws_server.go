package conn

import (
	"Aurora/internal/apps/access-server/svc"
	"fmt"
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

var (
	DefaultReadTimeout  int64 = 8 * 60
	DefaultWriteTimeout int64 = 8 * 60
)

type WsServer struct {
	options  *Option
	upgrader websocket.Upgrader
	handler  ConnectionHandler
	ctx      *svc.ServerCtx
}

func NewWsServer(ctx *svc.ServerCtx, options *Option) *WsServer {
	if options == nil {
		options = &Option{
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
		}
	}

	ws := &WsServer{
		options:  options,
		upgrader: Upgrader,
		ctx:      ctx,
	}
	return ws
}

func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	proxy := ConnProxy{
		Conn: NewWsConnection(conn, ws.options),
	}
	ws.handler(proxy)
}

func (ws *WsServer) SetConnHandler(handler ConnectionHandler) {
	ws.handler = handler
}

func (ws *WsServer) Run(host string, port int) error {
	http.HandleFunc("/ws", ws.wsHandler)
	addr := fmt.Sprintf("%s:%d", host, port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}
