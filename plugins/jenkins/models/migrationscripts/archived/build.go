package archived

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

// JenkinsBuild db entity for jenkins build
type JenkinsBuild struct {
	common.NoPKModel
	JobName           string  `gorm:"primaryKey;type:varchar(255)"`
	Duration          float64 // build time
	DisplayName       string  // "#7"
	EstimatedDuration float64
	Number            int64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	Result            string
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
	CommitSha         string
}
