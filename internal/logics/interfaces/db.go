package interfaces

import (
	"context"
	"database/sql"
)

type Metadata struct {
	ObjectId   string // 元数据id
	ParentId   string // 父对象id
	Name       string // 名称
	ObjectType int    // 类型
}

type MetadataClosure struct {
	Id         int    // 自增主键
	Ancestor   string // 祖先
	Descendant string // 后代
	Depth      int    //层级深度，从0开始
}

type MetadataNode struct {
	ObjectId   string // 元数据id
	ParentId   string // 父对象id
	Name       string // 名称
	ObjectType int    // 类型
	Ancestor   string // 祖先
	Descendant string // 后代
	Depth      int    //层级深度，从0开始
}

// IMetadata 元数据表操作
type DBMetadata interface {
	Add(ctx context.Context, meta Metadata, tx *sql.Tx) error
	QueryCountById(ctx context.Context, objectId string) (int, error)
	QueryById(ctx context.Context, objectId string) (*Metadata, error)
	UpdateParentId(ctx context.Context, objectId, newParentId string, tx *sql.Tx) (rowsAffected int64, err error)
}

// DBMetadataClosure 元数据闭包关系表
type DBMetadataClosure interface {
	// Add: 添加元数据关系
	// ancestor: 祖先节点，在那个节点下插入，即为那个节点
	// descendant: 后代节点，对象本身
	Add(ctx context.Context, ancestor, descendant string, tx *sql.Tx) (rowsAffected int64, err error)
	// Delete: 删除某个祖先的所有后代节点
	Delete(ctx context.Context, ancestor string, tx *sql.Tx) (rowsAffected int64, err error)
	// MoveSubTree：移动子树到另外一个节点
	MoveSubTree(ctx context.Context, objectId, parentId string, tx *sql.Tx) (deleteCount, insertCount int64, err error)
	// QueryAllDescendants: 查询所有后代（包含自己，其深度为0），按照节点深度升序排序
	QueryAllDescendants(ctx context.Context, ancestor string) ([]MetadataNode, error)
	// CheckIsDescendant: 检查某个节点是否是后代节点
	CheckIsDescendant(ctx context.Context, ancestor, descendant string, tx *sql.Tx) (bool, error)
}
