package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubIssueComment struct {
	GithubId        int `gorm:"primaryKey"`
	IssueId         int `gorm:"index;comment:References the Issue"`
	Body            string
	AuthorUsername  string
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time `gorm:"index"`
	common.NoPKModel
}
