package models

import (
	"github.com/merico-dev/lake/models/common"
)

type JiraBoard struct {
	common.NoPKModel
	SourceId  uint64 `gorm:"primaryKey"`
	BoardId   uint64 `gorm:"primaryKey"`
	ProjectId uint
	Name      string `gorm:"type:varchar(255)"`
	Self      string `gorm:"type:varchar(255)"`
	Type      string `gorm:"type:varchar(100)"`
}

func (JiraBoard) TableName() string {
	return "_tool_jira_boards"
}
