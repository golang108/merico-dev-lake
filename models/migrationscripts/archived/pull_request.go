package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type PullRequest struct {
	domainlayer.DomainEntity
	RepoId         string `gorm:"index"`
	Status         string `gorm:"comment:open/closed or other"`
	Title          string
	Description    string
	Url            string `gorm:"type:char(255)"`
	AuthorName     string `gorm:"type:char(100)"`
	AuthorId       int
	ParentPrId     string `gorm:"index;type:varchar(100)"`
	Key            int
	CreatedDate    time.Time
	MergedDate     *time.Time
	ClosedAt       *time.Time
	Type           string
	Component      string
	MergeCommitSha string `gorm:"type:char(40)"`
	HeadRef        string
	BaseRef        string
	BaseCommitSha  string
	HeadCommitSha  string
}

type PullRequestCommit struct {
	CommitSha     string `gorm:"primaryKey"`
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	common.NoPKModel
}

type PullRequestIssue struct {
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	IssueId       string `gorm:"primaryKey;type:varchar(255)"`
	PullNumber    int
	IssueNumber   int
	common.NoPKModel
}

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type PullRequestLabel struct {
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	LabelName     string `gorm:"primaryKey"`
	common.NoPKModel
}
