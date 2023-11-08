package logics

import (
	"context"
	"errors"

	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

type metadataService struct {
	metadata interfaces.DBMetadata
	closure  interfaces.DBMetadataClosure
}

func NewMetadataServer(metadata interfaces.DBMetadata, closure interfaces.DBMetadataClosure) interfaces.MetadataServer {
	return &metadataService{
		metadata: metadata,
		closure:  closure,
	}
}

// BeginUpload 开始文件上传
func (m *metadataService) BeginUpload(ctx context.Context, req interfaces.BeginUploadReq) (interfaces.BeginUploadRsp, error) {
	return interfaces.BeginUploadRsp{}, errors.New("unimplent")
}
func (m *metadataService) EndUpload(ctx context.Context, req interfaces.EndUploadReq) (interfaces.EndUploadRsp, error) {
	return interfaces.EndUploadRsp{}, errors.New("unimplent")
}

// CreateDir 创建目录
func (m *metadataService) CreateDir(ctx context.Context, req interfaces.CreateDirReq) (interfaces.CreateDirRsp, error) {

}
