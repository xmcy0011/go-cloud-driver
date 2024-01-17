package service

import (
	"context"
	"database/sql"

	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

type MetadataService interface {
	Create(ctx context.Context, metadata interfaces.Metadata) error
	QueryById(ctx context.Context, ObjectId string) (*interfaces.Metadata, error)
}

type metadataService struct {
	metadata interfaces.DBMetadata
	closure  interfaces.DBMetadataClosure
	db       *sql.DB
}

func NewMetadataService(db *sql.DB, metadata interfaces.DBMetadata, closure interfaces.DBMetadataClosure) MetadataService {
	return &metadataService{db: db, metadata: metadata, closure: closure}
}

func (m *metadataService) Create(ctx context.Context, metadata interfaces.Metadata) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	err = m.metadata.Add(ctx, metadata, tx)
	if err != nil {
		return err
	}
	_, err = m.closure.Add(ctx, metadata.ObjectId, metadata.ParentId, tx)
	if err != nil {
		return err
	}
	return nil
}

func (m *metadataService) QueryById(ctx context.Context, objectId string) (*interfaces.Metadata, error) {
	return m.metadata.QueryById(ctx, objectId)
}
