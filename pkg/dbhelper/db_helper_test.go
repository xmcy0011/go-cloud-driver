package dbhelper

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
)

func TestExecInTranscation(t *testing.T) {
	Convey("DBMetadataClosure.Add", t, func() {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		Convey("rollback", func() {
			// 回滚
			mock.ExpectBegin()
			mock.ExpectRollback()

			err = ExecInTranscation(db, func(tx *sql.Tx) error {
				return errors.New("error")
			})
			require.Error(t, err)
		})

		Convey("commit", func() {
			// 提交
			mock.ExpectBegin()
			mock.ExpectCommit()
			err = ExecInTranscation(db, func(tx *sql.Tx) error {
				return nil
			})
			require.NoError(t, err)
		})
	})
}
