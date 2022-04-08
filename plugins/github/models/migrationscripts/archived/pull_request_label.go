package archived

import (
	"github.com/merico-dev/lake/models/common"
)

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type GithubPullRequestLabel struct {
	PullId    int    `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey"`
	common.NoPKModel
}
