package access_server

import (
	userpb "Aurora/api/proto-go/user"
	_client "Aurora/internal/apps/access-server/pkg/client"
	_handler "Aurora/internal/apps/access-server/pkg/handler"
	_message "Aurora/internal/apps/access-server/pkg/message"
	_sony "Aurora/internal/apps/access-server/pkg/sonyflake"
	discovery "Aurora/internal/pkg/etcd"
	_grpc "Aurora/internal/pkg/grpc"
	_log "Aurora/internal/pkg/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/resolver"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	_conn "Aurora/internal/apps/access-server/pkg/conn"
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
	//NodeId uint16       `json:"nodeId"`
	Name           string                         `json:"name"`
	Host           string                         `json:"host"`
	Port           string                         `json:"port"`
	WorkID         uint16                         `json:"workId"`
	Etcd           _config.Etcd                   `yaml:"etcd"`
	Log            _log.Config                    `yaml:"log"`
	WsOpts         _conn.Option                   `yaml:"ws_opt"`
	GwOpts         _client.Options                `yaml:"gateway_opt"`
	MsgHandlerOpts _handler.MessageHandlerOptions `yaml:"msg_handler_opts"`
	//Address string
	// redis -- to get the conn situation
	//RedisConf redis.Config
}

type Server struct {
	stats
	//opts            *Options
	Config Config
	start  time.Time
	//connManager     *conn.ConnManager
	svcCtx          *svc.ServerCtx
	ipBlackList     map[string]uint64
	ipBlackListLock sync.RWMutex

	Gateway _client.Gateway
	// Server the ws_server(includes run and handler conn)
	Server _conn.Server

	//NodeSnowFlake *_pkg.Worker
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

	//if cfg.NodeId == "" {
	//	id := uuid.New()
	//	cfg.NodeId = id.String()
	//}

	logger := _log.InitLogger(&cfg.Log)
	// init snowflake
	//snowflakeWorker, err := _pkg.NewWorker(cfg.WorkID)
	//if err != nil {
	//	return nil, err
	//}

	// init connManager
	//connManager := _conn.NewConnManager()

	// TODO init timingWheel

	// TODO init redis

	// init sonyflake
	_sony.Init(cfg.WorkID)

	Init()

	// add grpc client
	etcdResolver := discovery.NewResolver([]string{cfg.Etcd.Address}, logger)
	resolver.Register(etcdResolver)
	defer etcdResolver.Close()

	userServer, err := _grpc.InitUserClient()
	if err != nil {
		logrus.Errorf("create user client err : %s", err)
	}

	ctx := &svc.ServerCtx{
		Logger: logger,
	}

	wsServer := _conn.NewWsServer(ctx, &cfg.WsOpts)
	gateway, err := _client.NewClientHub(ctx, &cfg.GwOpts)
	if err != nil {
		ctx.Logger.Errorf("New Gateway error : %s", err)
	}

	handler, err := _handler.NewHandlerWithOptions(gateway, ctx, &cfg.MsgHandlerOpts)
	if err != nil {
		ctx.Logger.Errorf("New Handler error : %s", err)
	}

	gateway.SetMessageHandler(func(cliInfo *_client.Info, message *_message.Message) {
		err = handler.Handle(cliInfo, message)
		if err != nil {
			ctx.Logger.Errorf("handler message error : %s", err)
		}
	})

	return &Server{
		Config: cfg,
		//NodeSnowFlake: snowflakeWorker,
		start: time.Now(),
		//connManager: connManager,
		Server:     wsServer,
		Gateway:    gateway,
		UserServer: userServer,
	}, nil
}

func (s *Server) Run() error {
	// start server to get
	//s.StartWSServer()
	s.Server.SetConnHandler(func(conn _conn.Conn) {
		s.handlerConn(conn)
	})
	port, err := strconv.Atoi(s.Config.Port)
	if err != nil {
		return err
	}
	return s.Server.Run(s.Config.Host, port)
}
