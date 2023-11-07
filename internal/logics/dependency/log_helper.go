package dependency

import (
	"github.com/xmcy0011/go-cloud-driver/pkg/logger"
	"go.uber.org/zap"
)

type CloudLog interface {
}

// MustNewLogger 创建logger实例
func MustNewLogger() *zap.Logger {
	l, err := logger.NewZapLogger("go-cloud-driver", false)
	if err != nil {
		panic(err)
	}
	return l
}
