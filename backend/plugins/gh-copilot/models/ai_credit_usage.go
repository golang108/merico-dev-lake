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

package models

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

// AI credit usage "level" identifies which billing endpoint a usage row came from.
const (
	AiCreditUsageLevelEnterprise   = "enterprise"
	AiCreditUsageLevelOrganization = "organization"
)

// GhCopilotAiCreditUsage stores a single AI credit billing usage line item for a
// time period, sourced from the GitHub Billing Usage "AI credit usage" reports
// (enterprise and organization endpoints). Each row corresponds to one usageItem
// (a product/sku/model/unitType combination) within a single day.
type GhCopilotAiCreditUsage struct {
	ConnectionId uint64 `gorm:"primaryKey" json:"connectionId"`
	ScopeId      string `gorm:"primaryKey;type:varchar(255)" json:"scopeId"`
	// Level distinguishes enterprise vs organization sourced usage.
	Level string `gorm:"primaryKey;type:varchar(20)" json:"level"`
	// Time period of the usage line (daily granularity).
	Year  int `gorm:"primaryKey" json:"year"`
	Month int `gorm:"primaryKey" json:"month"`
	Day   int `gorm:"primaryKey" json:"day"`
	// Usage line dimensions.
	Product      string `gorm:"primaryKey;type:varchar(100)" json:"product"`
	Sku          string `gorm:"primaryKey;type:varchar(150)" json:"sku"`
	Model        string `gorm:"primaryKey;type:varchar(150)" json:"model"`
	UnitType     string `gorm:"primaryKey;type:varchar(50)" json:"unitType"`
	CostCenterId string `gorm:"primaryKey;type:varchar(100)" json:"costCenterId"`

	// Date is the time period expressed as a date (year-month-day) for easier
	// time-series querying in dashboards.
	Date time.Time `gorm:"type:date;index" json:"date"`

	// Top-level context (may be empty for the unfiltered aggregate report).
	Enterprise     string `gorm:"type:varchar(100)" json:"enterprise"`
	Organization   string `gorm:"type:varchar(100)" json:"organization"`
	UserLogin      string `gorm:"type:varchar(255)" json:"userLogin"`
	CostCenterName string `gorm:"type:varchar(255)" json:"costCenterName"`

	// Billing metrics. GitHub returns fractional values, so these are float64.
	PricePerUnit     float64 `json:"pricePerUnit" gorm:"comment:Price per unit for the line item"`
	GrossQuantity    float64 `json:"grossQuantity" gorm:"comment:Gross quantity consumed (e.g. credits)"`
	GrossAmount      float64 `json:"grossAmount" gorm:"comment:Gross billed amount"`
	DiscountQuantity float64 `json:"discountQuantity" gorm:"comment:Discounted quantity"`
	DiscountAmount   float64 `json:"discountAmount" gorm:"comment:Discounted amount"`
	NetQuantity      float64 `json:"netQuantity" gorm:"comment:Net quantity after discounts"`
	NetAmount        float64 `json:"netAmount" gorm:"comment:Net billed amount after discounts"`

	common.NoPKModel
}

func (GhCopilotAiCreditUsage) TableName() string {
	return "_tool_copilot_ai_credit_usage"
}
