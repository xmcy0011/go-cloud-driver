package logics

import (
	"context"
	"database/sql"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/common"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/service"
)

type metadataLogic struct {
	metadataSvc service.MetadataService
}

func NewMetadataLogic(db *sql.DB, metadata interfaces.DBMetadata, closure interfaces.DBMetadataClosure) interfaces.MetadataLogic {
	return &metadataLogic{
		metadataSvc: service.NewMetadataService(db, metadata, closure),
	}
}

// BeginUpload 开始文件上传
func (m *metadataLogic) BeginUpload(ctx context.Context, req interfaces.BeginUploadReq) (*interfaces.BeginUploadRsp, error) {
	return nil, errors.New("unimplent")
}
func (m *metadataLogic) EndUpload(ctx context.Context, req interfaces.EndUploadReq) (*interfaces.EndUploadRsp, error) {
	return nil, errors.New("unimplent")
}

func (m *metadataLogic) CreateDir(ctx context.Context, req interfaces.CreateDirReq) (*interfaces.CreateDirRsp, error) {
	uid := ulid.Make().String()

	_, err := m.metadataSvc.QueryById(ctx, req.ParentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.WithCause(common.RestBadRequest, "parentId not found")
		}
		return nil, errors.WithStack(err)
	}

	err = m.metadataSvc.Create(ctx, interfaces.Metadata{
		ObjectId:   uid,
		ParentId:   req.ParentId,
		Name:       req.Name,
		ObjectType: int(common.ObjectTypeDir),
	})

	if err != nil {
		return nil, err
	}

	return &interfaces.CreateDirRsp{ObjectId: uid}, nil
}

func (m *metadataLogic) MoveDir(ctx context.Context, req interfaces.MoveDirReq) (*interfaces.MoveDirRsp, error) {
	err := m.metadataSvc.MoveDir(ctx, req.ObjectId, req.NewParentId)
	if err != nil {
		return nil, err
	}

	return &interfaces.MoveDirRsp{}, nil
}

func (m *metadataLogic) QuerySubTree(ctx context.Context, req interfaces.QuerySubTreeReq) (*interfaces.QuerySubTreeRsp, error) {
	// trees, err := m.metadataSvc.QuerySubTree(ctx, req.ObjectId)
	// if err != nil {
	// 	return nil, err
	// }

	rsp := &interfaces.QuerySubTreeRsp{}
	// rsp.SubTrees = metadatNodes2SubTrees(trees)
	return rsp, nil
}

func (m *metadataLogic) metadatNodes2SubTrees(subTree map[string]interfaces.SubTrees, closure []interfaces.MetadataNode, depth int) {

}

// func metadatNodes2SubTrees(subTree map[string]interfaces.SubTrees, nodes []interfaces.MetadataNode, depth int) {

// 	subTrees := make(map[string]interfaces.SubTrees, 0)

// 	subTree.SubTrees = make([]interfaces.SubTrees, 0)

// 	for _, node := range nodes {
// 		if node.Depth != depth {
// 			continue
// 		}

// 		subTree.SubTrees = append(subTree.SubTrees, interfaces.SubTrees{
// 			ObjectId:   node.ObjectId,
// 			ParentId:   node.ParentId,
// 			ObjectType: node.ObjectType,
// 			Name:       node.Name,
// 			Depth:      node.Depth,
// 		})
// 	}

// 	return metadatNodes2SubTrees()
// }
