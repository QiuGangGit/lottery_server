package main

import (
	"flag"
	"lottery_server/internal/config"
	"lottery_server/internal/server"
	"lottery_server/internal/svc"
)

var configFile = flag.String("config", "config.json", "config file")

func main() {
	flag.Parse()

	cfg := config.MustLoad(*configFile)

	c := svc.NewServiceContext(cfg)

	s := server.NewHttpServer(c)

	s.Start()
}
