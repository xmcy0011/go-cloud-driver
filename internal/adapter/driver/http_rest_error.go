package driver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/common"
	"go.uber.org/zap"
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

func (h *HttpRestHandler) responseError(g *gin.Context, err error) {
	h.log.Info("", zap.Error(errors.Cause(err)))

	log.Println(errors.Cause(err))

	if bizErr, ok := err.(common.BizError); ok {
		g.JSON(bizErr.StatusCode, HttpRestError{
			Code:  bizErr.Code,
			Cause: fmt.Sprintf("%s %s", bizErr.FileLine, bizErr.Cause),
		})
	} else {
		g.JSON(http.StatusInternalServerError, HttpRestError{
			Code:  common.RestServerInternalError,
			Cause: err.Error(),
		})
	}
}

func (h *HttpRestHandler) responseOk(g *gin.Context, httpCode int, data any) {
	g.JSON(httpCode, data)
}
