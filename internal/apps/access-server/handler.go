package access_server

import (
	userpb "Aurora/api/proto-go/user"
	_conn "Aurora/internal/apps/access-server/conn"
	"Aurora/internal/apps/access-server/internal/message"
	"context"
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

		msg, err := message.GetAuthMessage(data)
		if err != nil {
			return
		}

		// valid token
		verifyTokenResp, err := s.UserServer.VerifyToken(context.Background(), &userpb.VerifyTokenRequest{
			Token: msg.Token,
		})

		if err != nil {

		}

		if verifyTokenResp.Expire {

		}

		conn := _conn.NewConn(wsConn, verifyTokenResp.Id)
		_ = conn
	}
}
