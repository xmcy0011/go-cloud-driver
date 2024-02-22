package common

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type BizError struct {
	// HTTP 状态码
	StatusCode int
	// Code: 业务错误码
	Code string
	// Cause: 错误原因
	Cause string
	// fileLine: 错误行号
	FileLine string
}

const (
	RestBadRequest          string = "bad.request"
	RestServerInternalError string = "server.internal.error"
)

var statusCodeMap = map[string]int{
	RestBadRequest:          http.StatusOK,
	RestServerInternalError: http.StatusInternalServerError,
}

func WithCause(code, cause string) error {
	_, file, line, ok := runtime.Caller(1)
	fileLine := ""
	if ok {
		pos := strings.LastIndex(file, "/")
		if pos != -1 {
			file = file[pos+1:]
		}
		fileLine = fmt.Sprintf("[%s:%d]", file, line)
	}

	httpStateCode := http.StatusInternalServerError
	if v, ok := statusCodeMap[code]; ok {
		httpStateCode = v
	}

	return BizError{StatusCode: httpStateCode, Code: code, Cause: cause, FileLine: fileLine}
}

func (b BizError) Error() string {
	return fmt.Sprintf("%s statusCode: %d, code: %s, cause: %s", b.FileLine, b.StatusCode, b.Code, b.Cause)
}
