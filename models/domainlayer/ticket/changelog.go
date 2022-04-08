package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Changelog struct {
	domainlayer.DomainEntity

	// collected fields
	IssueId     string `gorm:"index"`
	AuthorId    string `gorm:"type:char(255)"`
	AuthorName  string `gorm:"type:char(255)"`
	FieldId     string `gorm:"type:char(255)"`
	FieldName   string `gorm:"type:char(255)"`
	From        string `gorm:"type:char(255)"`
	To          string `gorm:"type:char(255)"`
	CreatedDate time.Time
}
