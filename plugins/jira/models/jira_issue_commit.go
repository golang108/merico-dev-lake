package models

import "github.com/merico-dev/lake/models/common"

type JiraIssueCommit struct {
	common.NoPKModel
	SourceId  uint64 `gorm:"primaryKey"`
	IssueId   uint64 `gorm:"primaryKey"`
	CommitSha string `gorm:"primaryKey;type:char(40)"`
	CommitUrl string `gorm:"type:varchar(255)"`
}

func (JiraIssueCommit) TableName() string {
	return "_tool_jira_issue_commits"
}
