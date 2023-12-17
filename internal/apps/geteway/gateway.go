package geteway

import (
	userpb "Aurora/api/proto-go/user"
	"Aurora/internal/apps/geteway/router"
	discovery "Aurora/internal/pkg/etcd"
	_grpc "Aurora/internal/pkg/grpc"
	_config "Aurora/internal/tools/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	Logger *logrus.Logger
	//Etcd   *etcd
	// Grpc Server
	UserServer *userpb.UserServiceClient
}

func (s *Server) Run() error {

	svc := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Config.Port),
		Handler: s.Router,
	}

	return svc.ListenAndServe()

	//c := make(chan os.Signal, 1)
	//signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	//go func() {
	//	defer func() {
	//		logrus.Println("Shutting down server ...")
	//
	//		timeCtx, timeCancel := context.WithTimeout(context.Background(), 3*time.Second)
	//
	//		if err := s.Shutdown(timeCtx); err != nil {
	//			logrus.Fatalf("HTTP Server Shutdown Err: %s", err)
	//		}
	//	}()
	//
	//	//select {
	//	//case <-ctx.Done():
	//	//	return ctx
	//	//
	//	//}
	//}()
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

	// init router
	r := router.InitRouter()

	// TODO add log instance
	register := discovery.NewResolver([]string{cfg.Etcd.Address}, nil)
	resolver.Register(register)
	defer register.Close()
	userService := _grpc.InitUserClient()

	return &Server{
		opts:       o,
		Config:     cfg,
		Router:     r,
		UserServer: &userService,
	}, nil
}

func init() {

}
