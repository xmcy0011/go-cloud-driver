package logics

import (
	"database/sql"

	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

func MustInit(db *sql.DB, metadata interfaces.DBMetadata, closure interfaces.DBMetadataClosure) {

}
