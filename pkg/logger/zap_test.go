package logger

import (
	"fmt"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewZapLogger(t *testing.T) {
	develop, err := NewZapLogger("test", true)
	require.Equal(t, err, nil)

	develop.Info("hello world")
	develop.Warn("warn msg")

	prod, err := NewZapLogger("test", false)
	require.Equal(t, err, nil)

	prod.Info("hello world")
	prod.Warn("warn msg")

	// basic
	prod.Debug("debug msg")
	prod.Info("info msg")
	prod.Warn("warn msg")
	prod.Error("error msg")

	// format
	prod.Debug("hellor", zap.String("name", "go-cloud-dirver"))
	prod.Info(fmt.Sprintf("server ip=%s", "127.0.0.1"))

	// context
	l := prod.With(zap.Int("userId", 9527))
	l.Info("recv msg")
	l = l.With(zap.String("requestId", ulid.Make().String()))
	l.Debug("user login request")

	time.Sleep(time.Microsecond * 100)
}
