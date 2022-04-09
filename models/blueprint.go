package models

import (
	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

type Blueprint struct {
	Name       string         `json:"name" validate:"required"`
	Tasks      datatypes.JSON `json:"tasks" validate:"required"`
	Enable     bool           `json:"enable"`
	CronConfig string         `json:"cronConfig" validate:"required"`
	common.Model
}
