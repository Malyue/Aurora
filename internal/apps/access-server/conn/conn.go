package conn

import (
	_message "Aurora/internal/apps/access-server/pkg/message"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Conn struct {
	mutex  sync.Mutex
	Conn   *websocket.Conn
	UserId string
	//DeviceId  int64
	//inChan    chan []byte
	//outChan   chan []byte
	send      chan []byte
	closeChan chan []byte
	isClose   bool
	hub       *ConnManager
}

func NewConn(c *websocket.Conn, userId string, hub *ConnManager) *Conn {
	conn := &Conn{
		send:      make(chan []byte),
		closeChan: make(chan []byte, 1),
		Conn:      c,
		UserId:    userId,
		hub:       hub,
	}

	conn.hub.register <- conn

	go conn.ReadPump()
	go conn.WritePump()

	return conn
}

func (c *Conn) Write(bytes []byte) error {
	return c.WriteToWs(bytes)
}

func (c *Conn) WriteToWs(bytes []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
	if err != nil {
		return err
	}
	return c.Conn.WriteMessage(websocket.BinaryMessage, bytes)
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
	// TODO delete from conn manager
	c.mutex.Unlock()
	return c.Conn.Close()
}

func (c *Conn) GetAddr() string {
	return c.Conn.RemoteAddr().String()
}

// ReadPump pumps messages from the websocket connection to the hub.
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Conn) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, byteMsg, err := c.Conn.ReadMessage()
		if err != nil {
			//if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway,websocket.CloseAbnormalClosure){
			//
			//}
			logrus.Error("Read Message error: %s", err)
			break
		}

		// parse msg as message.Interface

		msg, _ := _message.Decode(byteMsg)

		switch msg.Type {
		case _message.Person:
			msg.Msg, _ = _message.ParsePersonMessage(byteMsg)
		case _message.Broadcast:
			msg.Msg, _ = _message.ParseBroadcastMessage(byteMsg)
		}

		c.hub.GetBroadcast() <- msg.Msg
	}
}

// WritePump pumps messages from the hub to the websocket connection.
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Conn) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
				return
			}

			w.Write(msg)

			// Add queued chat msg to the current ws message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, newline); err != nil {
				return
			}
		}
	}
}
