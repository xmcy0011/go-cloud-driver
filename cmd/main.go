package main

import (
	"flag"

	"github.com/xmcy0011/go-cloud-driver/internal/conf"
	"github.com/xmcy0011/go-cloud-driver/internal/server"
)

var (
	config = flag.String("conf", "config.yaml", "-conf fileName")
)

func main() {
	flag.Parse()

	c := conf.MustLoad(*config)
	httpServer := server.NewHttpServer(*c)
	httpServer.Run()
}
