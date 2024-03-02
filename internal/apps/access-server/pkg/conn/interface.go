package conn

import (
	"errors"
	"time"
)

var (
	ErrClosed      = errors.New("Conn closed")
	ErrReadTimeout = errors.New("i/o timeout")
	ErrBadPackage  = errors.New("bad package data")
)

type Conn interface {
	Write(msgType int, data []byte) error
	Read() ([]byte, error)
	Close() error
	GetConnInfo() *ConnInfo
	SetPongHandler(h func(appData string) error)
	SetReadDeadLine(t time.Time) error
}

type ConnInfo struct {
	Ip   string
	Port int
	Addr string
}

type ConnProxy struct {
	Conn Conn
}

func (c ConnProxy) Write(data []byte) error {
	return c.Conn.Write(data)
}

func (c ConnProxy) Read() ([]byte, error) {
	return c.Conn.Read()
}

func (c ConnProxy) Close() error {
	return c.Conn.Close()
}

func (c ConnProxy) GetConnInfo() *ConnInfo {
	return c.Conn.GetConnInfo()
}
