package driver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmcy0011/go-cloud-driver/pkg/bizerr"
)

type HttpRestError struct {
	// Code: 业务错误码
	Code string `json:"code"`
	// Cause: 错误原因
	Cause string `json:"cause"`
	// Description：错误描述
	//Description string
	// Solution: 解决方案
	//Solution string
}

const (
	RestBadRequest          string = "bad.request"
	RestServerInternalError string = "server.internal.error"
)

func BadRequest(err error) error {
	return bizerr.WithCause(http.StatusBadRequest, RestBadRequest, err.Error())
}

func (h *HttpRestHandler) responseError(g *gin.Context, err error) {
	if bizErr, ok := err.(bizerr.Error); ok {
		g.JSON(bizErr.StatusCode, HttpRestError{
			Code:  bizErr.Code,
			Cause: bizErr.Cause,
		})
	} else {
		g.JSON(http.StatusInternalServerError, HttpRestError{
			Code:  RestServerInternalError,
			Cause: err.Error(),
		})
	}
}
