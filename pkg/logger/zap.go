package logger

import "go.uber.org/zap"

// NewZapLogger new zap logger
func NewZapLogger(serviceName string, development bool) (*zap.Logger, error) {
	var config zap.Config
	var encoding = "json"
	if development {
		config = zap.NewDevelopmentConfig()
		encoding = "console"
	} else {
		config = zap.NewProductionConfig()
	}
	config.Encoding = encoding

	l, err := config.Build(
		//zap.AddCaller(),
		//zap.AddCallerSkip(callerSkip), //解决kratos 使用zap后 堆栈不正确的问题
		zap.Fields(
			zap.String("app", serviceName),
		),
	)

	return l, err
}
