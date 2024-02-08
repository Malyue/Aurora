package ai_proxy

import (
	"Aurora/internal/apps/ai-proxy/svc"
	discovery "Aurora/internal/pkg/etcd"
	_grpc "Aurora/internal/pkg/grpc"
	_log "Aurora/internal/pkg/log"
	_config "Aurora/internal/tools/config"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net"
)

type Config struct {
	_config.BaseConfig `yaml:",inline"`
	Etcd               _config.Etcd `yaml:"etcd"`
	Log                _log.Config  `yaml:"log"`
	AiProxy            string       `yaml:"aiProxy"`
}

type Server struct {
	opt    *Options
	Config Config `default:"config.yaml"`
	Addr   string
	//Router *gin.Engine
	AiProxyServer *grpc.Server
	svcCtx        *svc.ServerCtx
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	etcdRegister := discovery.NewRegister([]string{s.Config.Etcd.Address}, s.svcCtx.Logger)
	defer etcdRegister.Stop()

	node := discovery.Server{
		Name: s.Config.Name,
		Addr: s.Addr,
	}

	if _, err := etcdRegister.Register(node, 10); err != nil {
		return err
	}

	if err := s.AiProxyServer.Serve(lis); err != nil {
		s.svcCtx.Logger.Error("AiProxy Server exit!")
	}

	return nil
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

	logger := _log.InitLogger(&cfg.Log)

	svcCtx := &svc.ServerCtx{
		Ctx:     context.Background(),
		Logger:  logger,
		AiProxy: cfg.AiProxy,
	}

	etcdResolver := discovery.NewResolver([]string{cfg.Etcd.Address}, logger)
	resolver.Register(etcdResolver)
	defer etcdResolver.Close()

	svcCtx.UserClient, err = _grpc.InitUserClient()
	if err != nil {
		logrus.Errorf("create user client err : %s", err)
	}

	addr := cfg.Host + ":" + cfg.Port
	aiProxy := grpc.NewServer()
	// TODO register

	return &Server{
		Config:        cfg,
		Addr:          addr,
		AiProxyServer: aiProxy,
	}, nil
}
