package main

import (
	_gateway "Aurora/internal/apps/geteway"
	"os"
	"path/filepath"
)

const config = "config.yaml"

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(filepath.Dir(exePath), config)

	server, err := _gateway.New(
		_gateway.WithConfig(filePath))
	if err != nil {
		panic(err)
	}

	server.Run()
}
