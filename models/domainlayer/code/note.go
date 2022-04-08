package code

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Note struct {
	domainlayer.DomainEntity
	PrId        uint64 `gorm:"index;comment:References the pull request for this note"`
	Type        string `gorm:"type:char(100)"`
	Author      string `gorm:"type:char(255)"`
	Body        string
	Resolvable  bool `gorm:"comment:Is or is not a review comment"`
	System      bool `gorm:"comment:Is or is not auto-generated vs. human generated"`
	CreatedDate time.Time
}
