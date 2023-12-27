package im

import (
	"Aurora/internal/apps/im/server/conn"
	"Aurora/internal/apps/im/svc"
	"sync"
	"sync/atomic"
	"time"
)

type stats struct {
	inMsgs      atomic.Int64
	outMsgs     atomic.Int64
	inBytes     atomic.Int64
	outBytes    atomic.Int64
	slowClients atomic.Int64
}

type Config struct {
}

type Server struct {
	stats
	opts            *Options
	Config          Config
	start           time.Time
	connManager     *conn.ConnManager
	svcCtx          *svc.ServerCtx
	ipBlackList     map[string]uint64
	ipBlackListLock sync.RWMutex
}

func New(opts ...OptionFunc) (*Server, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	return &Server{}, nil
}
