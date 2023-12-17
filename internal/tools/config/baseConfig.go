package config

type BaseConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type GrpcMap map[string]*Service

type Service struct {
	Name        string   `yaml:"name"`
	LoadBalance bool     `yaml:"loadBalance"`
	Addr        []string `yaml:"addr"`
}

type Etcd struct {
	Address string `yaml:"address"`
}

type Kafka struct {
	Address []string `yaml:"address"`
}
