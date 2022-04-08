package archived

import (
	"github.com/merico-dev/lake/models/common"
)

type JiraBoard struct {
	common.NoPKModel
	SourceId  uint64 `gorm:"primaryKey"`
	BoardId   uint64 `gorm:"primaryKey"`
	ProjectId uint
	Name      string
	Self      string
	Type      string
}

func (JiraBoard) TableName() string {
	return "_tool_jira_boards"
}

type JiraBoardIssue struct {
	SourceId uint64 `gorm:"primaryKey"`
	BoardId  uint64 `gorm:"primaryKey"`
	IssueId  uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

func (JiraBoardIssue) TableName() string {
	return "_tool_jira_board_issues"
}
