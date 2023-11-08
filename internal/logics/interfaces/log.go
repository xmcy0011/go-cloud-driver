package interfaces

import (
	"os"

	"github.com/xmcy0011/go-cloud-driver/pkg/logger"
	"go.uber.org/zap"
)

const serviceName = "go-cloud-driver"

// MustNewLogger 创建logger实例
func MustNewLogger() *zap.Logger {
	isDevelopment := true
	if os.Getenv("env") == "prod" {
		isDevelopment = false
	}

	l, err := logger.NewZapLogger(serviceName, isDevelopment)
	if err != nil {
		panic(err)
	}
	return l
}
