package db

import (
	"context"
	"database/sql"

	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

type metadataClosure struct {
	db *sql.DB
}

func NewMetadataClosure(db *sql.DB) interfaces.DBMetadataClosure {
	return &metadataClosure{db: db}
}

func (m *metadataClosure) Add(ctx context.Context, ancestor, descendant string, tx *sql.Tx) (int, error) {
	sql := "insert into metadata_closure(ancestor, descendant,depth)" +
		"select t.ancestor, ?,t.depth+1 from metadata_path as t" + // 2. 把查询出的行的后代改为要插入的节点 id
		"where t.descendant = ?" + // 1. 查询出后代是 B3 的所有行
		"union all select ?,?,0;" // 3. 加上节点本身，深度为1"

	r, err := tx.Exec(sql, descendant, ancestor, descendant, descendant)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
