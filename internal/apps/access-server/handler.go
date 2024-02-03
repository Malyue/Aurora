package access_server

import (
	userpb "Aurora/api/proto-go/user"
	_conn "Aurora/internal/apps/access-server/conn"
	"Aurora/internal/apps/access-server/internal/message"
	_errorx "Aurora/internal/pkg/errorx"
	_resp "Aurora/internal/pkg/resp"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

func wsHandler(s *Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		//msg, err := message.GetAuthMessage(data)
		//if err != nil {
		//	return
		//}
		var msg *message.Msg

		// valid token
		verifyTokenResp, err := s.UserServer.VerifyToken(context.Background(), &userpb.VerifyTokenRequest{
			Token: msg.Token,
		})

		if err != nil {
			msg, _ := json.Marshal(_resp.ResponseCode{
				Code: _errorx.CodeServerBusy,
				Msg:  "verify user error",
				Data: nil,
			})
			wsConn.WriteMessage(websocket.TextMessage, msg)
			wsConn.Close()
			return
		}

		if verifyTokenResp.Expire {
			msg, _ := json.Marshal(_resp.ResponseCode{
				Code: _errorx.CodeTokenExpire,
				Msg:  _errorx.CodeTokenExpire.Msg(),
				Data: nil,
			})
			wsConn.WriteMessage(websocket.TextMessage, msg)
			wsConn.Close()
			return
		}

		// keep a conn in manager
		conn := _conn.NewConn(wsConn, verifyTokenResp.Id)
		s.connManager.AddConn(conn, conn.UserId)
		// TODO set a ack model
	}
}
