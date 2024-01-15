package main

import (
	_access "Aurora/internal/apps/access-server"
	"path/filepath"
)

const config = "cmd/access-server/config.yaml"

func main() {
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(dir, config)

	server, err := _access.New(
		_access.WithConfig(filePath))
	if err != nil {
		panic(err)
	}

	server.Run()
}
