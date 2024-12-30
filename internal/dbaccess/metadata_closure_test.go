package dbaccess

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/oklog/ulid/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmcy0011/go-cloud-driver/internal/common"
	"github.com/xmcy0011/go-cloud-driver/internal/conf"
	"github.com/xmcy0011/go-cloud-driver/internal/infra/mysqldb"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
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
// - go test -benchmem -run=^$ -bench ^BenchmarkMoveSubTree$ github.com/xmcy0011/go-cloud-driver/internal/dbaccess -benchtime=1000000x -cpu=1
// result:
// 1万：1s      分批移动 1s
// 10万：6s     分批移动 4s
// 100万：161s  分批移动
func BenchmarkMoveSubTree(b *testing.B) {
	b.Run("prepare", func(b *testing.B) {
		dbConn := MustInitDb()
		logger, _ := logger.NewZapLogger(true)
		common.SetLogger(logger)
		closure := NewMetadataClosure(dbConn)
		dbMetadata := NewMetdata(dbConn)

		b.StopTimer()

		//b.N = 10000

		// 单目录下模拟创建 b.N 个文件的闭包关系
		// /dir1/dir2/file....
		tx, err := dbConn.Begin()
		dir1 := ulid.Make().String()
		dbMetadata.Add(context.Background(), interfaces.Metadata{ObjectId: dir1, ParentId: dir1, Name: dir1, ObjectType: 1}, tx)
		closure.Add(context.Background(), dir1, dir1, tx)
		srcDirObjectId := ulid.Make().String()
		dbMetadata.Add(context.Background(), interfaces.Metadata{ObjectId: srcDirObjectId, ParentId: dir1, Name: srcDirObjectId, ObjectType: 1}, tx)
		closure.Add(context.Background(), dir1, srcDirObjectId, tx)

		logger.Info(fmt.Sprintf("start prepare data, parentId: %s", srcDirObjectId))

		assert.NoError(b, err)
		for i := 0; i < b.N; i++ {
			id := ulid.Make().String()
			dbMetadata.Add(context.Background(), interfaces.Metadata{ObjectId: id, ParentId: srcDirObjectId, Name: id, ObjectType: 100 + i}, tx)
			closure.Add(context.Background(), srcDirObjectId, id, tx)
			if (i%10000 == 0 && i > 0) || (i+1 == b.N) {
				assert.NoError(b, tx.Commit())
				tx, err = dbConn.Begin()
				assert.NoError(b, err)
				logger.Info(fmt.Sprintf("b.N: %d, prepare data: %d", b.N, i))
			}
		}

		subCount, err := closure.QueryCountByAncestor(context.Background(), srcDirObjectId)
		assert.NoError(b, err)
		logger.Info(fmt.Sprintf("b.N: %d, success create mock dir, count: %d", b.N, subCount))

		// 创建目标位置：/test/a/b
		tx, err = dbConn.Begin()
		assert.NoError(b, err)
		targetRoot := ulid.Make().String()
		closure.Add(context.Background(), targetRoot, targetRoot, tx)
		targetDir := ulid.Make().String()
		closure.Add(context.Background(), targetRoot, targetDir, tx)
		assert.NoError(b, tx.Commit())
		b.StartTimer()

		// 移动子树到目标位置（测试一次性移动）
		t1 := time.Now()
		tx, err = dbConn.Begin()
		assert.NoError(b, err)
		deleteCount, insertCount, err := closure.MoveSubTree(context.Background(), srcDirObjectId, targetDir, tx)
		assert.NoError(b, err)
		tx.Commit()

		b.Logf("b.N: %d, prepare data: %d, srcDirObjectId: %s, targetDirObjectId: %s",
			b.N, subCount, srcDirObjectId, targetDir)
		b.Logf("b.N: %d, moveSubTree, cost: %2.f s, deleteCount: %d, insertCount: %d",
			b.N, time.Since(t1).Seconds(), deleteCount, insertCount)
	})
}

// BenchmarkMoveLargeSubTree 测试闭包表目录移动性能
// example: 测试100万文件的移动
// - go test -benchmem -run=^$ -bench ^BenchmarkMoveSubTree$ github.com/xmcy0011/go-cloud-driver/internal/dbaccess -benchtime=1000000x -cpu=1
func BenchmarkMoveLargeSubTree(b *testing.B) {
	b.Run("prepare", func(b *testing.B) {
		dbConn := MustInitDb()
		logger, _ := logger.NewZapLogger(true)
		common.SetLogger(logger)
		closure := NewMetadataClosure(dbConn)

		b.StopTimer()

		//b.N = 10000

		// 单目录下模拟创建 b.N 个文件的闭包关系
		// /dir1/dir2/file....
		tx, err := dbConn.Begin()
		dir1 := ulid.Make().String()
		closure.Add(context.Background(), dir1, dir1, tx)
		srcDirObjectId := ulid.Make().String()
		closure.Add(context.Background(), dir1, srcDirObjectId, tx)
		assert.NoError(b, err)
		for i := 0; i < b.N; i++ {
			closure.Add(context.Background(), srcDirObjectId, ulid.Make().String(), tx)
			if (i%10000 == 0 && i > 0) || (i+1 == b.N) {
				assert.NoError(b, tx.Commit())
				tx, err = dbConn.Begin()
				assert.NoError(b, err)
				logger.Info(fmt.Sprintf("b.N: %d, prepare data: %d", b.N, i))
			}
		}

		subCount, err := closure.QueryCountByAncestor(context.Background(), srcDirObjectId)
		assert.NoError(b, err)
		logger.Info(fmt.Sprintf("b.N: %d, success create mock dir, count: %d", b.N, subCount))

		// 创建目标位置：/test/a/b
		tx, err = dbConn.Begin()
		assert.NoError(b, err)
		targetRoot := ulid.Make().String()
		closure.Add(context.Background(), targetRoot, targetRoot, tx)
		targetDir := ulid.Make().String()
		closure.Add(context.Background(), targetRoot, targetDir, tx)
		assert.NoError(b, tx.Commit())
		b.StartTimer()

		// /dir1/dir2 => /test/a/b
		// 移动子树到目标位置（测试分批移动）
		t1 := time.Now()
		deleteCount, insertCount, err := closure.MoveLargeSubTree(context.Background(), srcDirObjectId, targetDir, 10000)
		assert.NoError(b, err)

		b.Logf("b.N: %d, prepare data: %d, srcDirObjectId: %s, targetDirObjectId: %s",
			b.N, subCount, srcDirObjectId, targetDir)
		b.Logf("b.N: %d, moveSubTree, cost: %2.f s, deleteCount: %d, insertCount: %d",
			b.N, time.Since(t1).Seconds(), deleteCount, insertCount)
	})
}
