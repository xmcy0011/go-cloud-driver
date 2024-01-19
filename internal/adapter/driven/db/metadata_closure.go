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

func (m *metadataClosure) Add(ctx context.Context, ancestor, descendant string, tx *sql.Tx) (int64, error) {
	sql := "insert into metadata_closure(ancestor, descendant, depth) " +
		"select t.ancestor,?,t.depth+1 from metadata_closure as t " + // 2. 把查询出的行的后代改为要插入的节点 id
		"where t.descendant = ? " + // 1. 查询出后代是 B3 的所有行
		"union all select ?,?,0;" // 3. 加上节点本身，深度为1"

	r, err := tx.ExecContext(ctx, sql, descendant, ancestor, descendant, descendant)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

func (m *metadataClosure) Delete(ctx context.Context, ancestor string, tx *sql.Tx) (rowsAffected int64, err error) {
	sql := "delete from metadata_closure where descendant in (select a.id from (select descendant as id from metadata_closure where ancestor=?) as a )"
	r, err := tx.ExecContext(ctx, sql, ancestor)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func (m *metadataClosure) MoveSubTree(ctx context.Context, objectId, parentId string, tx *sql.Tx) (deleteCount, insertCount int64, err error) {

	// 第一步：先断开 x 这个子树和祖先们的关系，x 变成孤立的树
	sql := `delete FROM metadata_closure
	WHERE descendant IN (SELECT descendant FROM
						(SELECT ancestor,descendant FROM metadata_closure) as d WHERE d.ancestor=?)	--  后代节点(包括自己)
	AND ancestor IN (SELECT ancestor FROM 
					(SELECT ancestor,descendant FROM metadata_closure) as a WHERE a.descendant=? AND a.ancestor != a.descendant)	-- 祖先节点，不包括自己
	`
	r, err := tx.ExecContext(ctx, sql, objectId, objectId)
	if err != nil {
		return
	}

	deleteCount, err = r.RowsAffected()
	if err != nil {
		return
	}

	// 第二步：将上一步分离出的子树用笛卡尔积嫁接到节点1下
	sql = `INSERT INTO metadata_closure(ancestor,descendant,depth)
	SELECT
		T1.ancestor,
		T2.descendant,
		T1.depth + T2.depth + 1 as depth
	FROM
		metadata_closure as T1
		CROSS JOIN
		metadata_closure as T2
	WHERE
		T1.descendant=? AND T2.ancestor=?`
	r, err = tx.ExecContext(ctx, sql, parentId, objectId)
	if err != nil {
		return
	}

	insertCount, err = r.RowsAffected()
	if err != nil {
		return
	}

	return
}

func (m *metadataClosure) CheckIsDescendant(ctx context.Context, ancestor, descendant string, tx *sql.Tx) (bool, error) {
	sql := "select count(1) from metadata_closure where ancestor=? and descendant=?"
	r := tx.QueryRowContext(ctx, sql, ancestor, descendant)
	count := 0
	err := r.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *metadataClosure) QueryAllDescendants(ctx context.Context, ancestor string) ([]interfaces.MetadataNode, error) {
	sql := `select T1.ancestor, T1.descendant, T1.depth, T2.object_id, T2.parent_id, T2.name, T2.object_type from metadata_closure as T1
			inner join metadata as T2
			on T1.descendant=T2.object_id
			where T1.ancestor=? order by T1.depth`
	rows, err := m.db.Query(sql, ancestor)
	if err != nil {
		return nil, err
	}

	result := make([]interfaces.MetadataNode, 0)

	for rows.Next() {
		info := interfaces.MetadataNode{}
		err = rows.Scan(&info.Ancestor, &info.Descendant, &info.Depth, &info.ObjectId, &info.ParentId, &info.Name, &info.ObjectType)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	return result, nil
}
