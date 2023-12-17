package main

import (
	_gateway "Aurora/internal/apps/geteway"
	"path/filepath"
)

const config = "cmd/gateway/config.yaml"

func main() {
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(dir, config)

	server, err := _gateway.New(
		_gateway.WithConfig(filePath))
	if err != nil {
		panic(err)
	}

	server.Run()
}
