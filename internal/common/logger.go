package common

import "go.uber.org/zap"

var logger *zap.Logger

// SetLogger 设置 logger 实例
func SetLogger(l *zap.Logger) {
	logger = l
}

// GetLogger 获取 logger 实例
func GetLogger() *zap.Logger {
	return logger
}
