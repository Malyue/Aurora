package gateway

import (
	"Aurora/internal/apps/gateway/router"
	"Aurora/internal/apps/gateway/svc"
	discovery "Aurora/internal/pkg/etcd"
	_grpc "Aurora/internal/pkg/grpc"
	_log "Aurora/internal/pkg/log"
	_config "Aurora/internal/tools/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/resolver"
)

type Config struct {
	_config.BaseConfig `yaml:",inline"`
	//Services           _config.GrpcMap `yaml:"services"`
	Etcd _config.Etcd `yaml:"etcd"`
	//Jwt  jwt.Config   `yaml:"jwt"`
	Log _log.Config `yaml:"log"`
	//Redis              _redis
}

type Server struct {
	opts   *Options
	Config Config `default:"config.yaml"`
	Router *gin.Engine
	*svc.ServerCtx
}

func (s *Server) Run() error {

	//srv := &http.Server{
	//	Addr:    fmt.Sprintf(":%s", s.Config.Port),
	//	Handler: s.Router,
	//}

	err := s.Router.Run(fmt.Sprintf(":%s", s.Config.Port))
	if err != nil {
		panic(err)
	}

	s.Logger.Info("Start gateway server success, port : %d", s.Config.Port)

	return nil

	// TODO set signal to quit

}

func New(opts ...OptionFunc) (*Server, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	// init config
	var cfg Config
	err := _config.GetConfig(o.ConfigFilePath, &cfg)
	if err != nil {
		return nil, err
	}
	logrus.Info("Init config success")

	logger := _log.InitLogger(&cfg.Log)

	// set etcd resolver
	etcdResolver := discovery.NewResolver([]string{cfg.Etcd.Address}, logger)
	resolver.Register(etcdResolver)
	defer etcdResolver.Close()

	logger.Info("Set etcd resolver success")
	// init Client
	userServer, err := _grpc.InitUserClient()
	if err != nil {
		logrus.Errorf("create user client err : %s", err)
	}

	// create svc ctx
	ctx := &svc.ServerCtx{
		UserServer: userServer,
		Logger:     logger,
	}

	// init router
	r := router.InitRouter(ctx)

	logger.Info("Gateway init...")

	return &Server{
		opts:      o,
		Config:    cfg,
		Router:    r,
		ServerCtx: ctx,
	}, nil
}
