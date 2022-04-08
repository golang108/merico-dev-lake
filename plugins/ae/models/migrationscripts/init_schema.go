package migrationscripts

import (
	"context"

	"github.com/merico-dev/lake/plugins/ae/models/migrationscripts/archived"
	"gorm.io/gorm"
)

const (
	Owner = "AE"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.AECommit{},
		&archived.AEProject{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220407201133
}

func (*InitSchemas) Owner() string {
	return Owner
}

func (*InitSchemas) Comment() string {
	return "create init schemas"
}
