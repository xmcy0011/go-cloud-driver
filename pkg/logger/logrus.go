package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger 服务日志服务，可适配其他日志组件
type Logger interface {
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Tracef(format string, args ...interface{})
	Traceln(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
}

var (
	logOnce sync.Once
	l       *logrus.Logger
)

// NewLogger 获取日志句柄
func NewLogger() *logrus.Logger {
	logOnce.Do(initLogger)

	return l
}

func initLogger() {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fileName := path.Base(f.File)
			funcName := path.Base(f.Function)
			arr := strings.Split(funcName, ".")
			if len(arr) > 1 {
				funcName = arr[len(arr)-1]
			}
			return funcName, fmt.Sprintf("%s:%d", fileName, f.Line)
		},
	})

	if len(os.Getenv("LOGOUT")) > 0 {
		logger.SetOutput(os.Stdout)
	} else {
		logDir := "/var/log/metastore/"
		logName := "document.log"
		var filePerm os.FileMode = 0750
		err := os.MkdirAll(logDir, filePerm)
		if err != nil {
			fmt.Println("mkdir err", err)
		}
		logFileName := path.Join(logDir, logName)
		logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
		if err != nil {
			fmt.Println("open log file err", err)
		}
		logger.SetOutput(logFile)
	}

	l = logger
}
