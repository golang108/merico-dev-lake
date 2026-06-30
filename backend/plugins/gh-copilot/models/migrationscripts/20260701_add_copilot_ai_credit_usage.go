/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migrationscripts

import (
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addCopilotAiCreditUsage struct{}

type aiCreditUsage20260701 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(255)"`
	Level        string `gorm:"primaryKey;type:varchar(20)"`
	Year         int    `gorm:"primaryKey"`
	Month        int    `gorm:"primaryKey"`
	Day          int    `gorm:"primaryKey"`
	Product      string `gorm:"primaryKey;type:varchar(100)"`
	Sku          string `gorm:"primaryKey;type:varchar(150)"`
	Model        string `gorm:"primaryKey;type:varchar(150)"`
	UnitType     string `gorm:"primaryKey;type:varchar(50)"`
	CostCenterId string `gorm:"primaryKey;type:varchar(100)"`

	Date time.Time `gorm:"type:date;index"`

	Enterprise     string `gorm:"type:varchar(100)"`
	Organization   string `gorm:"type:varchar(100)"`
	UserLogin      string `gorm:"type:varchar(255)"`
	CostCenterName string `gorm:"type:varchar(255)"`

	PricePerUnit     float64
	GrossQuantity    float64
	GrossAmount      float64
	DiscountQuantity float64
	DiscountAmount   float64
	NetQuantity      float64
	NetAmount        float64

	archived.NoPKModel
}

func (aiCreditUsage20260701) TableName() string {
	return "_tool_copilot_ai_credit_usage"
}

func (script *addCopilotAiCreditUsage) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes,
		&aiCreditUsage20260701{},
	)
}

func (*addCopilotAiCreditUsage) Version() uint64 {
	return 20260701000000
}

func (*addCopilotAiCreditUsage) Name() string {
	return "Add Copilot AI credit usage table"
}
