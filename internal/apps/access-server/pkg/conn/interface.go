package conn

import "errors"

var (
	ErrClosed      = errors.New("conn closed")
	ErrReadTimeout = errors.New("i/o timeout")
	ErrBadPackage  = errors.New("bad package data")
)

type Conn interface {
	Write(data []byte) error
	Read() ([]byte, error)
	Close() error
	GetConnInfo() *ConnInfo
}

type ConnInfo struct {
	Ip   string
	Port int
	Addr string
}

type ConnProxy struct {
	conn Conn
}

func (c ConnProxy) Write(data []byte) error {
	return c.conn.Write(data)
}

func (c ConnProxy) Read() ([]byte, error) {
	return c.conn.Read()
}

func (c ConnProxy) Close() error {
	return c.conn.Close()
}

func (c ConnProxy) GetConnInfo() *ConnInfo {
	return c.conn.GetConnInfo()
}
