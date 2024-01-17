package driver

import (
	"github.com/gin-gonic/gin"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
	"go.uber.org/zap"
)

type HttpRestHandler struct {
	log *zap.Logger

	metadata interfaces.MetadataLogic
}

func NewHttpRestHandler(metadata interfaces.MetadataLogic) HttpRestHandler {
	return HttpRestHandler{
		log: interfaces.MustNewLogger(),
	}
}

func (h HttpRestHandler) RegisterRouter(g *gin.Engine) {
	g.POST("/api/clouddriver/dirs", h.createDir)
}

func (h HttpRestHandler) createDir(c *gin.Context) {

}
