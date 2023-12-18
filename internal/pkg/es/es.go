package es

import (
	"fmt"
	"github.com/elastic/go-elasticsearch"
)

type Config struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Index string `yaml:"index"`
}

// InitEs init es and return *elasticsearch.Client,error
func InitEs(cfg []Config) (*elasticsearch.Client, error) {
	config := elasticsearch.Config{}
	for _, c := range cfg {
		conn := fmt.Sprintf("%s:%s", c.Host, c.Port)
		config.Addresses = append(config.Addresses, conn)
	}
	return elasticsearch.NewClient(config)
}
