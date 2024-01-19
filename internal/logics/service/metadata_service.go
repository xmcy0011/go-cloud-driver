package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/common"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
	"github.com/xmcy0011/go-cloud-driver/pkg/dbhelper"
	"github.com/xmcy0011/go-cloud-driver/pkg/errhelper"
	"go.uber.org/zap"
)

type MetadataService interface {
	Create(ctx context.Context, metadata interfaces.Metadata) error
	QueryById(ctx context.Context, objectId string) (*interfaces.Metadata, error)
	MoveDir(ctx context.Context, objectId, newParentId string) error
	QuerySubTree(ctx context.Context, objectId string) ([]interfaces.MetadataNode, error)
}

type metadataService struct {
	metadata interfaces.DBMetadata
	closure  interfaces.DBMetadataClosure
	db       *sql.DB
	log      *zap.Logger
}

func NewMetadataService(db *sql.DB, metadata interfaces.DBMetadata, closure interfaces.DBMetadataClosure) MetadataService {
	return &metadataService{db: db, metadata: metadata, closure: closure, log: interfaces.MustNewLogger()}
}

func (m *metadataService) Create(ctx context.Context, metadata interfaces.Metadata) error {
	err := dbhelper.ExecInTranscation(m.db, func(tx *sql.Tx) error {
		err := m.metadata.Add(ctx, metadata, tx)
		if err != nil {
			return err
		}
		_, err = m.closure.Add(ctx, metadata.ParentId, metadata.ObjectId, tx)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (m *metadataService) QueryById(ctx context.Context, objectId string) (*interfaces.Metadata, error) {
	return m.metadata.QueryById(ctx, objectId)
}

func (m *metadataService) MoveDir(ctx context.Context, objectId, newParentId string) error {
	_, err := ulid.Parse(objectId)
	if err != nil {
		return common.WithCause(common.RestBadRequest, "invalid objectId")
	}

	_, err = ulid.Parse(newParentId)
	if err != nil {
		return common.WithCause(common.RestBadRequest, "invalid newParentId")
	}

	if objectId == newParentId {
		return common.WithCause(common.RestBadRequest, "invalid newParentId")
	}

	_, err = m.metadata.QueryById(ctx, objectId)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.WithCause(common.RestBadRequest, "objectId not found")
		}
		return err
	}

	info, err := m.metadata.QueryById(ctx, newParentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.WithCause(common.RestBadRequest, "parentId not found")
		}
		return err
	}
	if info.ObjectType != int(common.ObjectTypeDir) {
		return common.WithCause(common.RestBadRequest, "parentId msut be dir")
	}

	// 移动时，需要先删除原来的联系，再建立新的联系
	err = dbhelper.ExecInTranscation(m.db, func(tx *sql.Tx) error {
		// 不允许移动到其子目录
		isChild, err := m.closure.CheckIsDescendant(ctx, objectId, newParentId, tx)
		if err != nil {
			return errhelper.WithFileLine(err)
		}
		if isChild {
			return common.WithCause(common.RestBadRequest, "can not move to sub dir")
		}

		_, err = m.metadata.UpdateParentId(ctx, objectId, newParentId, tx)
		if err != nil {
			return errhelper.WithFileLine(err)
		}

		deleteCount, insertCount, err := m.closure.MoveSubTree(ctx, objectId, newParentId, tx)
		if err != nil {
			return errhelper.WithFileLine(err)
		}

		m.log.Info(fmt.Sprintf("move sucess, deleteCount: %d, insertCount: %d", deleteCount, insertCount))

		return nil
	})

	return err
}

func (m *metadataService) QuerySubTree(ctx context.Context, objectId string) ([]interfaces.MetadataNode, error) {
	_, err := m.metadata.QueryById(ctx, objectId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.WithCause(common.RestBadRequest, "objectId not found")
		}
		return nil, err
	}

	return m.closure.QueryAllDescendants(ctx, objectId)
}
