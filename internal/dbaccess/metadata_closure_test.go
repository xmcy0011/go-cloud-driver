package dbaccess

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

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

func MustInitDb() *sql.DB {
	l, _ := logger.NewZapLogger("", true)
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

func Mock(db *sql.DB) error {
	_, err := db.Exec("delete from hydra_oauth2_jti_blacklist WHERE nid = '078dada1-d130-11ee-b0b9-8af344cb7127' AND expires_at < ?",
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	signature := ulid.Make().String()
	_, err = db.Exec("INSERT INTO `hydra_oauth2_jti_blacklist` (`expires_at`, `nid`, `signature`) VALUES (?, '078dada1-d130-11ee-b0b9-8af344cb7127', ?)",
		time.Now().Format("2006-01-02 15:04:05"), signature)
	if err != nil {
		return err
	}
	return nil
}

func BenchmarkHyrd(b *testing.B) {
	dbConn := MustInitDb()
	for i := 0; i < b.N; i++ {

		group := sync.WaitGroup{}
		for j := 0; j < 20; j++ {
			group.Add(1)
			go func() {
				err := Mock(dbConn)
				assert.NoError(b, err)
				group.Done()
			}()
		}
		group.Wait()

	}
}
