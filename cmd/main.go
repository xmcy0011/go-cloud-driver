package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xmcy0011/go-cloud-driver/internal/adapter/driver"
	"github.com/xmcy0011/go-cloud-driver/internal/common"
	"github.com/xmcy0011/go-cloud-driver/internal/conf"
	"github.com/xmcy0011/go-cloud-driver/internal/dbaccess"
	"github.com/xmcy0011/go-cloud-driver/internal/infra/mysqldb"
	"github.com/xmcy0011/go-cloud-driver/internal/logics"
	"github.com/xmcy0011/go-cloud-driver/pkg/logger"
	"go.uber.org/zap"
)

var (
	config = flag.String("conf", "../config/config.yaml", "-conf fileName")
)

func main() {
	flag.Parse()

	config := conf.MustLoad(*config)

	logger, err := logger.NewZapLogger("go-cloud-driver", os.Getenv("env") == "prod")
	if err != nil {
		panic(err)
	}
	common.SetLogger(logger)

	engine := gin.Default()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	// 出站适配器
	myDb := mysqldb.MustInit(config.Db, logger)
	metadata := dbaccess.NewMetdata(myDb)
	metadataClosure := dbaccess.NewMetadataClosure(myDb)

	// 逻辑层依赖注入
	logics.SetDBPool(myDb)
	logics.SetDBMetadata(metadata)
	logics.SetDBMetadataClosure(metadataClosure)

	// 入栈适配器
	restHandler := driver.NewHttpRestHandler(logics.NewMetadataLogic())
	restHandler.RegisterRouter(engine)

	addr := fmt.Sprintf("%s:%d", config.Server.Listen, config.Server.Port)
	logger.Info("http server running", zap.String("addr", addr))
	if err := engine.Run(addr); err != nil {
		panic(err)
	}
}
