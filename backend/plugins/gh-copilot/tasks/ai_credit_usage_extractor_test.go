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

package tasks

import (
	"testing"
	"time"

	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitDay(t *testing.T) {
	y, m, d, err := splitDay("2026-06-08")
	require.Nil(t, err)
	assert.Equal(t, 2026, y)
	assert.Equal(t, 6, m)
	assert.Equal(t, 8, d)

	_, _, _, err = splitDay("not-a-date")
	assert.NotNil(t, err)
}

func TestMapAiCreditUsageReport(t *testing.T) {
	raw := []byte(`{
		"timePeriod": {"year": 2026, "month": 6, "day": 8},
		"enterprise": "Avocado Corp",
		"usageItems": [
			{
				"product": "Copilot",
				"sku": "Copilot AI Credits",
				"model": "Claude Opus 4.8",
				"unitType": "ai-credits",
				"pricePerUnit": 0.01,
				"grossQuantity": 5445.77595,
				"grossAmount": 54.4577595,
				"discountQuantity": 5445.77595,
				"discountAmount": 54.4577595,
				"netQuantity": 0.0,
				"netAmount": 0.0
			},
			{
				"product": "Copilot",
				"sku": "Copilot AI Credits",
				"model": "Claude Sonnet 4.6",
				"unitType": "ai-credits",
				"pricePerUnit": 0.01,
				"grossQuantity": 176.97417,
				"grossAmount": 1.7697417,
				"discountQuantity": 175.3533,
				"discountAmount": 1.753533,
				"netQuantity": 1.62087,
				"netAmount": 0.0162087
			}
		]
	}`)

	rows, err := mapAiCreditUsageReport(models.AiCreditUsageLevelEnterprise, 1, "avocado", raw)
	require.Nil(t, err)
	require.Len(t, rows, 2)

	first := rows[0]
	assert.Equal(t, uint64(1), first.ConnectionId)
	assert.Equal(t, "avocado", first.ScopeId)
	assert.Equal(t, models.AiCreditUsageLevelEnterprise, first.Level)
	assert.Equal(t, 2026, first.Year)
	assert.Equal(t, 6, first.Month)
	assert.Equal(t, 8, first.Day)
	assert.Equal(t, time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC), first.Date)
	assert.Equal(t, "Avocado Corp", first.Enterprise)
	assert.Equal(t, "Claude Opus 4.8", first.Model)
	assert.Equal(t, "ai-credits", first.UnitType)
	assert.InDelta(t, 5445.77595, first.GrossQuantity, 1e-9)
	assert.InDelta(t, 0.0, first.NetQuantity, 1e-9)

	second := rows[1]
	assert.Equal(t, "Claude Sonnet 4.6", second.Model)
	assert.InDelta(t, 1.62087, second.NetQuantity, 1e-9)
	assert.InDelta(t, 0.0162087, second.NetAmount, 1e-9)
}

func TestMapAiCreditUsageReportWithCostCenter(t *testing.T) {
	raw := []byte(`{
		"timePeriod": {"year": 2025, "month": 6, "day": 10},
		"enterprise": "acme",
		"organization": "platform",
		"user": "octocat",
		"costCenter": {"id": "cc-1", "name": "Engineering"},
		"usageItems": [
			{"product": "Copilot", "sku": "AI-CREDITS", "model": "gpt-4.1", "unitType": "credits",
			 "pricePerUnit": 1, "grossQuantity": 10, "grossAmount": 10, "discountQuantity": 0,
			 "discountAmount": 0, "netQuantity": 10, "netAmount": 10}
		]
	}`)

	rows, err := mapAiCreditUsageReport(models.AiCreditUsageLevelOrganization, 2, "platform", raw)
	require.Nil(t, err)
	require.Len(t, rows, 1)

	row := rows[0]
	assert.Equal(t, "cc-1", row.CostCenterId)
	assert.Equal(t, "Engineering", row.CostCenterName)
	assert.Equal(t, "platform", row.Organization)
	assert.Equal(t, "octocat", row.UserLogin)
	assert.InDelta(t, 10.0, row.NetAmount, 1e-9)
}

func TestMapAiCreditUsageReportEmpty(t *testing.T) {
	rows, err := mapAiCreditUsageReport(models.AiCreditUsageLevelEnterprise, 1, "x", []byte(`{"timePeriod":{"year":2026,"month":6,"day":8},"usageItems":[]}`))
	require.Nil(t, err)
	assert.Len(t, rows, 0)
}
