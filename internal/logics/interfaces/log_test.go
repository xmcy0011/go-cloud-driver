package interfaces

import (
	"fmt"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
)

func TestNewZapLogger(t *testing.T) {
	l := MustNewLogger()

	// basic
	l.Debug("debug msg")
	l.Info("info msg")
	l.Warn("warn msg")
	l.Error("error msg")

	// format
	l.Debug("hellor", zap.String("name", "go-cloud-dirver"))
	l.Info(fmt.Sprintf("server ip=%s", "127.0.0.1"))

	// context
	l = l.With(zap.Int("userId", 9527))
	l.Info("recv msg")
	l = l.With(zap.String("requestId", ulid.Make().String()))
	l.Debug("user login request")

	time.Sleep(time.Microsecond * 100)
}
