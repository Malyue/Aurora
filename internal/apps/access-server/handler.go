package access_server

import (
	"Aurora/internal/apps/access-server/pkg/client"
	_conn "Aurora/internal/apps/access-server/pkg/conn"
	_message "Aurora/internal/apps/access-server/pkg/message"
	"net/http"
)


func wsHandler(s *Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		proxy := _conn.ConnProxy{
			conn: _conn.NewWsConnection(conn, &s.Config.ConnOpt),
		}


		// block to wait the authorization package
		//_, data, err := wsConn.ReadMessage()
		//if err != nil {
		//	msg, _ := json.Marshal(_resp.ResponseCode{
		//		Code: _errorx.CodeServerBusy,
		//		Msg:  "verify user error",
		//		Data: nil,
		//	})
		//	wsConn.WriteMessage(websocket.TextMessage, msg)
		//	wsConn.Close()
		//	return
		//}

		//var msg *message.AuthMessage
		//msg, err = message.HandlerAuthMessage(data)
		//if err != nil {
		//	msg, _ := json.Marshal(_resp.ResponseCode{
		//		Code: _errorx.CodeServerBusy,
		//		Msg:  "verify user error",
		//		Data: nil,
		//	})
		//	wsConn.WriteMessage(websocket.TextMessage, msg)
		//	wsConn.Close()
		//	return
		//}

		// valid token
		//verifyTokenResp, err := s.UserServer.VerifyToken(context.Background(), &userpb.VerifyTokenRequest{
		//	Token: msg.Token,
		//})

		//if err != nil {
		//	msg, _ := json.Marshal(_resp.ResponseCode{
		//		Code: _errorx.CodeServerBusy,
		//		Msg:  "verify user error",
		//		Data: nil,
		//	})
		//	wsConn.WriteMessage(websocket.TextMessage, msg)
		//	wsConn.Close()
		//	return
		//}

		//if verifyTokenResp.Expire {
		//	msg, _ := json.Marshal(_resp.ResponseCode{
		//		Code: _errorx.CodeTokenExpire,
		//		Msg:  _errorx.CodeTokenExpire.Msg(),
		//		Data: nil,
		//	})
		//	wsConn.WriteMessage(websocket.TextMessage, msg)
		//	wsConn.Close()
		//	return
		//}

		// keep a conn in manager
		//_conn.NewConn(wsConn, verifyTokenResp.Id, s.connManager)

		// TODO set a ack model

	}
}

func HandlerRouter(callback func(action _message.Action, fn HandlerFunc) HandlerFunc) {
	m := map[_message.Action]HandlerFunc {
		_message.ActionChatMessage:
	}
}
