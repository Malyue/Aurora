package conn

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Conn struct {
	mutex  sync.Mutex
	WS     *websocket.Conn
	UserId string
	//DeviceId  int64
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan []byte
	isClose   bool
}

func NewConn(c *websocket.Conn, userId string) *Conn {
	conn := &Conn{
		inChan:    make(chan []byte, 1024),
		outChan:   make(chan []byte, 1024),
		closeChan: make(chan []byte, 1),
		WS:        c,
		UserId:    userId,
	}

	go conn.ReadMsgLoop()
	go conn.WriteMsgLoop()

	return conn
}

func (c *Conn) HandleMessage(bytes []byte) {
	// TODO Unmarshal the msg
	// switch case to judge the type of the message
}

func (c *Conn) Write(bytes []byte) error {
	return c.WriteToWs(bytes)
}

func (c *Conn) WriteToWs(bytes []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.WS.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
	if err != nil {
		return err
	}
	return c.WS.WriteMessage(websocket.BinaryMessage, bytes)
}

func (c *Conn) Close() error {
	// TODO subscribedRoom
	//go func() {
	//
	//}
	c.mutex.Lock()
	if !c.isClose {
		close(c.closeChan)
		c.isClose = true
	}
	// TODO delete from connmanager
	c.mutex.Unlock()
	return c.WS.Close()
}

func (c *Conn) GetAddr() string {
	return c.WS.RemoteAddr().String()
}

func (c *Conn) ReadMsgLoop() {
	for {
		var (
			data []byte
			err  error
		)
		// receive data
		if _, data, err = c.WS.ReadMessage(); err != nil {
			c.Close()
		}
		// write data
		if err = c.InChanWrite(data); err != nil {
			c.Close()
		}
	}
}

func (c *Conn) WriteMsgLoop() {
	for {
		var (
			data []byte
			err  error
		)
		if data, err = c.OutChanRead(); err != nil {
			c.Close()
		}
		if err = c.WS.WriteMessage(1, data); err != nil {
			c.Close()
		}
	}
}

func (c *Conn) InChanRead() (data []byte, err error) {
	select {
	case data = <-c.inChan:
	case <-c.closeChan:
		err = errors.New("conn is closed")
	}
	return nil, err
}

func (c *Conn) InChanWrite(data []byte) (err error) {
	select {
	case c.outChan <- data:
	case <-c.closeChan:
		err = errors.New("conn is closed")
	}
	return err
}

func (c *Conn) OutChanRead() (data []byte, err error) {
	select {
	case data = <-c.outChan:
	case <-c.closeChan:
		err = errors.New("conn is closed")
	}
	return nil, err
}

func (c *Conn) OutChanWrite(data []byte) (err error) {
	select {
	case c.outChan <- data:
	case <-c.closeChan:
		err = errors.New("conn is closed")
	}
	return err
}
