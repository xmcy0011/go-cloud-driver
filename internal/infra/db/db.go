package dbaccess

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xmcy0011/go-cloud-driver/internal/conf"
)

var Db *sql.DB = nil

func MustInit(c conf.Database) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&timeout=%ds&readTimeout=%ds&writeTimeout=%ds", c.UserName, c.Password, c.Host, c.Port, c.Db,
		int(c.Timeout.Seconds()), int(c.ReadTimeout.Seconds()), int(c.WriteTimeout.Seconds()))
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	Db = conn
}
