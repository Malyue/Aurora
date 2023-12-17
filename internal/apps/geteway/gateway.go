package geteway

import (
	"Aurora/internal/apps/geteway/router"
	"Aurora/internal/apps/geteway/svc"
	discovery "Aurora/internal/pkg/etcd"
	_grpc "Aurora/internal/pkg/grpc"
	_config "Aurora/internal/tools/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/resolver"
	"net/http"
)

type Config struct {
	_config.BaseConfig `yaml:",inline"`
	Services           _config.GrpcMap `yaml:"services"`
	Etcd               _config.Etcd    `yaml:"etcd"`
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

	// TODO signal

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

	// TODO init log

	// TODO add log instance
	register := discovery.NewResolver([]string{cfg.Etcd.Address}, nil)
	resolver.Register(register)
	defer register.Close()
	userService := _grpc.InitUserClient()

	ctx := &svc.ServerCtx{
		UserServer: userService,
	}

	// init router
	r := router.InitRouter(ctx)

	return &Server{
		opts:      o,
		Config:    cfg,
		Router:    r,
		ServerCtx: ctx,
	}, nil
}

func init() {

}
