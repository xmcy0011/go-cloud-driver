package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/oklog/ulid/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

func MustInitSqlmock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func TestMetadataAdd(t *testing.T) {
	Convey("DBMetadata.Add", t, func() {
		db, mock := MustInitSqlmock(t)

		uid := ulid.Make().String()
		metadata := NewMetdata(db)

		Convey("success", func() {
			mock.ExpectExec("insert into metadata").WithArgs(uid, "cc", "aa", 1).WillReturnResult(sqlmock.NewErrorResult(nil))

			tx, err := db.Begin()
			require.NoError(t, err)
			err = metadata.Add(context.Background(), interfaces.Metadata{
				ObjectId:  uid,
				ParentId:  "cc",
				Name:      "aa",
				BasicAttr: 1,
			}, tx)
			require.NoError(t, err)
		})

		Convey("error", func() {
			mock.ExpectExec("insert into metadata").WithArgs(uid, "cc", "aa", 1).WillReturnResult(sqlmock.NewErrorResult(nil))

			tx, err := db.Begin()
			require.NoError(t, err)
			err = metadata.Add(context.Background(), interfaces.Metadata{
				ObjectId:  uid,
				ParentId:  "cc",
				Name:      "aa",
				BasicAttr: 1,
			}, tx)
			require.NoError(t, err)
		})
	})
}

func TestQueryCountById(t *testing.T) {
	Convey("QueryCountById", t, func() {
		db, mock := MustInitSqlmock(t)
		uid := "id1"
		metadata := NewMetdata(db)

		Convey("success", func() {
			mock.ExpectCommit()
			// insert
			tx, err := db.Begin()
			require.NoError(t, err)
			err = metadata.Add(context.Background(), interfaces.Metadata{
				ObjectId:  uid,
				ParentId:  "cc",
				Name:      "aa",
				BasicAttr: 1,
			}, tx)
			require.NoError(t, err)

			// query
			mock.ExpectBegin()
			mock.ExpectExec("select").WithArgs(uid).WillReturnError(nil)
			mock.ExpectCommit()

			count, err := metadata.QueryCountById(context.Background(), uid)
			require.NoError(t, err)
			require.Equal(t, count, 1)
		})
	})
}
