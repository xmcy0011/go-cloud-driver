package main

import (
	"flag"

	"github.com/xmcy0011/go-cloud-driver/internal/conf"
	"github.com/xmcy0011/go-cloud-driver/internal/server"
)

var (
	config = flag.String("conf", "../config/config.yaml", "-conf fileName")
)

func main() {
	flag.Parse()

	// go func() {
	// 	r := gin.Default()
	// 	r.GET("/ping", func(c *gin.Context) {
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"message": "pong",
	// 		})
	// 	})
	// 	r.Run(":9800")
	// }()

	c := conf.MustLoad(*config)
	httpServer := server.NewHttpServer(*c)
	httpServer.Run()
}
