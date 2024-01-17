package logics

import (
	"context"
	"database/sql"
	"errors"

	"github.com/oklog/ulid/v2"
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

// CreateDir 创建目录
func (m *metadataLogic) CreateDir(ctx context.Context, req interfaces.CreateDirReq) (*interfaces.CreateDirRsp, error) {
	uid := ulid.Make().String()

	_, err := m.metadataSvc.QueryById(ctx, req.ParentId)
	if err != nil {
		return nil, err
	}

	err = m.metadataSvc.Create(ctx, interfaces.Metadata{
		ObjectId:  uid,
		ParentId:  req.ParentId,
		Name:      req.Name,
		BasicAttr: int(common.BasicAttrDir),
	})

	if err != nil {
		return nil, err
	}

	return &interfaces.CreateDirRsp{ObjectId: uid}, nil
}
