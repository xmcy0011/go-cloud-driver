package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/oklog/ulid/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

func TestMetadataAdd(t *testing.T) {
	Convey("DBMetadata.Add", t, func() {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		uid := ulid.Make().String()
		metadata := NewMetdata(db)

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("insert into metadata").WithArgs(uid, "cc", "aa", 1).WillReturnResult(sqlmock.NewErrorResult(nil))
			mock.ExpectCommit()

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
			mock.ExpectBegin()
			mock.ExpectExec("insert into metadata").WithArgs(uid, "cc", "aa", 1).WillReturnResult(sqlmock.NewErrorResult(nil))
			mock.ExpectCommit()

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
