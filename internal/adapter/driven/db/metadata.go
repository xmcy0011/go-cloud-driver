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
	sql := "insert into metadata(`object_id`,`parent_id`,`name`,`object_type`) values(?,?,?,?)"
	_, err := tx.ExecContext(ctx, sql, meta.ObjectId, meta.ParentId, meta.Name, meta.ObjectId)
	return err
}

func (d *dbMetadata) QueryCountById(ctx context.Context, objectId string) (count int, err error) {
	sql := "select count(1) from metadata where objectId=?"
	row := d.db.QueryRow(sql, objectId)
	err = row.Scan(&count)
	return
}

func (d *dbMetadata) QueryById(ctx context.Context, objectId string) (*interfaces.Metadata, error) {
	sql := "select object_id,parent_id,name,object_type from metadata where object_id=? limit 1"
	row := d.db.QueryRow(sql, objectId)
	metdata := interfaces.Metadata{}
	if err := row.Scan(&metdata.ObjectId, &metdata.ParentId, &metdata.Name, &metdata.ObjectType); err != nil {
		return nil, err
	}
	return &metdata, nil
}

func (d *dbMetadata) UpdateParentId(ctx context.Context, objectId, newParentId string, tx *sql.Tx) (rowsAffected int64, err error) {
	sql := "update metadata set parent_id=? where object_id=?"
	r, err := tx.ExecContext(ctx, sql, newParentId, objectId)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
