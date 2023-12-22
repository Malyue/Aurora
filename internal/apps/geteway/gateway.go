package geteway

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/resolver"

	"Aurora/internal/apps/geteway/router"
	"Aurora/internal/apps/geteway/svc"
	discovery "Aurora/internal/pkg/etcd"
	_grpc "Aurora/internal/pkg/grpc"
	"Aurora/internal/pkg/jwt"
	_log "Aurora/internal/pkg/log"
	_config "Aurora/internal/tools/config"
)

type Config struct {
	_config.BaseConfig `yaml:",inline"`
	Services           _config.GrpcMap `yaml:"services"`
	Etcd               _config.Etcd    `yaml:"etcd"`
	Jwt                jwt.Config      `yaml:"jwt"`
	Log                _log.Config     `yaml:"log"`
	//Redis              _redis
}

type Server struct {
	opts   *Options
	Config Config `default:"config.yaml"`
	Router *gin.Engine
	*svc.ServerCtx
}

func (s *Server) Run() error {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Config.Port),
		Handler: s.Router,
	}

	return srv.ListenAndServe()

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

	logger := _log.InitLogger(&cfg.Log)

	// set etcd resolver
	etcdResolver := discovery.NewResolver([]string{cfg.Etcd.Address}, logger)
	resolver.Register(etcdResolver)
	defer etcdResolver.Close()

	// init Client
	userServer := _grpc.InitUserClient()

	// create svc ctx
	ctx := &svc.ServerCtx{
		UserServer: userServer,
		Logger:     logger,
	}

	// init jwt config
	jwt.InitJWTConfig(&cfg.Jwt)

	// init router
	r := router.InitRouter(ctx)

	return &Server{
		opts:      o,
		Config:    cfg,
		Router:    r,
		ServerCtx: ctx,
	}, nil
}
