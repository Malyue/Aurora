package client

import (
	"Aurora/internal/apps/access-server/pkg/conn"
	_message "Aurora/internal/apps/access-server/pkg/message"
	"github.com/sirupsen/logrus"
	"sync"
)

var messageReader MessageReader

var defautlCodec _message.Codec = _message.JsonCodec

// recyclePool reuse readerRes
var recyclePool sync.Pool

// MessageReader implements a reader to get message from connection
type MessageReader interface {
	// Read with block
	Read(conn conn.Conn) (*_message.Message, error)
	ReadWithCodec(conn conn.Conn, codec _message.Codec) (*_message.Message, error)
	// ReadCh returns two channel, the first use to get message and the second is used to stop the message sent
	ReadCh(conn conn.Conn) (<-chan *readerRes, chan<- struct{})
	ReadChWithCodec(conn conn.Conn, codec _message.Codec) (<-chan *readerRes, chan<- struct{})
}

type defaultReader struct{}

func InitReader() {
	recyclePool = sync.Pool{
		New: func() interface{} {
			return &readerRes{}
		},
	}
	messageReader = &defaultReader{}
}

type readerRes struct {
	err error
	m   *_message.Message
}

func (r *defaultReader) Read(conn conn.Conn) (*_message.Message, error) {
	return r.ReadWithCodec(conn, defautlCodec)
}

func (r *defaultReader) ReadWithCodec(conn conn.Conn, codec _message.Codec) (*_message.Message, error) {
	bytes, err := conn.Read()
	if err != nil {
		return nil, err
	}
	msg := _message.NewEmptyMessage()
	err = codec.Decode(bytes, msg)
	return msg, err
}

func (r *defaultReader) ReadCh(conn conn.Conn) (<-chan *readerRes, chan<- struct{}) {
	return r.ReadChWithCodec(conn, defautlCodec)
}

func (r *defaultReader) ReadChWithCodec(conn conn.Conn, codec _message.Codec) (<-chan *readerRes, chan<- struct{}) {
	c := make(chan *readerRes, 10)
	done := make(chan struct{})
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				logrus.Error("Err on readPump from connection :%v", e)
			}
		}()

		for {
			select {
			case <-done:
				goto CLOSE
			default:
				msg, err := r.Read(conn)
				res := recyclePool.Get().(*readerRes)
				if err != nil {
					res.err = err
					c <- res
					if _message.IsDecodeError(err) {
						continue
					}
					goto CLOSE
				}

				res.m = msg
				c <- res
			}
		}
	CLOSE:
		close(c)
	}()
	return c, done
}

func (r *readerRes) Recycle() {
	r.m = nil
	r.err = nil
	recyclePool.Put(r)
}
