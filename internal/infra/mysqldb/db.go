package mysqldb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xmcy0011/go-cloud-driver/internal/conf"
	"go.uber.org/zap"
)

func MustInit(c conf.Database, log *zap.Logger) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&timeout=%ds&readTimeout=%ds&writeTimeout=%ds", c.UserName, c.Password, c.Host, c.Port, c.Db,
		int(c.Timeout.Seconds()), int(c.ReadTimeout.Seconds()), int(c.WriteTimeout.Seconds()))
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err := conn.Ping(); err != nil {
		panic(err)
	}
	log.Info("success connect to mysql", zap.String("addr", fmt.Sprintf("%s:%d", c.Host, c.Port)))
	return conn
}
