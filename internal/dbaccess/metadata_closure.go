package dbaccess

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
		"select t.ancestor,'%s',t.depth+1 from metadata_closure as t " + // 2. 把查询出的行的后代改为要插入的节点 id
		"where t.descendant = '%s' " + // 1. 查询出后代是 B3 的所有行
		"union all select '%s','%s',0;" // 3. 加上节点本身，深度为1"

	r, err := tx.ExecContext(ctx, fmt.Sprintf(sql, descendant, ancestor, descendant, descendant))
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

func (m *metadataClosure) Delete(ctx context.Context, ancestor string, tx *sql.Tx) (rowsAffected int64, err error) {
	sql := "delete from metadata_closure where descendant in (select a.id from (select descendant as id from metadata_closure where ancestor='%s') as a )"
	r, err := tx.ExecContext(ctx, fmt.Sprintf(sql, ancestor))
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func (m *metadataClosure) MoveSubTree(ctx context.Context, objectId, parentId string, tx *sql.Tx) (deleteCount, insertCount int64, err error) {
	var (
		rows *sql.Rows  = nil
		row  sql.Result = nil
	)

	// 第一步：先断开 x 这个子树和祖先们的关系，x 变成孤立的树
	// 如下目录结构：
	// test
	// |-a
	//   |-1.txt
	// |-b
	// 1）查询 a 的所有祖先：SELECT ancestor FROM metadata_closure WHERE descendant='%s' AND ancestor != descendant
	// 2）查询 a 的所有后代：SELECT descendant FROM metadata_closure WHERE ancestor='%s'
	// 3）组合2种情况：删除 a 和所有祖先的关系，以及 a 子树和所有祖先的关系。

	sql := `select id,ancestor,descendant,depth from metadata_closure 
	where descendant IN (SELECT descendant FROM 
						(SELECT descendant FROM metadata_closure WHERE ancestor='%s') as d)	--  后代节点(包括自己)
		AND ancestor IN (SELECT ancestor FROM 
						(SELECT ancestor FROM metadata_closure WHERE descendant='%s' AND ancestor != descendant) as a)	-- 祖先节点，不包括自己`

	rows, err = tx.QueryContext(ctx, sql, objectId, objectId)
	if err != nil {
		return
	}
	defer rows.Close()

	deletedClosure := make([]*interfaces.MetadataClosure, 0)
	for rows.Next() {
		item := &interfaces.MetadataClosure{}
		err = rows.Scan(&item.Id, &item.Ancestor, &item.Descendant, &item.Depth)
		if err != nil {
			return
		}
		deletedClosure = append(deletedClosure, item)
	}

	// 删除和祖先的关系
	args := make([]interface{}, 0, len(deletedClosure))
	sqlPlacehoder := ""
	for _, item := range deletedClosure {
		sqlPlacehoder = "?,"
		args = append(args, item.Id)
	}
	sql = fmt.Sprintf(`delete from metadata_closure where id in(%s)`, strings.Trim(sqlPlacehoder, ","))
	row, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return
	}
	deleteCount, err = row.RowsAffected()
	if err != nil {
		return
	}

	// 插入到新的路径
	// 1) 先找到目标路径的所有祖先
	// 2）每个节点都需要建立和这些祖先的关系，其中深度是 parentId 的深度 + 节点之前的深度
	// sql = `select ancestor,depth from metadata_closure where descendant=?`
	// rows, err = tx.QueryContext(ctx, sql, parentId)
	// if err != nil {
	// 	return
	// }

	// ancestors := make([]*interfaces.MetadataClosure, 0)
	// defer rows.Close()
	// for rows.Next() {
	// 	item := &interfaces.MetadataClosure{}
	// 	err = rows.Scan(&item.Ancestor, &item.Depth)
	// 	if err != nil {
	// 		return
	// 	}
	// 	ancestors = append(ancestors, item)
	// }

	// sql = `INSERT INTO metadata_closure(ancestor,descendant,depth) values()`
	// args = make([]interface{}, 0, len(deletedClosure))
	// sqlPlacehoder = ""
	// for _, item := range deletedClosure {
	// 	args = append(args, item)
	// }

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
		T1.descendant='%s' AND T2.ancestor='%s'`
	r, err := tx.ExecContext(ctx, fmt.Sprintf(sql, parentId, objectId))
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

func (m *metadataClosure) QueryCountByAncestor(ctx context.Context, ancestor string) (int, error) {
	sql := `select count(1) from metadata_closure where ancestor=?`
	row := m.db.QueryRow(sql, ancestor)
	count := 0
	err := row.Scan(&count)
	return count, err
}

func (m *metadataClosure) QueryExistByPairTx(ctx context.Context, closures []interfaces.MetadataClosure, tx *sql.Tx) (bool, error) {
	if len(closures) == 0 {
		return false, nil
	}

	sql := "select count(1) from metadata_closure where"
	args := make([]any, 0)
	for _, value := range closures {
		sql += " ancestor=? and descendant=?"
		args = append(args, value.Ancestor)
		args = append(args, value.Descendant)
	}
	row := tx.QueryRowContext(ctx, sql, args...)
	count := 0
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *metadataClosure) QueryExistByPair(ctx context.Context, closures []interfaces.MetadataClosure) (bool, error) {
	if len(closures) == 0 {
		return false, nil
	}

	sql := "select count(1) from metadata_closure where"
	args := make([]any, 0)
	for _, value := range closures {
		sql += " ancestor=? and descendant=?"
		args = append(args, value.Ancestor)
		args = append(args, value.Descendant)
	}
	row := m.db.QueryRowContext(ctx, sql, args...)
	count := 0
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
