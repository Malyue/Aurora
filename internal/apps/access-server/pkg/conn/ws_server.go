package conn

import (
	"Aurora/internal/apps/access-server/pkg/client"
	"Aurora/internal/apps/access-server/svc"
	"fmt"
	"github.com/gorilla/mux"
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
	options   *Option
	upgrader  websocket.Upgrader
	handler   ConnectionHandler
	ctx       *svc.ServerCtx
	h         client.MessageHandler
	decorator client.Gateway
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

func (ws *WsServer) Run(host string, port int, endpoints []Router) error {
	r := mux.NewRouter()
	r.HandleFunc("/ws", ws.wsHandler)
	for _, endpoint := range endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}

//func (ws *WsServer) SetMessageHandler(h MessageHandler) {
//	ws.h = h
//}
//
//func (ws *WsServer) GetClient(id ID) Client {
//	return ws.decorator.GetClient(id)
//}
//
//func (ws *WsServer) GetAll() map[ID]Info {
//	return ws.decorator.GetAll()
//}
//
//func (ws *WsServer) AddClient(c Client) {
//	ws.decorator.AddClient(c)
//}
//
//func (ws *WsServer) SetClientID(old ID, new ID) error {
//	return ws.decorator.SetClientID(old, new)
//}
//
//func (ws *WsServer) UpdateClient(id ID, secret *ClientSecrets) error {
//	return ws.decorator.UpdateClient(id, secret)
//}
