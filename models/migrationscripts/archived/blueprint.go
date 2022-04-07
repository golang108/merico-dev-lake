package archived

import (
	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

type Blueprint struct {
	Name       string
	Tasks      datatypes.JSON
	Enable     bool
	CronConfig string
	common.Model
}
