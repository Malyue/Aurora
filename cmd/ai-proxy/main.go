package main

import (
	"path/filepath"

	_aiProxy "Aurora/internal/apps/ai-proxy"
)

const config = "cmd/ai-proxy/config.yaml"

func main() {
	dir, err := filepath.Abs(filepath.Dir("."))
	filePath := filepath.Join(dir, config)
	if err != nil {
		panic(err)
	}

	server, err := _aiProxy.New(
		_aiProxy.WithConfig(filePath))
	if err != nil {
		panic(err)
	}

	if err := server.Run(); err != nil {
		panic(err)
	}
}
