package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func GetConfig(filepath string, cfg interface{}) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	return err
}
