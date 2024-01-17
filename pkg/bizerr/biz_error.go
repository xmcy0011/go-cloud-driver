package bizerr

import "fmt"

type Error struct {
	// HTTP 状态码
	StatusCode int
	// Code: 业务错误码
	Code string
	// Cause: 错误原因
	Cause string
}

func WithCause(statusCode int, code, cause string) error {
	return Error{StatusCode: statusCode, Code: code, Cause: cause}
}

func (b Error) Error() string {
	return fmt.Sprintf("statusCode: %d, code: %s, cause: %s", b.StatusCode, b.Code, b.Cause)
}
