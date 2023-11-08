package interfaces

import "context"

type BeginUploadReq struct{}
type BeginUploadRsp struct{}

type EndUploadReq struct{}
type EndUploadRsp struct{}

type CreateDirReq struct {
	ObjectId string `json:"object_id"` // 父目录 objectId
	Name     string `json:"name"`      // 文件夹名称
	Ondup    string `json:"ondup"`     // 冲突处理。1:重名报错 2:自动重命名 3:覆盖
}
type CreateDirRsp struct {
	ObjectId string `json:"object_id"` // 创建目录的 objectId
}

type MetadataServer interface {
	// BeginUpload 开始文件上传
	BeginUpload(ctx context.Context, req BeginUploadReq) (BeginUploadRsp, error)
	EndUpload(ctx context.Context, req EndUploadReq) (EndUploadRsp, error)
	// CreateDir 创建目录
	CreateDir(ctx context.Context, req CreateDirReq) (CreateDirRsp, error)
}
