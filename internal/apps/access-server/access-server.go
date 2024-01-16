package access_server

import (
	userpb "Aurora/api/proto-go/user"
	"Aurora/internal/apps/access-server/conn"
	discovery "Aurora/internal/pkg/etcd"
	_grpc "Aurora/internal/pkg/grpc"
	_log "Aurora/internal/pkg/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/resolver"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	_conn "Aurora/internal/apps/access-server/conn"
	_pkg "Aurora/internal/apps/access-server/pkg/snowflake"
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
	NodeId string       `json:"nodeId"`
	Name   string       `json:"name"`
	Host   string       `json:"host"`
	Port   string       `json:"port"`
	WorkID int64        `json:"workId"`
	Etcd   _config.Etcd `yaml:"etcd"`
	Log    _log.Config  `yaml:"log"`
	//Address string
	// redis -- to get the conn situation
	//RedisConf redis.Config
}

type Server struct {
	stats
	//opts            *Options
	Config          Config
	start           time.Time
	connManager     *conn.ConnManager
	svcCtx          *svc.ServerCtx
	ipBlackList     map[string]uint64
	ipBlackListLock sync.RWMutex

	NodeSnowFlake *_pkg.Worker
	// Node Manager select a node to send msg

	// grpc client
	UserServer userpb.UserServiceClient
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

	logger := _log.InitLogger(&cfg.Log)
	// init snowflake
	snowflakeWorker, err := _pkg.NewWorker(cfg.WorkID)
	if err != nil {
		return nil, err
	}

	// init connManager
	connManager := _conn.NewConnManager()

	// TODO init timingWheel

	// TODO init redis

	// add grpc client
	etcdResolver := discovery.NewResolver([]string{cfg.Etcd.Address}, logger)
	resolver.Register(etcdResolver)
	defer etcdResolver.Close()

	userServer, err := _grpc.InitUserClient()
	if err != nil {
		logrus.Errorf("create user client err : %s", err)
	}

	return &Server{
		Config:        cfg,
		NodeSnowFlake: snowflakeWorker,
		start:         time.Now(),
		connManager:   connManager,
		UserServer:    userServer,
	}, nil
}

func (s *Server) Run() error {
	// start server to get
	s.StartWSServer()

	return nil
}