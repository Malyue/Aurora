package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

func GetConfig(filepath string, cfg interface{}) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &cfg)

}