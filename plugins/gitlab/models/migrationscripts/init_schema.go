package migrationscripts

import (
	"context"

	"github.com/merico-dev/lake/plugins/gitlab/models/migrationscripts/archived"
	"gorm.io/gorm"
)

const (
	Owner = "Gitlab"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.GitlabProject{},
		&archived.GitlabMergeRequest{},
		&archived.GitlabCommit{},
		&archived.GitlabTag{},
		&archived.GitlabProjectCommit{},
		&archived.GitlabPipeline{},
		&archived.GitlabReviewer{},
		&archived.GitlabMergeRequestNote{},
		&archived.GitlabMergeRequestCommit{},
		&archived.GitlabUser{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220407201136
}

func (*InitSchemas) Owner() string {
	return Owner
}

func (*InitSchemas) Comment() string {
	return "create init schemas"
}
