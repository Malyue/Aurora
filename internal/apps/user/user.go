package user

import (
	userpb "Aurora/api/proto-go/user"
	"Aurora/internal/apps/user/service"
	"Aurora/internal/pkg/jwt"
	"context"
	"net"

	"google.golang.org/grpc"

	"Aurora/internal/apps/user/svc"
	discovery "Aurora/internal/pkg/etcd"
	_log "Aurora/internal/pkg/log"
	_mysql "Aurora/internal/pkg/mysql"
	_redis "Aurora/internal/pkg/redis"
	_config "Aurora/internal/tools/config"
)

type Config struct {
	_config.BaseConfig `yaml:",inline"`
	//Services _config.GrpcMap `yaml:"services"`
	Etcd  _config.Etcd  `yaml:"etcd"`
	Log   _log.Config   `yaml:"log"`
	Mysql _mysql.Config `yaml:"mysql"`
	Redis _redis.Config `yaml:"redis"`
	Jwt   jwt.Config    `yaml:"jwt"`
}

type Server struct {
	Cfg        Config
	Addr       string
	UserServer *grpc.Server
	SvcCtx     *svc.ServerCtx
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	// register in etcd
	etcdRegister := discovery.NewRegister([]string{s.Cfg.Etcd.Address}, s.SvcCtx.Logger)
	defer etcdRegister.Stop()

	node := discovery.Server{
		Name: s.Cfg.Name,
		Addr: s.Addr,
	}

	// register in etcd
	if _, err := etcdRegister.Register(node, 10); err != nil {
		return err
	}

	// init jwt config
	jwt.InitJWTConfig(&s.Cfg.Jwt)

	if err := s.UserServer.Serve(lis); err != nil {
		s.SvcCtx.Logger.Error("User Service exit!")
	}

	return nil
}

func New(opts ...OptionFunc) (*Server, error) {

	s := &Server{}

	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	// init config
	err := _config.GetConfig(o.ConfigFilePath, &s.Cfg)
	if err != nil {
		return nil, err
	}

	// init log
	logger := _log.InitLogger(&s.Cfg.Log)

	s.SvcCtx = &svc.ServerCtx{
		Ctx:    context.Background(),
		Logger: logger,
	}

	// init mysql conn
	s.SvcCtx.DBClient, err = _mysql.NewMysql(&s.Cfg.Mysql)
	if err != nil {
		s.SvcCtx.Logger.Error("init db error")
		return nil, err
	}

	// init redis cli
	s.SvcCtx.RedisCli, err = _redis.NewRedis(&s.Cfg.Redis)
	if err != nil {
		s.SvcCtx.Logger.Error("init redis error")
		return nil, err
	}

	s.SvcCtx.Cache = make(map[string]svc.Item)

	// init grpc internal
	s.Addr = s.Cfg.Host + ":" + s.Cfg.Port
	s.UserServer = grpc.NewServer()
	userpb.RegisterUserServiceServer(s.UserServer, service.NewUserServer(s.SvcCtx))

	s.SvcCtx.Logger.Info("User Service Init...")

	return s, nil
}
