package service

import (
	"context"
	"database/sql"

	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

type MetadataService interface {
	Create(ctx context.Context, metadata interfaces.Metadata) error
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

}
