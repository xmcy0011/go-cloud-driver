package interfaces

import "context"

type BeginUploadReq struct{}
type BeginUploadRsp struct{}

type EndUploadReq struct{}
type EndUploadRsp struct{}

type CreateDirReq struct {
	ParentId string `json:"parent_id"` // 父目录位置
	Name     string `json:"name"`      // 文件夹名称
	Ondup    int    `json:"ondup"`     // 冲突处理。1:重名报错 2:自动重命名 3:覆盖
}
type CreateDirRsp struct {
	ObjectId string `json:"object_id"` // 创建目录的 objectId
}

type MoveDirReq struct {
	ObjectId    string // 要移动的目录 objectId
	NewParentId string // 新的目标父目录 objectId
}

type MoveDirRsp struct{}

type QuerySubTreeReq struct {
	ObjectId string
}

type QuerySubTreeRsp struct {
	SubTrees *MetadataTreeNode `json:"sub_trees"` // 子树列表
}

type MetadataTreeNode struct {
	ObjectId   string `json:"object_id"`   // 元数据id
	Name       string `json:"name"`        // 名称
	ObjectType int    `json:"object_type"` // 类型
	Depth      int    `json:"depth"`       //层级深度，从0开始

	Children []*MetadataTreeNode `json:"children,omitempty"` // 子树
}

type MetadataLogic interface {
	// BeginUpload 开始文件上传
	BeginUpload(ctx context.Context, req BeginUploadReq) (*BeginUploadRsp, error)
	EndUpload(ctx context.Context, req EndUploadReq) (*EndUploadRsp, error)
	// CreateDir 创建目录
	CreateDir(ctx context.Context, req CreateDirReq) (*CreateDirRsp, error)
	// MoveDir 移动目录
	MoveDir(ctx context.Context, req MoveDirReq) (*MoveDirRsp, error)
	// QuerySubTree 查询子树
	QuerySubTree(ctx context.Context, req QuerySubTreeReq) (*QuerySubTreeRsp, error)
}
