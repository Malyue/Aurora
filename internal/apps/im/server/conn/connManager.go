package conn

import (
	"Aurora/internal/apps/im/svc"
	"github.com/gorilla/websocket"
	"sync"
)

type ConnManager struct {
	UserConnMap map[string][]int64
	//ConnMap     map[*websocket.Conn]struct{}
	ConnMap map[string]*websocket.Conn
	sync.RWMutex
	s *svc.ServerCtx
}

func NewConnManager(s *svc.ServerCtx) *ConnManager {
	return &ConnManager{
		//UserConnMap: make(map[string][]int64),
		ConnMap: make(map[string]*websocket.Conn),
		s:       s,
	}
}

func (c *ConnManager) AddConn(conn *websocket.Conn, id string) {
	c.Lock()
	defer c.Unlock()
	c.ConnMap[id] = conn
}

func (c *ConnManager) RemoveConn(id string) {
	c.Lock()
	defer c.Unlock()
	delete(c.ConnMap, id)
}

func (c *ConnManager) GetConn(id string) *websocket.Conn {
	c.RLock()
	defer c.RUnlock()
	return c.ConnMap[id]
}
