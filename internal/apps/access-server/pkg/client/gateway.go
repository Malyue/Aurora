package client

import (
	_message "Aurora/internal/apps/access-server/pkg/message"
	"Aurora/internal/apps/access-server/svc"
	"errors"
	"github.com/panjf2000/ants/v2"
	"sync"
)

// Gateway is the basic and common interface for all gate implementations
// As the basic gate, it is used to provide a common gate interface for other modules to interact with the gate
type Gateway interface {
	// SetClientID sets the client id with the new id
	SetClientID(old ID, new ID) error

	//UpdateClient(id Info,info *)

	ExitClient(id ID) error

	EnqueueMessage(id ID, message *_message.Message) error

	GetClient(id ID) Client

	GetAll() map[ID]Info

	SetMessageHandler(h MessageHandler)

	AddClient(cli Client)
}

type ClientsHub struct {
	id string

	clients map[ID]Client
	mu      sync.RWMutex

	msgHandler MessageHandler

	pool *ants.Pool
	ctx  *svc.ServerCtx
}

type Options struct {
	ID                    string `yaml:"id"`
	SecretKey             string `yaml:"secretKey"`
	MaxMessageConcurrency int    `yaml:"maxMessageConcurrency"`
}

func NewClientHub(ctx *svc.ServerCtx, options *Options) (*ClientsHub, error) {
	ret := &ClientsHub{
		clients: map[ID]Client{},
		mu:      sync.RWMutex{},
		id:      options.ID,
		ctx:     ctx,
	}
	//if options.SecretKey != "" {
	//
	//}

	pool, err := ants.NewPool(options.MaxMessageConcurrency,
		ants.WithNonblocking(true),
		ants.WithPanicHandler(func(i interface{}) {
			ctx.Logger.Printf("panic : %v \n", i)
		}),
		ants.WithPreAlloc(false),
	)

	if err != nil {
		return nil, err
	}
	ret.pool = pool
	return ret, nil
}

func (c *ClientsHub) GetClient(id ID) Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.clients[id]
}

func (c *ClientsHub) GetAll() map[ID]Info {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[ID]Info)
	for id, client := range c.clients {
		result[id] = client.GetInfo()
	}
	return result
}

func (c *ClientsHub) SetMessageHandler(h MessageHandler) {
	c.msgHandler = h
}

func (c *ClientsHub) AddClient(cli Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := cli.GetInfo().ID
	id.SetGateway(c.id)

	//dc, ok := cli.(Gateway)
	//if ok {
	//
	//}

	c.clients[id] = cli
	info := cli.GetInfo()
	c.msgHandler(&info, _message.NewMessage(0, _message.ActionInternalOnline, id))

}

// SetClientID if the old is not exist or the new is exists, return error
// otherwise, the old offline and the new online
func (c *ClientsHub) SetClientID(oldID, newID ID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	oldID.SetGateway(c.id)
	newID.SetGateway(c.id)

	// check if the client is exists
	cli, ok := c.clients[oldID]
	if !ok || cli == nil {
		return errors.New(errClientNotExist)
	}

	cliLogged, exist := c.clients[newID]
	if exist && cliLogged != nil {
		return errors.New(errClientAlreadyExist)
	}

	oldInfo := cli.GetInfo()
	cli.SetID(newID)
	newInfo := cli.GetInfo()
	delete(c.clients, oldID)
	c.msgHandler(&oldInfo, _message.NewMessage(0, _message.ActionInternalOffline, oldID))
	c.msgHandler(&newInfo, _message.NewMessage(0, _message.ActionInternalOnline, newID))

	c.clients[newID] = cli
	return nil
}

func (c *ClientsHub) ExitClient(id ID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	id.SetGateway(c.id)

	cli, ok := c.clients[id]
	if !ok || cli == nil {
		return errors.New(errClientNotExist)
	}

	info := cli.GetInfo()
	cli.SetID("")
	delete(c.clients, id)
	c.msgHandler(&info, _message.NewMessage(0, _message.ActionInternalOffline, id))
	cli.Exit()

	return nil
}

func (c *ClientsHub) EnqueueMessage(id ID, msg *_message.Message) error {
	c.mu.RLock()
	defer c.mu.Unlock()

	id.SetGateway(c.id)
	cli, ok := c.clients[id]
	if !ok || cli == nil {
		return errors.New(errClientNotExist)
	}

	return c.enqueueMessage(cli, msg)
}

func (c *ClientsHub) enqueueMessage(cli Client, msg *_message.Message) error {
	if !cli.IsRunning() {
		return errors.New(errClientClosed)
	}
	err := c.pool.Submit(func() {
		_ = cli.EnqueueMessage(msg)
	})
	if err != nil {
		return errors.New("enqueue message to client failed")
	}
	return nil
}
