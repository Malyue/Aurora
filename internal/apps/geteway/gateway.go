package geteway

import (
	_config "Aurora/internal/tools/config"
)

type Config struct {
	_config.BaseConfig
}

type Server struct {
	opts   *Options
	Config *Config `default:"config.yaml"`
}

func (s *Server) Run() error {

	return nil
}

func New(opts ...OptionFunc) (*Server, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	// init config
	var cfg *Config
	err := _config.GetConfig(o.ConfigFilePath, cfg)
	if err != nil {
		return nil, err
	}

	// init router

	return &Server{
		opts:   o,
		Config: cfg,
	}, nil
}

func init() {

}
