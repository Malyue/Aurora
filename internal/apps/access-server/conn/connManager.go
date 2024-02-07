package conn

import (
	"Aurora/internal/apps/access-server/internal/message"
	"encoding/json"
	"sync"
)

type ConnManager struct {
	UserConnMap map[string][]int64
	//ConnMap     map[*websocket.Conn]struct{}
	Clients    map[string]*Conn
	register   chan *Conn
	unregister chan *Conn
	broadcast  chan message.Interface
	// GroupClient
	// UserClient
	sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		//UserConnMap: make(map[string][]int64),
		Clients:    make(map[string]*Conn),
		register:   make(chan *Conn),
		unregister: make(chan *Conn),
		broadcast:  make(chan message.Interface),
	}
}

func (cm *ConnManager) Run() {
	for {
		select {
		case client := <-cm.register:
			cm.addConn(client)
		case client := <-cm.unregister:
			cm.removeConn(client.UserId)
		case msg := <-cm.broadcast:
			cm.msgBroadcast(msg)
		}
	}
}

func (cm *ConnManager) addConn(conn *Conn) {
	cm.Lock()
	defer cm.Unlock()
	cm.Clients[conn.UserId] = conn
}

func (cm *ConnManager) removeConn(id string) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.Clients, id)
}

func (cm *ConnManager) GetConn(id string) *Conn {
	cm.RLock()
	defer cm.RUnlock()
	return cm.Clients[id]
}

func (cm *ConnManager) GetBroadcast() chan message.Interface {
	return cm.broadcast
}

func (cm *ConnManager) msgBroadcast(msg message.Interface) {
	byteMessage, _ := json.Marshal(&msg)
	receivers, flag := msg.GetReceiverID()
	if flag {
		// get receivers id from db
		switch msg.GetMsgType() {
		case message.Group:
			// TODO get receiver ids
		case message.Broadcast:

		}

	}
	for i := 0; i < len(receivers); i++ {
		client, ifExist := cm.Clients[receivers[i]]
		if ifExist {
			client.send <- byteMessage
		}
	}
}
