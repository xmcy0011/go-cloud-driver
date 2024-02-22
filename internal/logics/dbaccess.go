package logics

import (
	"database/sql"

	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

var dbPool *sql.DB
var dbMetadata interfaces.DBMetadata
var dbMetadataClosure interfaces.DBMetadataClosure

// SetDBPool 注入 DB 连接池对象
func SetDBPool(pool *sql.DB) {
	dbPool = pool
}

// SetDBMetadata 注入 DBMetadata 依赖
func SetDBMetadata(metadata interfaces.DBMetadata) {
	dbMetadata = metadata
}

// SetDBMetadataClosure 注入 DBMetadataClosure 依赖
func SetDBMetadataClosure(closure interfaces.DBMetadataClosure) {
	dbMetadataClosure = closure
}
