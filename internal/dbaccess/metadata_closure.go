package dbaccess

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/xmcy0011/go-cloud-driver/internal/common"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
	"go.uber.org/zap"
)

type metadataClosure struct {
	db  *sql.DB
	log *zap.Logger
}

func NewMetadataClosure(db *sql.DB) interfaces.DBMetadataClosure {
	return &metadataClosure{db: db, log: common.GetLogger()}
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
		row sql.Result = nil
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
	sql := fmt.Sprintf(`delete from metadata_closure 
				where descendant IN (SELECT descendant FROM (SELECT descendant FROM metadata_closure WHERE ancestor='%s') as d)	--  后代节点(包括自己)
					  AND ancestor IN (SELECT ancestor FROM (SELECT ancestor FROM metadata_closure WHERE descendant='%s' AND ancestor != descendant) as a)	-- 祖先节点，不包括自己)`,
		objectId, objectId)
	if row, err = tx.ExecContext(ctx, sql); err != nil {
		return
	}
	if deleteCount, err = row.RowsAffected(); err != nil {
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

func (m *metadataClosure) MoveLargeSubTree(ctx context.Context, objectId, parentId string, pageSize int) (deleteCount, insertCount int64, err error) {
	var (
		rows   *sql.Rows  = nil
		result sql.Result = nil
		count  int64      = 0
	)

	// 第一步：先断开 x 这个子树和祖先们的关系，x 变成孤立的树
	// 如下目录结构：
	// test
	// |-a
	//   |- c
	//   	|- 2.txt
	// |-b
	// 按照路径枚举（深度降序）总结其闭包如下：
	// test/a/c/2.txt => (2.txt, 2.txt, 0), (2.txt, c, 1), (2.txt, a, 2), (2.txt, test, 3)
	// test/a/c       => (c, c, 0), (c, a, 1), (c, test, 2)
	// test/a         => (a, a, 0), (a, test, 1)
	// test/b         => (b, b, 0), (b, test, 1)
	// test           => (test, test, 0)
	// 移动就是删除前缀和后代的笛卡尔积组合，把 c 移动到 b 下，则需要删除 test/a 前缀，把 c 变成孤立的树：
	// c/2.txt        => (2.txt, 2.txt, 0), (2.txt, c, 1)
	// c              => (c, c, 0)
	// 所以，SQL：
	// 1）先查询 c 的所有祖先：SELECT ancestor FROM metadata_closure WHERE descendant='c' AND ancestor != descendant
	// 2）再查询 c 的所有后代：SELECT descendant FROM metadata_closure WHERE ancestor='c'
	// 3）取笛卡尔积(n * n)：删除 c 自己以及所有后代和祖先的关系，把 c 变成孤立的树

	//
	// 考虑到超大目录移动超时的情况，我们需要循环从叶子节点自底向上删除，将其变成孤立的树，插入时自上而下插入
	//
	// 查询源的所有祖先，不包括自己
	sqlStr := "select ancestor from metadata_closure where descendant=? and ancestor != descendant"
	rows, err = m.db.QueryContext(ctx, sqlStr, objectId)
	if err != nil {
		return
	}
	srcParentPaths := make([]string, 0)
	for rows.Next() {
		var ancestor string
		if err = rows.Scan(&ancestor); err != nil {
			rows.Close()
			return
		}
		srcParentPaths = append(srcParentPaths, ancestor)
	}
	rows.Close()

	// 查询后代，最大深度
	sqlStr = "select max(depth) from metadata_closure where ancestor=?"
	row := m.db.QueryRowContext(ctx, sqlStr, objectId)
	maxDepth := 0
	if err = row.Scan(&maxDepth); err != nil {
		return
	}

	// 从叶子开始删除直到自己，如果此时没有查询到祖先，表明已经是一颗孤立的树
	alreadyIsolatedTree := len(srcParentPaths) == 0

	if !alreadyIsolatedTree {
		// depth=0 自己到祖先也需要删除
		for depth := maxDepth; depth >= 0; depth-- {
			m.log.Info("MoveLargeSubTree delete depth", zap.Int("depth", depth), zap.Int("pageSize", pageSize))
			lastDescendant := ""
			startDescendant := ""
			for {
				// 查询指定深度的后代的分页标记，从叶子向上分页删除
				sqlStr := fmt.Sprintf(`select max(descendant) from 
							(select descendant from metadata_closure
							where ancestor = '%s' and depth = %d and descendant > '%s'
							order by descendant asc limit %d) as tmp`, objectId, depth, lastDescendant, pageSize)
				row = m.db.QueryRowContext(ctx, sqlStr)
				var v sql.NullString
				if err = row.Scan(&v); err != nil {
					return
				}
				if v.Valid {
					startDescendant = lastDescendant
					lastDescendant = v.String
				} else {
					break
				}

				// 和祖先的所有关系对
				m.log.Info("MoveLargeSubTree delete closure start")
				args1 := make([]string, len(srcParentPaths))
				for i := range srcParentPaths {
					args1[i] = fmt.Sprintf("'%s'", srcParentPaths[i])
				}

				// 祖先相同的，一次性批量删除
				sqlStr = fmt.Sprintf(`delete from metadata_closure where 
					ancestor in (%s) and 
					descendant in (select descendant from metadata_closure
									where ancestor = '%s' and depth = %d and descendant > '%s' and descendant <='%s')`,
					strings.Join(args1, ","), objectId, depth, startDescendant, lastDescendant)
				if result, err = m.db.ExecContext(ctx, sqlStr); err != nil {
					return
				}
				if count, err = result.RowsAffected(); err != nil {
					return
				}
				deleteCount += count

				m.log.Info("MoveLargeSubTree delete closure end", zap.Int64("deleteCount", deleteCount))
			}
		}

		if deleteCount == 0 {
			m.log.Warn("empty src subtree", zap.String("objectId", objectId), zap.String("parentId", parentId))
		} else {
			m.log.Info(fmt.Sprintf("total delete closure: %d", deleteCount),
				zap.Int("srcParentPathsCount", len(srcParentPaths)),
				zap.Int("maxDescendantDepth", maxDepth))
		}
	}

	// 第二步：将上一步分离出的子树用笛卡尔积嫁接到新位置
	//
	// 考虑到超大目录，先查询孤立的树的后代，从根开始分批插入（深度由小到大）
	//
	// 1) 查询目标父路径
	// 2）查询孤立的树的后代，深度由小到大
	// 3）取笛卡尔积，插入到新目标位置
	sqlStr = "select ancestor, descendant, depth from metadata_closure where descendant=?"
	rows, err = m.db.QueryContext(ctx, sqlStr, parentId)
	if err != nil {
		return
	}

	dstParentPaths := make([]interfaces.MetadataClosure, 0)
	for rows.Next() {
		item := interfaces.MetadataClosure{}
		if err = rows.Scan(&item.Ancestor, &item.Descendant, &item.Depth); err != nil {
			rows.Close()
			return
		}
		dstParentPaths = append(dstParentPaths, item)
	}
	rows.Close()
	m.log.Info("MoveLargeSubTree before insert", zap.Any("dstParentPaths.depth", len(dstParentPaths)))

	// 分批插入
	for depth := 0; depth <= maxDepth; depth++ {
		lastDescendant := ""  // 包含
		startDescendant := "" // 不包含
		for {
			// 查孤立树的所有后代，深度由小到大
			// 先查插入的分页标记
			sqlStr := `select max(descendant) from 
					(select descendant from metadata_closure where 
				   	 ancestor = ? and depth = ? and descendant > ? limit ?) as tmp`
			row = m.db.QueryRowContext(ctx, sqlStr, objectId, depth, lastDescendant, pageSize)
			var v sql.NullString
			if err = row.Scan(&v); err != nil {
				return
			}

			if v.Valid {
				startDescendant = lastDescendant
				lastDescendant = v.String
			} else {
				break
			}

			// 一次性批量插入
			sqlStr = fmt.Sprintf(`insert into metadata_closure (ancestor, descendant,  depth)
								  select T1.ancestor, T2.descendant, T1.depth + T2.depth + 1
								  from metadata_closure as T1 
								  cross join
								  metadata_closure as T2
								  where 
								  -- T1 目标路径的所有祖先
								  T1.descendant='%s' and
								  -- T2 源路径指定深度的后代，分页
								  T2.ancestor='%s' and T2.depth = %d and T2.descendant > '%s' and T2.descendant <= '%s';`,
				parentId, objectId, depth, startDescendant, lastDescendant)

			result, err = m.db.ExecContext(ctx, sqlStr)
			if err != nil {
				return
			}
			if count, err = result.RowsAffected(); err != nil {
				return
			}
			insertCount += count
			m.log.Debug("MoveLargeSubTree insert closure", zap.Int64("count", count), zap.Int("depth", depth))
		}
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
