package interfaces

import (
	"context"
	"database/sql"
)

type Metadata struct {
	ObjectId  string // 元数据id
	ParentId  string // 父对象id
	Name      string // 名称
	BasicAttr int    // 类型
}

type MetadataClosure struct {
	Id         int    // 自增主键
	Ancestor   string // 祖先
	Descendant string // 后代
	Depth      int    //层级深度，从0开始
}

// IMetadata 元数据表操作
type DBMetadata interface {
	Add(ctx context.Context, meta Metadata, tx *sql.Tx) error
	QueryCountById(ctx context.Context, objectId string) (int, error)
	QueryById(ctx context.Context, objectId string) (*Metadata, error)
}

// DBMetadataClosure 元数据闭包关系表
type DBMetadataClosure interface {
	Add(ctx context.Context, ancestor, descendant string, tx *sql.Tx) (int, error)
}
