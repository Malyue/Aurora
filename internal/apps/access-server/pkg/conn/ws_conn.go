package conn

import (
	"github.com/gorilla/websocket"
	"net"
	"strings"
	"time"
)

type Option struct {
	ReadTimeout  int64 `yaml:"read_timeout"`
	WriteTimeout int64 `yaml:"write_timeout"`
}

type WsConn struct {
	conn    *websocket.Conn
	options *Option
}

func NewWsConnection(conn *websocket.Conn, opt *Option) *WsConn {
	c := &WsConn{
		conn:    conn,
		options: opt,
	}
	c.conn.SetCloseHandler(func(code int, text string) error {
		return ErrClosed
	})

	return c
}

func (c *WsConn) Write(data []byte) error {
	ddl := time.Now().Add(time.Duration(c.options.WriteTimeout))
	_ = c.conn.SetWriteDeadline(ddl)

	err := c.conn.WriteMessage(websocket.TextMessage, data)
	return c.wrapError(err)
}

func (c *WsConn) Read() ([]byte, error) {
	ddl := time.Now().Add(time.Duration(c.options.ReadTimeout))
	_ = c.conn.SetReadDeadline(ddl)

	msgType, bytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, c.wrapError(err)
	}

	switch msgType {
	case websocket.TextMessage:
	case websocket.PingMessage:
	case websocket.BinaryMessage:
	default:
		return nil, ErrBadPackage
	}
	return bytes, err
}

func (c *WsConn) Close() error {
	return c.wrapError(c.conn.Close())
}

func (c *WsConn) GetConnInfo() *ConnInfo {
	c.conn.NetConn()
	remoteAddr := c.conn.RemoteAddr().(*net.TCPAddr)
	info := ConnInfo{
		Ip:   remoteAddr.IP.String(),
		Port: remoteAddr.Port,
		Addr: c.conn.RemoteAddr().String(),
	}
	return &info
}

func (c *WsConn) wrapError(err error) error {
	if err == nil {
		return nil
	}
	if websocket.IsUnexpectedCloseError(err) {
		return ErrClosed
	}
	if websocket.IsCloseError(err) {
		return ErrClosed
	}
	if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
		_ = c.conn.Close()
		return ErrClosed
	}
	if strings.Contains(err.Error(), "use of closed network conn") {
		_ = c.conn.Close()
		return ErrClosed
	}
	if strings.Contains(err.Error(), "i/o timeout") {
		return ErrReadTimeout
	}
	return err
}
