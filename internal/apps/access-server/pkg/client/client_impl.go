package client

import (
	"Aurora/internal/apps/access-server/pkg/conn"
	_message "Aurora/internal/apps/access-server/pkg/message"
	"Aurora/internal/apps/access-server/pkg/timingWheel"
	"Aurora/internal/apps/access-server/svc"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var tw = timingWheel.NewTimingWheel(500*time.Millisecond, 2000)

type MessageInterceptor = func(client Client, msg *_message.Message) bool

const (
	defaultServerHeartbeatDuration = time.Second * 30
	defaultHeartbeatDuration       = time.Second * 20
	defaultHeartbeatLostLimit      = 3
	defaultCloseImmediately        = false
)

// ClientConfig client config
type Config struct {

	// ClientHeartbeatDuration is the duration of heartbeat.
	ClientHeartbeatDuration time.Duration

	// ServerHeartbeatDuration is the duration of server heartbeat.
	ServerHeartbeatDuration time.Duration

	// HeartbeatLostLimit is the max lost heartbeat count.
	HeartbeatLostLimit int

	// CloseImmediately true express when client exit, discard all message in queue, and close connection immediately,
	// otherwise client will close runRead, and mark as stateClosing, the client cannot receive and enqueue message,
	// after all message in queue is sent, client will close runWrite and connection.
	CloseImmediately bool
}

// client state
const (
	_ int32 = iota
	// stateRunning client is running, can runRead and runWrite message.
	stateRunning
	// stateClosed client is closed, cannot do anything.
	stateClosed
)

type UserClient struct {
	conn  conn.Conn
	state int32
	// queueMessage message count in the messages channel
	queueMessage int64
	messages     chan *_message.Message
	// closeReadCh is the channel for runRead goroutine to close
	closeReadCh chan struct{}
	// closeWriteCh is the channel for runWrite goroutine to close
	closeWriteCh chan struct{}

	closeWriteOnce sync.Once
	closeReadOnce  sync.Once

	hbC    *timingWheel.Timer
	hbS    *timingWheel.Timer
	hbLost int

	info *Info

	//hub conn.ConnManager
	// mgr the client manager which manage this client
	mgr        Gateway
	msgHandler MessageHandler
	ctx        *svc.ServerCtx
	cfg        *Config

	credentials *ClientAuthCredentials
}

func NewClient(conn conn.Conn, mgr Gateway, msgHandler MessageHandler) Client {
	return NewClientWithConfig(conn, mgr, msgHandler, nil)
}

func NewClientWithConfig(conn conn.Conn, mgr Gateway, msgHandler MessageHandler, cfg *Config) Client {
	if cfg == nil {
		cfg = &Config{
			ClientHeartbeatDuration: defaultHeartbeatDuration,
			ServerHeartbeatDuration: defaultServerHeartbeatDuration,
			HeartbeatLostLimit:      defaultHeartbeatLostLimit,
			CloseImmediately:        false,
		}
	}

	ret := UserClient{
		conn:         conn,
		messages:     make(chan *_message.Message, 100),
		closeReadCh:  make(chan struct{}),
		closeWriteCh: make(chan struct{}),
		msgHandler:   msgHandler,
		// TODO init hbC and hbS
		info: &Info{
			ConnectionAt: time.Now().UnixMilli(),
			CliAddr:      conn.GetConnInfo().Addr,
		},
		cfg: cfg,
		mgr: mgr,
	}

	return &ret
}

func (c *UserClient) IsRunning() bool {
	return atomic.LoadInt32(&c.state) == stateRunning
}

func (c *UserClient) GetInfo() Info {
	return *c.info
}

func (c *UserClient) SetID(id ID) {
	c.info.ID = id
}

// readPump read message from readChan
func (c *UserClient) readPump() {
	defer func() {
		err := recover()
		if err != nil {
			c.ctx.Logger.Error("read message panic: %v", err)
		}
	}()

	readChan, done := messageReader.ReadCh(c.conn)
	var closeReason string
	for {
		select {
		case <-c.closeReadCh:
			if closeReason == "" {
				closeReason = "closed initiative"
			}
			goto STOP
		case msg := <-readChan:
			if msg == nil {
				closeReason = "readCh closed"
				c.Exit()
				continue
			}
			if msg.err != nil {
				if _message.IsDecodeError(msg.err) {
					_ = c.EnqueueMessage(_message.NewMessage(0, _message.ActionNotifyError, msg.err.Error()))
					continue
				}
				closeReason = msg.err.Error()
				c.Exit()
				continue
			}
			if c.info.ID == "" {
				closeReason = "client not login"
				c.Exit()
				break
			}
			// c.hbLost = 0
			// c.hbC.Cancel()
			// c.hbC = tw.After(c.cfg.ClientHeartbeatDuration)

			if msg.m.GetAction() == _message.ActionHello {
				c.handleHello(msg.m)
			} else {
				// handler message and send the msg to the receivers
				c.msgHandler(c.info, msg.m)
			}
			// recycle readerRes object
			msg.Recycle()
		}

	}
STOP:
	close(done)
	// c.hbC.Cancel()
	c.ctx.Logger.Info("Read exit,reason : %s", closeReason)
}

// writePump write msg to the client
func (c *UserClient) writePump() {
	defer func() {
		err := recover()
		if err != nil {
			c.ctx.Logger.Debugf("write message error, exit client: %v", err)
			c.Exit()
		}
	}()
	var closeReason string
	for {
		select {
		case <-c.closeWriteCh:
			if closeReason == "" {
				closeReason = "closed initiative"
			}
			goto STOP
		//case <-c.hbS.C:
		//	if !c.IsRunning() {
		//		closeReason = "client not running"
		//		goto STOP
		//	}
		//	_ = c.EnqueueMessage(messages.NewMessage(0, messages.ActionHeartbeat, nil))
		//	c.hbS.Cancel()
		//	c.hbS = tw.After(c.config.ServerHeartbeatDuration)
		case m := <-c.messages:
			if m == nil {
				closeReason = "message is nil, maybe client has closed"
				c.Exit()
				break
			}
			c.write2Conn(m)
			//c.hbS.Cancel()
			//c.hbS = tw.After(c.config.ServerHeartbeatDuration)
		}
	}
STOP:
	//c.hbS.Cancel()
	c.ctx.Logger.Debugf("write exit, addr=%s, reason:%s", c.info.CliAddr, closeReason)
}

// EnqueueMessage send msg in messages
func (c *UserClient) EnqueueMessage(msg *_message.Message) error {
	if atomic.LoadInt32(&c.state) == stateClosed {
		return errors.New("client has closed")
	}
	c.ctx.Logger.Info("EnqueueMessage ID = %s , msg = %v", c.info.ID, msg)
	select {
	case c.messages <- msg:
		atomic.AddInt64(&c.queueMessage, 1)
	default:
		c.ctx.Logger.Info("The msg channel is full, id = %v", c.info.ID)
	}
	return nil
}

func (c *UserClient) Exit() {
	// if it is closed, return it
	if atomic.LoadInt32(&c.state) == stateClosed {
		return
	}
	// set it as closed
	atomic.StoreInt32(&c.state, stateClosed)

	// remove client from manager
}

func (c *UserClient) Run() {
	c.ctx.Logger.Info("new client running")
	atomic.StoreInt32(&c.state, stateRunning)
	c.closeWriteOnce = sync.Once{}
	c.closeReadOnce = sync.Once{}

	go c.readPump()
	go c.writePump()
}

func (c *UserClient) isClosed() bool {
	return atomic.LoadInt32(&c.state) == stateClosed
}

func (c *UserClient) close() {
	close(c.messages)
	_ = c.conn.Close()
}

func (c *UserClient) write2Conn(m *_message.Message) {
	// encode msg as byte
	b, err := defautlCodec.Encode(m)
	if err != nil {
		c.ctx.Logger.Error("serialize output message", err)
		return
	}

	err = c.conn.Write(b)
	// sub queue message
	atomic.AddInt64(&c.queueMessage, -1)
	if err != nil {
		c.ctx.Logger.Error("runWrite err : %s", err)
		c.closeWriteOnce.Do(func() {
			close(c.closeWriteCh)
		})
	}
}

func (c *UserClient) SetCredentials(credentials *ClientAuthCredentials) {
	c.credentials = credentials
	c.info.ConnectionId = credentials.ConnectID
	if credentials.ConnectConfig != nil {
		c.cfg.HeartbeatLostLimit = credentials.ConnectConfig.AllowMaxHeartbeatLost
		c.cfg.CloseImmediately = credentials.ConnectConfig.CloseImmediately
		c.cfg.ClientHeartbeatDuration = time.Duration(credentials.ConnectConfig.HeartbeatDuration) * time.Second
	}
}

func (c *UserClient) GetCredentials() *ClientAuthCredentials {
	return c.credentials
}

func (c *UserClient) AddMessageInterceptor(interceptor MessageInterceptor) {
	h := c.msgHandler
	c.msgHandler = func(cliInfo *Info, msg *_message.Message) {
		if interceptor(c, msg) {
			return
		}
		h(cliInfo, msg)
	}
}
