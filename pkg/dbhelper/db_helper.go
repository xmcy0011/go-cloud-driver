package dbhelper

import (
	"database/sql"
)

type FuncCallback func(tx *sql.Tx) error

// ExecInTranscation：在事务中执行
// db: 数据库连接
// f: 回调函数，当返回 error 时，自动回滚事务，否则自动提交
func ExecInTranscation(db *sql.DB, f FuncCallback) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	err = f(tx)
	return err
}
