package main

import (
	_user "Aurora/internal/apps/user"
	"path/filepath"
)

const config = "cmd/user/config.yaml"

func main() {
	dir, err := filepath.Abs(filepath.Dir("."))
	filePath := filepath.Join(dir, config)
	if err != nil {
		panic(err)
	}

	server, err := _user.New(
		_user.WithConfig(filePath))
	if err != nil {
		panic(err)
	}

	if err := server.Run(); err != nil {
		panic(err)
	}
}
