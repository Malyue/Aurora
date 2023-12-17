package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func GetConfig(filepath string, cfg any) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &cfg)
}
