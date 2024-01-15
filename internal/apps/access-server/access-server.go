package access_server

import (
	"Aurora/internal/apps/access-server/server"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"Aurora/internal/apps/access-server/internal/conn"
	"Aurora/internal/apps/access-server/svc"
	_config "Aurora/internal/tools/config"
)

type stats struct {
	inMsgs      atomic.Int64
	outMsgs     atomic.Int64
	inBytes     atomic.Int64
	outBytes    atomic.Int64
	slowClients atomic.Int64
}

type Config struct {
	NodeId string `json:"nodeId"`
	Name   string `json:"name"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	//Address string
	// redis -- to get the conn situation
	//RedisConf redis.Config
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
	// Node Manager select a node to send msg
}

func New(opts ...OptionFunc) (*Server, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	var cfg Config
	err := _config.GetConfig(o.ConfigFilePath, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.NodeId == "" {
		id := uuid.New()
		cfg.NodeId = id.String()
	}

	return &Server{}, nil
}

func (s *Server) Run() error {
	// start server to get
	server.StartWSServer(s.Config.Host + ":" + s.Config.Port)

	return nil
}
