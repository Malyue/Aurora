package logic_server

import (
	"Aurora/internal/apps/push-server/svc"
	"context"
	"net"

	"google.golang.org/grpc"

	grouppb "Aurora/api/proto-go/group"
	"Aurora/internal/apps/logic-server/service"
	discovery "Aurora/internal/pkg/etcd"
	_log "Aurora/internal/pkg/log"
	_mysql "Aurora/internal/pkg/mysql"
	_redis "Aurora/internal/pkg/redis"
	_config "Aurora/internal/tools/config"
)

type Config struct {
	_config.BaseConfig `yaml:",inline"`
	Etcd               _config.Etcd  `yaml:"etcd"`
	Log                _log.Config   `yaml:"log"`
	Mysql              _mysql.Config `yaml:"mysql"`
	Redis              _redis.Config `yaml:"redis"`
}

type Server struct {
	Cfg         Config
	Addr        string
	GroupServer *grpc.Server
	SvcCtx      *svc.ServerCtx
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	etcdRegister := discovery.NewRegister([]string{s.Cfg.Etcd.Address}, s.SvcCtx.Logger)
	defer etcdRegister.Stop()

	node := discovery.Server{
		Name: s.Cfg.Name,
		Addr: s.Addr,
	}

	if _, err := etcdRegister.Register(node, 10); err != nil {
		return err
	}

	s.SvcCtx.Logger.Info("Group Service Start ...")

	return s.GroupServer.Serve(lis)
}

func New(opts ...OptionFunc) (*Server, error) {
	s := &Server{}

	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	err := _config.GetConfig(o.ConfigFilePath, &s.Cfg)
	if err != nil {
		return nil, err
	}

	logger := _log.InitLogger(&s.Cfg.Log)

	s.SvcCtx = &svc.ServerCtx{
		Ctx:    context.Background(),
		Logger: logger,
	}

	s.SvcCtx.DBClient, err = _mysql.NewMysql(&s.Cfg.Mysql)
	if err != nil {
		s.SvcCtx.Logger.Error("init db error")
		return nil, err
	}

	s.SvcCtx.RedisCli, err = _redis.NewRedis(&s.Cfg.Redis)
	if err != nil {
		s.SvcCtx.Logger.Error("init redis error")
		return nil, err
	}

	s.SvcCtx.Cache = make(map[string]svc.Item)

	s.Addr = s.Cfg.Host + ":" + s.Cfg.Port
	grpcServer := grpc.NewServer()
	grouppb.RegisterGroupServiceServer(grpcServer, service.NewGroupServer(s.SvcCtx))

	s.SvcCtx.Logger.Info("Group Service Init...")

	return s, nil
}
