package migration

import (
	"time"
)

const (
	tableName = "migration_history"
)

type MigrationHistory struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	ScriptOwner   string
	ScriptVersion uint64
	ScriptComment string
}

func (MigrationHistory) TableName() string {
	return tableName
}
