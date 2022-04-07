package archived

import (
	"time"

	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

type Pipeline struct {
	Name          string `gorm:"index"`
	BlueprintId   uint64
	Tasks         datatypes.JSON
	TotalTasks    int
	FinishedTasks int
	BeganAt       *time.Time
	FinishedAt    *time.Time `gorm:"index"`
	Status        string
	Message       string
	SpentSeconds  int
	common.Model
}
