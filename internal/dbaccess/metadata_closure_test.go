package dbaccess

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/oklog/ulid/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmcy0011/go-cloud-driver/internal/conf"
	"github.com/xmcy0011/go-cloud-driver/internal/infra/mysqldb"
	"github.com/xmcy0011/go-cloud-driver/pkg/logger"
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

func TestMoveLargeSubTree(t *testing.T) {
	Convey("MoveLargeSubTree", t, func() {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		closure := NewMetadataClosure(db)
		ctx := context.Background()

		objectId := "c"
		parentId := "b"

		Convey("empty subtree", func() {
			mock.ExpectQuery("^select").WithArgs(objectId, objectId, 10).WillReturnError(sql.ErrNoRows)
			_, _, err := closure.MoveLargeSubTree(ctx, objectId, parentId, 10)
			require.Equal(t, "empty subtree", err.Error())
		})

		Convey("success", func() {
			// test
			// |-a
			//   |- c
			//   	|- 2.txt
			// |-b
			// 按照路径枚举（深度降序）总结其闭包如下：
			// test/a/c/2.txt => (2.txt, 2.txt, 0): 9, (2.txt, c, 1): 10, (2.txt, a, 2): 11, (2.txt, test, 3): 12
			// test/a/c       => (c, c, 0): 6, (c, a, 1): 7, (c, test, 2): 8
			// test/a         => (a, a, 0): 4, (a, test, 1): 5
			// test/b         => (b, b, 0): 2, (b, test, 1): 3
			// test           => (test, test, 0): 1
			// 移动 test/a/c 到 test/b 下，相当于删除 test/a 的前缀，删除顺序从叶子开始：
			// 2.txt 深度最深为3，删除 (2.txt, a, 2): 10, (2.txt, test, 3): 11
			// c 深度次之为2，删除 (c, a, 1): 7, (c, test, 2): 8
			mock.ExpectQuery("^select id from metadata_closure where descendant").WithArgs(objectId, objectId, 10).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(11).AddRow(10).AddRow(8).AddRow(7))
			mock.ExpectExec("^delete from metadata_closure where id").WithArgs(11, 10, 8, 7)
			_, _, err := closure.MoveLargeSubTree(ctx, objectId, parentId, 10)
			require.NoError(t, err)
		})
	})
}

func MustInitDb() *sql.DB {
	l, _ := logger.NewZapLogger(false)
	port, _ := strconv.Atoi(os.Getenv("mysql_port"))
	dbConn := mysqldb.MustInit(
		conf.Database{
			UserName: os.Getenv("mysql_user"),
			Password: os.Getenv("mysql_pwd"),
			Host:     os.Getenv("mysql_host"),
			Port:     port,
			Db:       os.Getenv("mysql_db"),
		}, l)
	return dbConn
}

// BenchmarkMoveSubTree 测试闭包表目录移动性能
// example: 测试100万文件的移动
// - go test -benchmem -run=^$ -bench ^BenchmarkMoveSubTree$ github.com/xmcy0011/go-cloud-driver/internal/adapter/driven/db -benchtime=1000000x
// result:
// 1万：1s
// 10万：2s
// 100万：161s
func BenchmarkMoveSubTree(b *testing.B) {
	b.Run("prepare", func(b *testing.B) {
		root := "00000000000000000000000000"
		dbConn := MustInitDb()
		closure := NewMetadataClosure(dbConn)

		b.StopTimer()

		// 确保所有的祖先已存在
		num, err := closure.QueryCountByAncestor(context.Background(), root)
		assert.NoError(b, err)
		if num <= 0 {
			tx, err := dbConn.Begin()
			assert.NoError(b, err)
			closure.Add(context.Background(), root, root, tx)
			assert.NoError(b, tx.Commit())
		}

		dirObjectId := ulid.Make().String()

		// 单目录下模拟创建 b.N 个文件的闭包关系
		tx, err := dbConn.Begin()
		closure.Add(context.Background(), root, dirObjectId, tx)
		assert.NoError(b, err)
		for i := 0; i < b.N; i++ {
			closure.Add(context.Background(), dirObjectId, ulid.Make().String(), tx)
			if (i/1000 == 0 && i > 0) || (i+1 == b.N) {
				assert.NoError(b, tx.Commit())
				tx, err = dbConn.Begin()
				assert.NoError(b, err)
			}
		}

		subCount, err := closure.QueryCountByAncestor(context.Background(), dirObjectId)
		assert.NoError(b, err)

		// 创建目标位置
		target := ulid.Make().String()
		tx, err = dbConn.Begin()
		assert.NoError(b, err)
		closure.Add(context.Background(), target, target, tx)
		assert.NoError(b, tx.Commit())
		b.StartTimer()

		// 移动子树到目标位置
		// t1 := time.Now()
		// tx, err = dbConn.Begin()
		// assert.NoError(b, err)
		// deleteCount, insertCount, err := closure.MoveSubTree(context.Background(), root, target, tx)
		// assert.NoError(b, err)
		// assert.NoError(b, tx.Commit())

		b.Logf("b.N: %d, prepare data: %d, objectId: %s, moveSubTree, newObjectId: %s",
			b.N, subCount, dirObjectId, target)
		// b.Logf("b.N: %d, cost: %2.f s, prepare data: %d, objectId: %s, moveSubTree, deleteCount: %d, insertCount: %d, newObjectId: %s",
		// 	b.N, time.Since(t1).Seconds(), subCount, root, deleteCount, insertCount, target)
	})
}
