package errhelper

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// WithFileLine: 给错误增加堆栈
func WithFileLine(err error) error {
	if err == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		pos := strings.LastIndex(file, "/")
		if pos != -1 {
			file = file[pos+1:]
		}
		return errors.Wrap(err, fmt.Sprintf("[%s:%d]", file, line))
	}
	return err
}
