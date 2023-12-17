package user

import (
	userpb "Aurora/api/proto-go/user"
	"Aurora/internal/apps/user/service"
	discovery "Aurora/internal/pkg/etcd"
	_config "Aurora/internal/tools/config"
	"google.golang.org/grpc"
	"net"
)

type Config struct {
	_config.BaseConfig `yaml:",inline"`
	//Services _config.GrpcMap `yaml:"services"`
	Etcd _config.Etcd `yaml:"etcd"`
}

type Server struct {
	Cfg        Config
	UserServer userpb.UserServiceServer
}

func (s *Server) Run() error {
	//grpcAddr := s.Cfg.Services[_const.UserServiceName].Addr[0]
	addr := s.Cfg.Host + ":" + s.Cfg.Port
	// init grpc server
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, service.GetUserSvc())
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// init etcd
	// TODO add logger instance
	etcdRegister := discovery.NewRegister([]string{s.Cfg.Etcd.Address}, nil)
	defer etcdRegister.Stop()

	node := discovery.Server{
		Name: s.Cfg.Name,
		Addr: addr,
	}

	// register in etcd
	if _, err := etcdRegister.Register(node, 10); err != nil {
		return err
	}

	// TODO add log output

	if err := grpcServer.Serve(lis); err != nil {
		return err
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

	// init grpcServer

	return s, nil
}
