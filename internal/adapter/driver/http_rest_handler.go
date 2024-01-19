package driver

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
	"go.uber.org/zap"
)

type CreateDirReq struct {
	ParentId string `json:"parent_id"`
	Name     string `json:"name"`
	Ondup    int    `json:"ondup"`
}

type CreateDirRes struct {
	ObjectId string `json:"object_id"`
}

type MoveDirReq struct {
	NewParentId string `json:"new_parent_id"`
}

type HttpRestHandler struct {
	log *zap.Logger

	metadata interfaces.MetadataLogic
}

func NewHttpRestHandler(metadata interfaces.MetadataLogic) HttpRestHandler {
	return HttpRestHandler{
		log:      interfaces.MustNewLogger(),
		metadata: metadata,
	}
}

func (h *HttpRestHandler) RegisterRouter(g *gin.Engine) {
	g.POST("/api/metadata/dirs", h.createDir)
	g.GET("/api/metadata/dirs/:objectId/sub-trees", h.listDirSubTrees)
	g.PUT("/api/metadata/dirs/:objectId/move", h.moveDir)
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
	if _, err := ulid.Parse(req.ParentId); err != nil {
		h.responseError(g, BadRequest(errors.New("bad ulid format with parent_id")))
		return
	}

	ctx := context.Background()
	param := interfaces.CreateDirReq{Name: req.Name, ParentId: req.ParentId, Ondup: req.Ondup}
	result, err := h.metadata.CreateDir(ctx, param)
	if err != nil {
		h.responseError(g, err)
	} else {
		g.Header("Location", result.ObjectId)
		res := CreateDirRes{ObjectId: result.ObjectId}
		h.responseOk(g, http.StatusCreated, res)
	}
}

func (h *HttpRestHandler) listDirSubTrees(g *gin.Context) {

}

func (h *HttpRestHandler) moveDir(g *gin.Context) {
	req := MoveDirReq{}
	if err := g.ShouldBindJSON(&req); err != nil {
		h.responseError(g, BadRequest(err))
		return
	}

	objectId, ok := g.Params.Get("objectId")
	if !ok {
		h.responseError(g, BadRequest(errors.New("url miss objectId field")))
		return
	}

	if _, err := ulid.Parse(objectId); err != nil {
		h.responseError(g, BadRequest(errors.New("invalid objectId")))
		return
	}

	if _, err := ulid.Parse(req.NewParentId); err != nil {
		h.responseError(g, BadRequest(errors.New("invalid new_parent_id")))
		return
	}

	ctx := context.Background()
	_, err := h.metadata.MoveDir(ctx, interfaces.MoveDirReq{ObjectId: objectId, NewParentId: req.NewParentId})
	if err != nil {
		h.responseError(g, err)
	} else {
		h.responseOk(g, http.StatusNoContent, nil)
	}
}
