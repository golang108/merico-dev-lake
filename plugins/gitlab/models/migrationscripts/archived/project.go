package archived

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type GitlabProject struct {
	GitlabId                int    `gorm:"primaryKey"`
	Name                    string `gorm:"type:varchar(255)"`
	Description             string
	DefaultBranch           string `gorm:"varchar(255)"`
	PathWithNamespace       string
	WebUrl                  string
	CreatorId               int
	Visibility              string
	OpenIssuesCount         int
	StarCount               int
	ForkedFromProjectId     int
	ForkedFromProjectWebUrl string

	CreatedDate time.Time
	UpdatedDate *time.Time
	common.NoPKModel
}
