package main

import (
	"musiclib/internal/config"
	"musiclib/internal/service"
)

func main() {
	cfg := config.GetConfig()
	service, err := service.NewService(cfg)
	if err != nil {
		return
	}
	service.Run()
}
