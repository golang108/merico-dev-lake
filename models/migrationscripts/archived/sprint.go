package archived

import (
	"time"

	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/models/domainlayer"
)

type Sprint struct {
	domainlayer.DomainEntity
	Name            string
	Url             string
	Status          string
	Title           string
	StartedDate     *time.Time
	EndedDate       *time.Time
	CompletedDate   *time.Time
	OriginalBoardID string
}

type SprintIssue struct {
	common.NoPKModel
	SprintId      string `gorm:"primaryKey"`
	IssueId       string `gorm:"primaryKey"`
	IsRemoved     bool
	AddedDate     *time.Time
	RemovedDate   *time.Time
	AddedStage    *string `gorm:"type:varchar(255)"`
	ResolvedStage *string `gorm:"type:varchar(255)"`
}
