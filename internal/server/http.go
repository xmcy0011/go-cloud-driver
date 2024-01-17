package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xmcy0011/go-cloud-driver/internal/adapter/driven/db"
	"github.com/xmcy0011/go-cloud-driver/internal/adapter/driver"
	"github.com/xmcy0011/go-cloud-driver/internal/conf"
	"github.com/xmcy0011/go-cloud-driver/internal/infra/mysqldb"
	"github.com/xmcy0011/go-cloud-driver/internal/logics"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
	"go.uber.org/zap"
)

type HttpServer struct {
	config      conf.Config
	engine      *gin.Engine
	log         *zap.Logger
	restHandler driver.HttpRestHandler
}

func NewHttpServer(config conf.Config) *HttpServer {
	h := &HttpServer{
		engine: gin.Default(),
		config: config,
		log:    interfaces.MustNewLogger(),
	}

	h.engine.Use(gin.Recovery())
	h.engine.Use(gin.Logger())

	// db
	myDb := mysqldb.MustInit(config.Db, h.log)

	// 出站适配器
	metadata := db.NewMetdata(myDb)
	metadataClosure := db.NewMetadataClosure(myDb)

	// 逻辑层
	metadataLogic := logics.NewMetadataLogic(myDb, metadata, metadataClosure)

	// 入栈适配器
	h.restHandler = driver.NewHttpRestHandler(metadataLogic)

	return h
}

func (h *HttpServer) Run() {
	addr := fmt.Sprintf("%s:%d", h.config.Server.Listen, h.config.Server.Port)
	h.log.Info("http server running", zap.String("addr", addr))
	h.engine.Run(addr)
}
