package driver

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
	"go.uber.org/zap"
)

type CreateDirReq struct {
	ParentId string `json:"parent_id"`
	Name     string `json:"name"`
	Ondup    int    `json:"ondup"`
}

type HttpRestHandler struct {
	log *zap.Logger

	metadata interfaces.MetadataLogic
}

func NewHttpRestHandler(metadata interfaces.MetadataLogic) HttpRestHandler {
	return HttpRestHandler{
		log: interfaces.MustNewLogger(),
	}
}

func (h *HttpRestHandler) RegisterRouter(g *gin.Engine) {
	g.POST("/api/metadata/dirs", h.createDir)
}

func (h *HttpRestHandler) createDir(g *gin.Context) {
	req := CreateDirReq{}
	if err := g.ShouldBindJSON(&req); err != nil {
		h.responseError(g, BadRequest(err))
		return
	}

	if req.Name == "" {
		h.responseError(g, BadRequest(errors.New("invalid object name")))
		return
	}

	ctx := context.Background()
	param := interfaces.CreateDirReq{Name: req.Name, ParentId: req.ParentId, Ondup: req.Ondup}
	result, err := h.metadata.CreateDir(ctx, param)
	if err != nil {
		h.responseError(g, err)
	} else {
		g.Header("Location", result.ObjectId)
		g.JSON(http.StatusCreated, gin.H{})
	}
}
