package access_server

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_conn "Aurora/internal/apps/access-server/pkg/conn"
	_message "Aurora/internal/apps/access-server/pkg/message"
	"strconv"
	"time"
)

func (s *Server) handlerConn(conn _conn.Conn) _client.ID {
	// TODO get user id and device
	id := _client.NewID(strconv.Itoa(int(s.Config.WorkID)), "", "")
	ret := _client.NewClientWithConfig(conn, s.Gateway, s.handlerMessage, &_client.Config{
		HeartbeatLostLimit:      3,
		ClientHeartbeatDuration: time.Second * 30,
		ServerHeartbeatDuration: time.Second * 30,
		CloseImmediately:        false,
	})
	ret.SetID(id)

	s.Gateway.AddClient(ret)

	// handler conn
	ret.Run()

	//hello := _message.ServerHello{
	//	//TempID:
	//	HeartbeatInterval: 30,
	//}

	//m := _message.NewMessage(0,_message.ActionHello,hello)
	//_ = ret.EnqueueMessage(m)

	return id
}

func (s *Server) handlerMessage(cliInfo *_client.Info, message *_message.Message) {

}

//func wsHandler(s *Server) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		conn, err := Upgrader.Upgrade(w, r, nil)
//		if err != nil {
//			return
//		}
//		proxy := _conn.ConnProxy{
//			Conn: _conn.NewWsConnection(conn, &s.Config.WsOpts),
//		}
//
//
//		// block to wait the authorization package
//		//_, data, err := wsConn.ReadMessage()
//		//if err != nil {
//		//	msg, _ := json.Marshal(_resp.ResponseCode{
//		//		Code: _errorx.CodeServerBusy,
//		//		Msg:  "verify user error",
//		//		Data: nil,
//		//	})
//		//	wsConn.WriteMessage(websocket.TextMessage, msg)
//		//	wsConn.Close()
//		//	return
//		//}
//
//		//var msg *message.AuthMessage
//		//msg, err = message.HandlerAuthMessage(data)
//		//if err != nil {
//		//	msg, _ := json.Marshal(_resp.ResponseCode{
//		//		Code: _errorx.CodeServerBusy,
//		//		Msg:  "verify user error",
//		//		Data: nil,
//		//	})
//		//	wsConn.WriteMessage(websocket.TextMessage, msg)
//		//	wsConn.Close()
//		//	return
//		//}
//
//		// valid token
//		//verifyTokenResp, err := s.UserServer.VerifyToken(context.Background(), &userpb.VerifyTokenRequest{
//		//	Token: msg.Token,
//		//})
//
//		//if err != nil {
//		//	msg, _ := json.Marshal(_resp.ResponseCode{
//		//		Code: _errorx.CodeServerBusy,
//		//		Msg:  "verify user error",
//		//		Data: nil,
//		//	})
//		//	wsConn.WriteMessage(websocket.TextMessage, msg)
//		//	wsConn.Close()
//		//	return
//		//}
//
//		//if verifyTokenResp.Expire {
//		//	msg, _ := json.Marshal(_resp.ResponseCode{
//		//		Code: _errorx.CodeTokenExpire,
//		//		Msg:  _errorx.CodeTokenExpire.Msg(),
//		//		Data: nil,
//		//	})
//		//	wsConn.WriteMessage(websocket.TextMessage, msg)
//		//	wsConn.Close()
//		//	return
//		//}
//
//		// keep a conn in manager
//		//_conn.NewConn(wsConn, verifyTokenResp.Id, s.connManager)
//
//		// TODO set a ack model
//
//	}
//}

//func HandlerRouter(callback func(action _message.Action, fn HandlerFunc) HandlerFunc) {
//	m := map[_message.Action]HandlerFunc {
//		_message.ActionChatMessage:
//	}
//}
