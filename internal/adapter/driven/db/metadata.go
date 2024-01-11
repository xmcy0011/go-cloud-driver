package db

import (
	"context"
	"database/sql"

	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

type dbMetadata struct {
	db *sql.DB
}

func NewMetdata(db *sql.DB) interfaces.DBMetadata {
	return &dbMetadata{db: db}
}

func (d *dbMetadata) Add(ctx context.Context, meta interfaces.Metadata, tx *sql.Tx) error {
	sql := "insert into metadata(`object_id`,`parent_id`,`name`,`basic_attr`) values(?,?,?,?)"
	_, err := tx.Exec(sql, meta.ObjectId, meta.ParentId, meta.Name, meta.BasicAttr)
	return err
}

func (d *dbMetadata) QueryCountById(ctx context.Context, objectId string) (count int, err error) {
	sql := "select count(1) from metadata where objectId=?"
	row := d.db.QueryRow(sql, objectId)
	err = row.Scan(&count)
	return
}