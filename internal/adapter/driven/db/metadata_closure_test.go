package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
)

func TestClosureAdd(t *testing.T) {
	Convey("DBMetadataClosure.Add", t, func() {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		closure := NewMetadataClosure(db)

		ancestor := "f"
		descendant := "d"

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("insert into").WithArgs(descendant, ancestor, descendant, descendant).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			tx, err := db.Begin()
			require.NoError(t, err)
			id, err := closure.Add(context.Background(), ancestor, descendant, tx)
			require.NoError(t, err)
			require.Equal(t, true, id > 0)
		})
	})
}
