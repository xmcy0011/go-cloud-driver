package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
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
}
