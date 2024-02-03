package ai_proxy

import (
	userpb "Aurora/api/proto-go/user"
	_log "Aurora/internal/pkg/log"
	_config "Aurora/internal/tools/config"
)

type Config struct {
	_config.BaseConfig `yaml:",inline"`
	Etcd               _config.Etcd `yaml:"etcd"`
	Log                _log.Config  `yaml:"log"`
}

type Server struct {
	opt    *Options
	Config Config `default:"config.yaml"`
	//Router *gin.Engine
	UserServer userpb.UserServiceClient
}

func Run() {
}

func New() *Server {
	var s *Server

	return s
}
